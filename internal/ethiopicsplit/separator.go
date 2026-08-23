package ethiopicsplit

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"math"
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

	"github.com/yohannes/kidasie-backend/internal/flatanaphora"
)

type Options struct {
	DPI           int
	Workers       int
	MinConfidence float64
	MetadataOnly  bool
	Progress      func(done int, total int)
}

type Report struct {
	SourcePDF  string       `json:"source_pdf"`
	Total      int          `json:"total_ethiopic_entries"`
	Separated  int          `json:"separated"`
	Unresolved int          `json:"unresolved"`
	Entries    []ReportItem `json:"entries"`
}

type ReportItem struct {
	Index      int     `json:"index"`
	Page       int     `json:"page"`
	Role       string  `json:"role"`
	Status     string  `json:"status"`
	Confidence float64 `json:"confidence,omitempty"`
	Note       string  `json:"note"`
}

type pageOCR struct {
	page   int
	blocks []string
}

type suggestion struct {
	geez       string
	amharic    string
	confidence float64
}

type pdfDocument struct {
	Pages []pdfPage `xml:"page"`
}

type pdfPage struct {
	Number int           `xml:"number,attr"`
	Fonts  []pdfFontSpec `xml:"fontspec"`
	Texts  []pdfText     `xml:"text"`
}

type pdfFontSpec struct {
	ID     string `xml:"id,attr"`
	Family string `xml:"family,attr"`
	Color  string `xml:"color,attr"`
}

type pdfText struct {
	Top    int    `xml:"top,attr"`
	Left   int    `xml:"left,attr"`
	Height int    `xml:"height,attr"`
	Font   string `xml:"font,attr"`
	Value  string
}

type evidenceBlock struct {
	text   string
	color  string
	glyphs int
	top    int
	bottom int
}

var blankLines = regexp.MustCompile(`\r?\n[\t ]*\r?\n+`)

func Separate(
	ctx context.Context,
	pdfPath string,
	entries []flatanaphora.Entry,
	options Options,
) ([]flatanaphora.Entry, Report, error) {
	if strings.TrimSpace(pdfPath) == "" {
		return nil, Report{}, errors.New("PDF path is required")
	}
	if options.DPI == 0 {
		options.DPI = 220
	}
	if options.DPI < 120 || options.DPI > 400 {
		return nil, Report{}, errors.New("DPI must be between 120 and 400")
	}
	if options.Workers == 0 {
		options.Workers = min(runtime.NumCPU(), 4)
	}
	if options.Workers < 1 || options.Workers > 16 {
		return nil, Report{}, errors.New("workers must be between 1 and 16")
	}
	if options.MinConfidence == 0 {
		options.MinConfidence = 0.72
	}
	if options.MinConfidence < 0.5 || options.MinConfidence > 1 {
		return nil, Report{}, errors.New("minimum confidence must be between 0.5 and 1")
	}
	if err := requireCommands("pdftohtml", "pdftoppm", "tesseract"); err != nil {
		return nil, Report{}, err
	}

	result := append([]flatanaphora.Entry(nil), entries...)
	report := Report{SourcePDF: filepath.Base(pdfPath)}
	for index := range result {
		entry := &result[index]
		combined := strings.TrimSpace(entry.EthiopicText)
		if combined == "" {
			continue
		}

		report.Total++
		if hasGeez(*entry) && strings.TrimSpace(entry.AmharicText) == "" {
			entry.AmharicText = combined
			entry.OriginalEthiopicText = combined
			entry.EthiopicText = ""
			entry.SeparationConfidence = 1
			entry.SeparationNote = "existing Ge'ez field confirms Ethiopic text is Amharic"
			report.Separated++
			report.Entries = append(report.Entries, ReportItem{
				Index: index, Page: entry.Page, Role: entry.Role,
				Status: "separated", Confidence: 1,
				Note: "existing Ge'ez field confirms Ethiopic text is Amharic",
			})
			continue
		}

	}

	metadata, err := metadataSuggestions(ctx, pdfPath, result)
	if err != nil {
		return nil, Report{}, err
	}

	pagesNeeded := make(map[int]struct{})
	for index := range result {
		entry := &result[index]
		combined := strings.TrimSpace(entry.EthiopicText)
		if combined == "" || hasGeez(*entry) {
			continue
		}
		if candidate, ok := metadata[index]; ok &&
			candidate.confidence >= options.MinConfidence {
			applySplit(
				entry,
				combined,
				candidate.geez,
				candidate.amharic,
				candidate.confidence,
				"separated from PDF font/color blocks; requires human review",
			)
			report.Separated++
			report.Entries = append(report.Entries, ReportItem{
				Index: index, Page: entry.Page, Role: entry.Role,
				Status: "separated", Confidence: candidate.confidence,
				Note: "separated from PDF font/color blocks; requires human review",
			})
			continue
		}
		pagesNeeded[entry.Page] = struct{}{}
	}

	pages := make(map[int][]string)
	if !options.MetadataOnly {
		pages, err = extractPages(ctx, pdfPath, pagesNeeded, options)
		if err != nil {
			return nil, Report{}, err
		}
	}

	for index := range result {
		entry := &result[index]
		combined := strings.TrimSpace(entry.EthiopicText)
		if combined == "" {
			continue
		}
		if hasGeez(*entry) {
			report.Unresolved++
			report.Entries = append(report.Entries, ReportItem{
				Index: index, Page: entry.Page, Role: entry.Role,
				Status: "unresolved",
				Note:   "entry already contains both explicit language fields and ambiguous Ethiopic text; original retained",
			})
			continue
		}
		if options.MetadataOnly {
			report.Unresolved++
			report.Entries = append(report.Entries, ReportItem{
				Index: index, Page: entry.Page, Role: entry.Role,
				Status: "unresolved",
				Note:   "PDF metadata did not prove a boundary; original text retained for OCR or human review",
			})
			continue
		}

		geez, amharic, confidence, ok := splitCombined(
			combined,
			pages[entry.Page],
			options.MinConfidence,
		)
		if !ok {
			report.Unresolved++
			report.Entries = append(report.Entries, ReportItem{
				Index: index, Page: entry.Page, Role: entry.Role,
				Status: "unresolved",
				Note:   "PDF OCR did not identify a reliable Ge'ez/Amharic boundary; original text retained",
			})
			continue
		}

		applySplit(
			entry,
			combined,
			geez,
			amharic,
			confidence,
			"separated from adjacent PDF text blocks; requires human review",
		)
		report.Separated++
		report.Entries = append(report.Entries, ReportItem{
			Index: index, Page: entry.Page, Role: entry.Role,
			Status: "separated", Confidence: confidence,
			Note: "separated from adjacent PDF text blocks; requires human review",
		})
	}

	sort.Slice(report.Entries, func(i, j int) bool {
		return report.Entries[i].Index < report.Entries[j].Index
	})

	return result, report, nil
}

func (value *pdfText) UnmarshalXML(
	decoder *xml.Decoder,
	start xml.StartElement,
) error {
	for _, attribute := range start.Attr {
		switch attribute.Name.Local {
		case "top":
			value.Top, _ = strconv.Atoi(attribute.Value)
		case "left":
			value.Left, _ = strconv.Atoi(attribute.Value)
		case "height":
			value.Height, _ = strconv.Atoi(attribute.Value)
		case "font":
			value.Font = attribute.Value
		}
	}

	var text strings.Builder
	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		switch token := token.(type) {
		case xml.CharData:
			text.Write([]byte(token))
		case xml.EndElement:
			if token.Name == start.Name {
				value.Value = strings.TrimSpace(text.String())
				return nil
			}
		}
	}
}

func metadataSuggestions(
	ctx context.Context,
	pdfPath string,
	entries []flatanaphora.Entry,
) (map[int]suggestion, error) {
	temporaryDirectory, err := os.MkdirTemp("", "kidasie-pdf-metadata-*")
	if err != nil {
		return nil, fmt.Errorf("create PDF metadata directory: %w", err)
	}
	defer os.RemoveAll(temporaryDirectory)

	outputPath := filepath.Join(temporaryDirectory, "document.xml")
	command := exec.CommandContext(
		ctx,
		"pdftohtml",
		"-xml",
		"-hidden",
		"-nodrm",
		pdfPath,
		outputPath,
	)
	if output, err := command.CombinedOutput(); err != nil {
		return nil, fmt.Errorf(
			"extract PDF metadata: %w: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	file, err := os.Open(outputPath)
	if err != nil {
		return nil, fmt.Errorf("open PDF metadata: %w", err)
	}
	var document pdfDocument
	decodeErr := xml.NewDecoder(file).Decode(&document)
	closeErr := file.Close()
	if decodeErr != nil {
		return nil, fmt.Errorf("decode PDF metadata: %w", decodeErr)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("close PDF metadata: %w", closeErr)
	}

	needed := make(map[int]struct{})
	for _, entry := range entries {
		if strings.TrimSpace(entry.EthiopicText) != "" && !hasGeez(entry) {
			needed[entry.Page] = struct{}{}
		}
	}

	pageBlocks := make(map[int][]evidenceBlock, len(needed))
	fonts := make(map[string]pdfFontSpec)
	for _, page := range document.Pages {
		for _, font := range page.Fonts {
			fonts[font.ID] = font
		}
		if _, ok := needed[page.Number]; !ok {
			continue
		}
		pageBlocks[page.Number] = buildEvidenceBlocks(page.Texts, fonts)
	}

	result := make(map[int]suggestion)
	for index, entry := range entries {
		combined := strings.TrimSpace(entry.EthiopicText)
		if combined == "" || hasGeez(entry) {
			continue
		}
		geez, amharic, confidence, ok := splitByEvidence(
			combined,
			pageBlocks[entry.Page],
		)
		if ok {
			result[index] = suggestion{
				geez: geez, amharic: amharic, confidence: confidence,
			}
		}
	}
	return result, nil
}

func buildEvidenceBlocks(
	texts []pdfText,
	fonts map[string]pdfFontSpec,
) []evidenceBlock {
	sort.Slice(texts, func(i, j int) bool {
		if texts[i].Top == texts[j].Top {
			return texts[i].Left < texts[j].Left
		}
		return texts[i].Top < texts[j].Top
	})

	type line struct {
		text   string
		color  string
		glyphs int
		top    int
		height int
	}
	lines := make([]line, 0)
	for _, text := range texts {
		font, ok := fonts[text.Font]
		if !ok || !isEthiopicFont(font.Family) {
			continue
		}
		value := strings.TrimSpace(text.Value)
		if value == "" {
			continue
		}
		if len(lines) > 0 &&
			lines[len(lines)-1].top == text.Top &&
			lines[len(lines)-1].color == font.Color {
			last := &lines[len(lines)-1]
			last.text = strings.TrimSpace(last.text + " " + value)
			last.glyphs += glyphCount(value)
			last.height = max(last.height, text.Height)
			continue
		}
		lines = append(lines, line{
			text: value, color: font.Color, glyphs: glyphCount(value),
			top: text.Top, height: text.Height,
		})
	}

	blocks := make([]evidenceBlock, 0, len(lines))
	for _, line := range lines {
		if line.glyphs < 2 || isPDFRoleLine(line.text) {
			continue
		}
		actualText := ""
		if hasEthiopic(line.text) {
			actualText = strings.Join(ethiopicWords(line.text), " ")
		}
		if len(blocks) > 0 {
			last := &blocks[len(blocks)-1]
			if last.color == line.color && line.top-last.bottom <= 110 {
				last.glyphs += line.glyphs
				last.bottom = line.top + line.height
				last.text = strings.TrimSpace(last.text + " " + actualText)
				continue
			}
		}
		blocks = append(blocks, evidenceBlock{
			text: actualText, color: line.color, glyphs: line.glyphs,
			top: line.top, bottom: line.top + line.height,
		})
	}
	return blocks
}

func isEthiopicFont(family string) bool {
	family = strings.ToLower(family)
	for _, name := range []string{
		"ethiopia", "washra", "wookianos", "nyala",
	} {
		if strings.Contains(family, name) {
			return true
		}
	}
	return false
}

func isPDFRoleLine(text string) bool {
	text = strings.ToLower(strings.TrimSpace(text))
	if strings.Contains(text, "/") {
		return true
	}
	for _, role := range []string{
		"ካህን", "ሕዝብ", "ምእመናን", "ዲያቆን", "priest", "people", "deacon",
	} {
		if strings.HasPrefix(text, role) && len([]rune(text)) < 45 {
			return true
		}
	}
	return false
}

func ethiopicWords(text string) []string {
	result := make([]string, 0)
	for _, word := range strings.Fields(text) {
		if hasEthiopic(word) {
			result = append(result, word)
		}
	}
	return result
}

func glyphCount(text string) int {
	count := 0
	for _, value := range text {
		if unicode.IsSpace(value) || unicode.IsPunct(value) ||
			unicode.In(value, unicode.Latin) || unicode.IsDigit(value) {
			continue
		}
		count++
	}
	return count
}

func splitByEvidence(
	combined string,
	blocks []evidenceBlock,
) (string, string, float64, bool) {
	words := strings.Fields(combined)
	totalRunes := len(normalizeEthiopic(combined))
	if len(words) < 2 || totalRunes == 0 {
		return "", "", 0, false
	}

	bestScore := 0.0
	bestBoundary := 0
	for index := 0; index+1 < len(blocks); index++ {
		leftBlock := blocks[index]
		rightBlock := blocks[index+1]
		if leftBlock.color == rightBlock.color ||
			leftBlock.glyphs < 5 || rightBlock.glyphs < 5 ||
			rightBlock.top-leftBlock.bottom > 120 {
			continue
		}

		if leftBlock.text != "" && rightBlock.text != "" {
			geez, amharic, confidence, ok := splitCombined(
				combined,
				[]string{leftBlock.text, rightBlock.text},
				0.72,
			)
			if ok && confidence > bestScore {
				bestScore = confidence
				bestBoundary = len(strings.Fields(geez))
				_ = amharic
			}
			continue
		}

		blockTotal := leftBlock.glyphs + rightBlock.glyphs
		totalScore := float64(min(totalRunes, blockTotal)) /
			float64(max(totalRunes, blockTotal))
		if totalScore < 0.6 {
			continue
		}
		targetRatio := float64(leftBlock.glyphs) / float64(blockTotal)
		for boundary := 1; boundary < len(words); boundary++ {
			leftRunes := len(normalizeEthiopic(
				strings.Join(words[:boundary], " "),
			))
			ratio := float64(leftRunes) / float64(totalRunes)
			ratioScore := 1 - min(1.0, math.Abs(ratio-targetRatio)*3)
			score := totalScore*0.65 + ratioScore*0.35
			if score > bestScore {
				bestScore = score
				bestBoundary = boundary
			}
		}
	}

	if bestBoundary == 0 || bestScore < 0.72 {
		return "", "", bestScore, false
	}
	return strings.Join(words[:bestBoundary], " "),
		strings.Join(words[bestBoundary:], " "),
		bestScore,
		true
}

func applySplit(
	entry *flatanaphora.Entry,
	original string,
	geez string,
	amharic string,
	confidence float64,
	note string,
) {
	entry.GeezText = strings.TrimSpace(geez)
	entry.AmharicText = strings.TrimSpace(amharic)
	entry.OriginalEthiopicText = original
	entry.EthiopicText = ""
	entry.SeparationConfidence = math.Round(confidence*1000) / 1000
	entry.SeparationNote = note
}

func hasGeez(entry flatanaphora.Entry) bool {
	return strings.TrimSpace(entry.GeezText) != "" ||
		strings.TrimSpace(entry.TextGeez) != ""
}

func extractPages(
	ctx context.Context,
	pdfPath string,
	pageSet map[int]struct{},
	options Options,
) (map[int][]string, error) {
	pages := make([]int, 0, len(pageSet))
	for page := range pageSet {
		pages = append(pages, page)
	}
	sort.Ints(pages)

	temporaryDirectory, err := os.MkdirTemp("", "kidasie-ethiopic-split-*")
	if err != nil {
		return nil, fmt.Errorf("create temporary OCR directory: %w", err)
	}
	defer os.RemoveAll(temporaryDirectory)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan int)
	results := make(chan pageOCR)
	errorsFound := make(chan error, 1)
	var workers sync.WaitGroup

	for range options.Workers {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for page := range jobs {
				blocks, err := extractPage(
					ctx,
					pdfPath,
					temporaryDirectory,
					page,
					options.DPI,
				)
				if err != nil {
					select {
					case errorsFound <- err:
					default:
					}
					cancel()
					return
				}
				select {
				case results <- pageOCR{page: page, blocks: blocks}:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for _, page := range pages {
			select {
			case jobs <- page:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		workers.Wait()
		close(results)
	}()

	extracted := make(map[int][]string, len(pages))
	var completed atomic.Int64
	for result := range results {
		extracted[result.page] = result.blocks
		done := int(completed.Add(1))
		if options.Progress != nil {
			options.Progress(done, len(pages))
		}
	}

	select {
	case err := <-errorsFound:
		return nil, err
	default:
	}
	if err := ctx.Err(); err != nil && !errors.Is(err, context.Canceled) {
		return nil, err
	}

	return extracted, nil
}

func extractPage(
	ctx context.Context,
	pdfPath string,
	temporaryDirectory string,
	page int,
	dpi int,
) ([]string, error) {
	prefix := filepath.Join(temporaryDirectory, fmt.Sprintf("page-%04d", page))
	render := exec.CommandContext(
		ctx,
		"pdftoppm",
		"-f",
		strconv.Itoa(page),
		"-l",
		strconv.Itoa(page),
		"-png",
		"-gray",
		"-r",
		strconv.Itoa(dpi),
		"-singlefile",
		pdfPath,
		prefix,
	)
	if output, err := render.CombinedOutput(); err != nil {
		return nil, fmt.Errorf(
			"render PDF page %d: %w: %s",
			page,
			err,
			strings.TrimSpace(string(output)),
		)
	}

	imagePath := prefix + ".png"
	defer os.Remove(imagePath)
	ocr := exec.CommandContext(
		ctx,
		"tesseract",
		imagePath,
		"stdout",
		"-l",
		"amh+eng",
		"--psm",
		"3",
		"-c",
		"preserve_interword_spaces=1",
	)
	output, err := ocr.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(
			"OCR PDF page %d: %w: %s",
			page,
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return extractEthiopicBlocks(string(output)), nil
}

func extractEthiopicBlocks(text string) []string {
	parts := blankLines.Split(strings.TrimSpace(text), -1)
	blocks := make([]string, 0, len(parts))
	for _, part := range parts {
		words := make([]string, 0)
		for _, line := range strings.Split(part, "\n") {
			line = strings.TrimSpace(line)
			if isRoleLine(line) {
				continue
			}
			for _, word := range strings.Fields(line) {
				if hasEthiopic(word) {
					words = append(words, word)
				}
			}
		}
		block := strings.Join(words, " ")
		if len(normalizeEthiopic(block)) >= 3 {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

func isRoleLine(line string) bool {
	line = strings.ToLower(strings.TrimSpace(line))
	for _, prefix := range []string{
		"ካህን",
		"ሕዝብ",
		"ምእመናን",
		"ዲያቆን",
		"ንፍቅ ዲያቆን",
		"priest",
		"people",
		"deacon",
		"assistant",
	} {
		if strings.HasPrefix(line, prefix) && len([]rune(line)) < 45 {
			return true
		}
	}
	return false
}

func splitCombined(
	combined string,
	blocks []string,
	minimum float64,
) (string, string, float64, bool) {
	words := strings.Fields(combined)
	if len(words) < 2 || len(blocks) < 2 {
		return "", "", 0, false
	}

	bestScore := 0.0
	bestBoundary := 0
	for index := 0; index+1 < len(blocks); index++ {
		leftOCR := normalizeEthiopic(blocks[index])
		rightOCR := normalizeEthiopic(blocks[index+1])
		if len(leftOCR) == 0 || len(rightOCR) == 0 {
			continue
		}

		expected := int(math.Round(
			float64(len(words)) * float64(len(leftOCR)) /
				float64(len(leftOCR)+len(rightOCR)),
		))
		window := max(5, len(words)/5)
		start := max(1, expected-window)
		end := min(len(words)-1, expected+window)

		for boundary := start; boundary <= end; boundary++ {
			left := normalizeEthiopic(strings.Join(words[:boundary], " "))
			right := normalizeEthiopic(strings.Join(words[boundary:], " "))
			leftScore := similarity(left, leftOCR)
			rightScore := similarity(right, rightOCR)
			total := append(append([]rune(nil), left...), right...)
			totalOCR := append(append([]rune(nil), leftOCR...), rightOCR...)
			totalScore := similarity(total, totalOCR)
			score := totalScore*0.5 + leftScore*0.25 + rightScore*0.25
			if min(leftScore, rightScore) < 0.48 {
				continue
			}
			if score > bestScore {
				bestScore = score
				bestBoundary = boundary
			}
		}
	}

	if bestBoundary == 0 || bestScore < minimum {
		return "", "", bestScore, false
	}

	return strings.Join(words[:bestBoundary], " "),
		strings.Join(words[bestBoundary:], " "),
		bestScore,
		true
}

func normalizeEthiopic(text string) []rune {
	result := make([]rune, 0, len(text))
	for _, value := range text {
		if value >= '\u1200' && value <= '\u2DDF' &&
			!unicode.IsPunct(value) {
			result = append(result, value)
		}
	}
	return result
}

func hasEthiopic(text string) bool {
	for _, value := range text {
		if value >= '\u1200' && value <= '\u2DDF' {
			return true
		}
	}
	return false
}

func similarity(left []rune, right []rune) float64 {
	maximum := max(len(left), len(right))
	if maximum == 0 {
		return 1
	}
	return 1 - float64(editDistance(left, right))/float64(maximum)
}

func editDistance(left []rune, right []rune) int {
	previous := make([]int, len(right)+1)
	current := make([]int, len(right)+1)
	for index := range previous {
		previous[index] = index
	}
	for leftIndex, leftValue := range left {
		current[0] = leftIndex + 1
		for rightIndex, rightValue := range right {
			cost := 1
			if leftValue == rightValue {
				cost = 0
			}
			current[rightIndex+1] = min(
				previous[rightIndex+1]+1,
				current[rightIndex]+1,
				previous[rightIndex]+cost,
			)
		}
		previous, current = current, previous
	}
	return previous[len(right)]
}

func requireCommands(commands ...string) error {
	for _, command := range commands {
		if _, err := exec.LookPath(command); err != nil {
			return fmt.Errorf(
				"required command %q was not found: %w",
				command,
				err,
			)
		}
	}
	return nil
}
