package slideextract

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"

	"github.com/yohannes/kidasie-backend/internal/domain"
)

type Options struct {
	PDFPath   string
	AudioPath string
	StartPage int
	EndPage   int
	DPI       int
	Workers   int
	Progress  func(done int, total int)
}

var pagesPattern = regexp.MustCompile(`(?m)^Pages:\s+(\d+)\s*$`)

func Extract(ctx context.Context, options Options) (*Draft, error) {
	if strings.TrimSpace(options.PDFPath) == "" {
		return nil, errors.New("PDF path is required")
	}

	if options.DPI == 0 {
		options.DPI = 180
	}
	if options.DPI < 72 || options.DPI > 400 {
		return nil, errors.New("DPI must be between 72 and 400")
	}

	if options.Workers == 0 {
		options.Workers = min(runtime.NumCPU(), 4)
	}
	if options.Workers < 1 || options.Workers > 16 {
		return nil, errors.New("workers must be between 1 and 16")
	}

	if err := requireCommands("pdfinfo", "pdftoppm", "tesseract"); err != nil {
		return nil, err
	}

	totalPages, err := pdfPageCount(ctx, options.PDFPath)
	if err != nil {
		return nil, err
	}

	startPage := options.StartPage
	if startPage == 0 {
		startPage = 1
	}

	endPage := options.EndPage
	if endPage == 0 {
		endPage = totalPages
	}

	if startPage < 1 || endPage > totalPages || startPage > endPage {
		return nil, fmt.Errorf(
			"page range %d-%d is outside PDF pages 1-%d",
			startPage,
			endPage,
			totalPages,
		)
	}

	temporaryDirectory, err := os.MkdirTemp("", "kidasie-slide-extract-*")
	if err != nil {
		return nil, fmt.Errorf("create temporary extraction directory: %w", err)
	}
	defer os.RemoveAll(temporaryDirectory)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pageTotal := endPage - startPage + 1
	pages := make([]Page, pageTotal)

	var (
		nextPage  atomic.Int64
		completed atomic.Int64
		firstErr  error
		errOnce   sync.Once
		waitGroup sync.WaitGroup
	)

	nextPage.Store(int64(startPage - 1))

	for range options.Workers {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			for {
				pageNumber := int(nextPage.Add(1))
				if pageNumber > endPage {
					return
				}

				if err := ctx.Err(); err != nil {
					return
				}

				page, err := extractPage(
					ctx,
					options.PDFPath,
					temporaryDirectory,
					pageNumber,
					options.DPI,
				)
				if err != nil {
					errOnce.Do(func() {
						firstErr = err
						cancel()
					})
					return
				}

				pages[pageNumber-startPage] = page

				done := int(completed.Add(1))
				if options.Progress != nil {
					options.Progress(done, pageTotal)
				}
			}
		}()
	}

	waitGroup.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	sort.Slice(pages, func(i int, j int) bool {
		return pages[i].Number < pages[j].Number
	})

	propagateRoleSuggestions(pages)

	draft := &Draft{
		SchemaVersion: 1,
		SourcePDF:     filepath.Base(options.PDFPath),
		Pages:         pages,
	}

	if strings.TrimSpace(options.AudioPath) != "" {
		audio, err := inspectAudio(options.AudioPath)
		if err != nil {
			return nil, err
		}
		draft.Audio = audio
	}

	return draft, nil
}

func extractPage(
	ctx context.Context,
	pdfPath string,
	temporaryDirectory string,
	pageNumber int,
	dpi int,
) (Page, error) {
	outputPrefix := filepath.Join(
		temporaryDirectory,
		fmt.Sprintf("page-%04d", pageNumber),
	)

	renderCommand := exec.CommandContext(
		ctx,
		"pdftoppm",
		"-f",
		strconv.Itoa(pageNumber),
		"-l",
		strconv.Itoa(pageNumber),
		"-png",
		"-gray",
		"-r",
		strconv.Itoa(dpi),
		"-singlefile",
		pdfPath,
		outputPrefix,
	)

	if output, err := renderCommand.CombinedOutput(); err != nil {
		return Page{}, fmt.Errorf(
			"render PDF page %d: %w: %s",
			pageNumber,
			err,
			strings.TrimSpace(string(output)),
		)
	}

	imagePath := outputPrefix + ".png"

	ocrCommand := exec.CommandContext(
		ctx,
		"tesseract",
		imagePath,
		"stdout",
		"-l",
		"amh+eng",
		"--psm",
		"6",
		"-c",
		"preserve_interword_spaces=1",
	)

	ocrOutput, err := ocrCommand.Output()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return Page{}, fmt.Errorf(
				"OCR PDF page %d: %w: %s",
				pageNumber,
				err,
				strings.TrimSpace(string(exitError.Stderr)),
			)
		}

		return Page{}, fmt.Errorf("OCR PDF page %d: %w", pageNumber, err)
	}

	_ = os.Remove(imagePath)

	rawText := string(ocrOutput)
	ethiopicText := selectEthiopicLines(rawText)
	englishText := selectLines(rawText, isMostlyLatin)
	roles := detectRoles(englishText + "\n" + ethiopicText)

	page := Page{
		Number:          pageNumber,
		TextEthiopicOCR: ethiopicText,
		TextEnglishOCR:  englishText,
		NeedsReview:     true,
	}

	switch {
	case ethiopicText == "" && englishText == "":
		page.Kind = KindEmpty
		page.ExtractionWarning = "OCR found no text"
	case hasEmbeddedRole(rawText):
		page.Kind = KindMixed
		page.ExtractionWarning = "speaker role appears after content; split this slide"

		if len(roles) == 1 {
			page.RoleSuggestion = roles[0]
		}
	case len(roles) > 1:
		page.Kind = KindMixed
		page.ExtractionWarning = "multiple speaker roles detected"
	case len(roles) == 1:
		page.RoleSuggestion = roles[0]

		if wordCount(ethiopicText)+wordCount(englishText) <= 24 {
			page.Kind = KindRoleHeader
		} else {
			page.Kind = KindContent
		}
	default:
		page.Kind = KindContent
	}

	return page, nil
}

func propagateRoleSuggestions(pages []Page) {
	currentRole := ""

	for i := range pages {
		page := &pages[i]

		if page.Kind == KindRoleHeader && page.RoleSuggestion != "" {
			currentRole = page.RoleSuggestion
			continue
		}

		if page.RoleSuggestion != "" {
			currentRole = page.RoleSuggestion
			continue
		}

		if page.Kind == KindContent && currentRole != "" {
			page.RoleSuggestion = currentRole
			page.ExtractionWarning = appendWarning(
				page.ExtractionWarning,
				"role inherited from an earlier role header",
			)
		}
	}
}

func detectRoles(text string) []string {
	seen := make(map[string]struct{})

	for _, rawLine := range strings.Split(text, "\n") {
		line := strings.ToLower(strings.TrimSpace(rawLine))
		line = strings.TrimLeft(line, "-–—| ")

		switch {
		case strings.HasPrefix(line, "ረዳተ ካህን"),
			strings.HasPrefix(line, "assistant priest:"),
			strings.HasPrefix(line, "assistant priest :-"),
			strings.HasPrefix(line, "asst. priest:"),
			strings.HasPrefix(line, "asst. priest :-"):
			seen[domain.RoleAssistantPriest] = struct{}{}
		case strings.HasPrefix(line, "ካህን"),
			strings.HasPrefix(line, "priest:"),
			strings.HasPrefix(line, "priest :-"):
			seen[domain.RolePriest] = struct{}{}
		case strings.HasPrefix(line, "ዲያቆን"),
			strings.HasPrefix(line, "deacon:"),
			strings.HasPrefix(line, "deacon :-"):
			seen[domain.RoleDeacon] = struct{}{}
		case strings.HasPrefix(line, "ሕዝብ"),
			strings.HasPrefix(line, "ምእመናን"),
			strings.HasPrefix(line, "people:"),
			strings.HasPrefix(line, "congregation:"):
			seen[domain.RoleCongregation] = struct{}{}
		case strings.HasPrefix(line, "reader:"):
			seen[domain.RoleReader] = struct{}{}
		case strings.HasPrefix(line, "chanter:"):
			seen[domain.RoleChanter] = struct{}{}
		}
	}

	roles := make([]string, 0, len(seen))
	for role := range seen {
		roles = append(roles, role)
	}
	sort.Strings(roles)

	return roles
}

func hasEmbeddedRole(text string) bool {
	contentSeen := false

	for _, rawLine := range strings.Split(text, "\n") {
		line := strings.Join(strings.Fields(rawLine), " ")
		if line == "" || isPageNumber(line) {
			continue
		}

		if len(detectRoles(line)) > 0 {
			if contentSeen {
				return true
			}

			continue
		}

		if wordCount(line) > 2 {
			contentSeen = true
		}
	}

	return false
}

func selectEthiopicLines(text string) string {
	lines := make([]string, 0)

	for _, rawLine := range strings.Split(text, "\n") {
		words := make([]string, 0)

		for _, word := range strings.Fields(rawLine) {
			if hasEthiopic(word) {
				words = append(words, word)
			}
		}

		if len(words) > 0 {
			lines = append(lines, strings.Join(words, " "))
		}
	}

	return strings.Join(lines, "\n")
}

func isMostlyLatin(text string) bool {
	latinCount := 0
	ethiopicCount := 0

	for _, value := range text {
		switch {
		case unicode.In(value, unicode.Latin):
			latinCount++
		case value >= '\u1200' && value <= '\u139F':
			ethiopicCount++
		}
	}

	return latinCount > 0 && latinCount >= ethiopicCount
}

func selectLines(text string, predicate func(string) bool) string {
	lines := make([]string, 0)

	for _, rawLine := range strings.Split(text, "\n") {
		line := strings.Join(strings.Fields(rawLine), " ")
		if line == "" || isPageNumber(line) || !predicate(line) {
			continue
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func hasEthiopic(text string) bool {
	for _, value := range text {
		if value >= '\u1200' && value <= '\u139F' {
			return true
		}
	}

	return false
}

func hasLatin(text string) bool {
	for _, value := range text {
		if unicode.In(value, unicode.Latin) {
			return true
		}
	}

	return false
}

func isPageNumber(text string) bool {
	_, err := strconv.Atoi(strings.TrimSpace(text))
	return err == nil
}

func wordCount(text string) int {
	return len(strings.Fields(text))
}

func appendWarning(existing string, warning string) string {
	if existing == "" {
		return warning
	}

	return existing + "; " + warning
}

func requireCommands(names ...string) error {
	for _, name := range names {
		if _, err := exec.LookPath(name); err != nil {
			return fmt.Errorf("required command %q is not installed", name)
		}
	}

	return nil
}

func pdfPageCount(ctx context.Context, path string) (int, error) {
	command := exec.CommandContext(ctx, "pdfinfo", path)

	output, err := command.Output()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return 0, fmt.Errorf(
				"read PDF metadata: %w: %s",
				err,
				strings.TrimSpace(string(exitError.Stderr)),
			)
		}

		return 0, fmt.Errorf("read PDF metadata: %w", err)
	}

	match := pagesPattern.FindSubmatch(output)
	if len(match) != 2 {
		return 0, errors.New("PDF page count was not found")
	}

	count, err := strconv.Atoi(string(match[1]))
	if err != nil {
		return 0, fmt.Errorf("parse PDF page count: %w", err)
	}

	return count, nil
}

func inspectAudio(path string) (*AudioAsset, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open audio file: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("read audio file information: %w", err)
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, fmt.Errorf("calculate audio SHA-256: %w", err)
	}

	mimeType := "application/octet-stream"
	if strings.EqualFold(filepath.Ext(path), ".mp3") {
		mimeType = "audio/mpeg"
	}

	return &AudioAsset{
		FileName:  filepath.Base(path),
		SizeBytes: info.Size(),
		MIMEType:  mimeType,
		SHA256:    hex.EncodeToString(hash.Sum(nil)),
	}, nil
}
