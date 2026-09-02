package liturgydoc

import (
	"fmt"
	"strings"

	"github.com/yohannes/kidasie-backend/internal/contentimport"
	"github.com/yohannes/kidasie-backend/internal/domain"
)

const (
	kindBeginning = "beginning"
	kindAnaphora  = "anaphora"
)

// Convert joins the shared beginning with one anaphora and paginates the
// result into an importable document.
//
// The beginning is stored once, in content/updated/Qidase_serate.json, rather
// than copied into all fourteen anaphoras. Every anaphora restarts its source
// numbering at ፩, so the joined pages are renumbered continuously from 1: the
// beginning occupies the first pages and the anaphora continues from where it
// ends.
func Convert(
	beginning *Document,
	anaphora *Document,
	options Options,
) (*contentimport.Document, Stats, error) {
	if beginning == nil || len(beginning.Entries) == 0 {
		return nil, Stats{}, fmt.Errorf("beginning document is required")
	}

	if anaphora == nil || len(anaphora.Entries) == 0 {
		return nil, Stats{}, fmt.Errorf("anaphora document is required")
	}

	stats := Stats{
		BeginningEntries: len(beginning.Entries),
		AnaphoraEntries:  len(anaphora.Entries),
		SourceGroups: lastGroup(beginning.Entries) +
			lastGroup(anaphora.Entries),
	}

	// Paginated separately so the anaphora always opens a fresh page instead
	// of sharing one with the tail of the shared beginning.
	beginningPages := paginate(beginning.Entries, options.Pagination, 1)
	anaphoraPages := paginate(
		anaphora.Entries,
		options.Pagination,
		len(beginningPages)+1,
	)

	verses := make([]contentimport.Verse, 0, stats.BeginningEntries+stats.AnaphoraEntries)

	for _, section := range []struct {
		kind  string
		pages []Page
	}{
		{kind: kindBeginning, pages: beginningPages},
		{kind: kindAnaphora, pages: anaphoraPages},
	} {
		for _, page := range section.pages {
			if err := appendPage(&verses, &stats, page, section.kind); err != nil {
				return nil, Stats{}, err
			}
		}
	}

	stats.Pages = len(beginningPages) + len(anaphoraPages)
	stats.Verses = len(verses)

	target := options.Pagination.target()
	for _, page := range append(append([]Page{}, beginningPages...), anaphoraPages...) {
		size := page.Runes()
		if size > target {
			stats.OversizePages++
		}
		if size > stats.LargestPageRunes {
			stats.LargestPageRunes = size
		}
	}

	document := &contentimport.Document{
		Slug:   strings.TrimSpace(options.Slug),
		Name:   strings.TrimSpace(options.Name),
		NameAm: strings.TrimSpace(options.NameAm),
		Sections: []contentimport.Section{
			{
				Order:   1,
				Title:   strings.TrimSpace(options.SectionTitle),
				TitleAm: strings.TrimSpace(options.SectionTitleAm),
				Verses:  verses,
			},
		},
	}

	if err := document.Validate(); err != nil {
		return nil, Stats{}, fmt.Errorf("validate joined liturgy: %w", err)
	}

	return document, stats, nil
}

func appendPage(
	verses *[]contentimport.Verse,
	stats *Stats,
	page Page,
	kind string,
) error {
	for index, entry := range page.Entries {
		role, err := normalizeRole(entry.Role)
		if err != nil {
			return fmt.Errorf(
				"page %d entry %d (source number %q): %w",
				page.Number,
				index+1,
				entry.Number,
				err,
			)
		}

		part := fmt.Sprintf("%s-%d-%d", kind, entry.group, index+1)
		sourceKind := kind
		if subtype := strings.TrimSpace(entry.Subtype); subtype != "" {
			sourceKind += ":" + subtype
		}

		// A rubric attached to an utterance is a stage direction for it, so it
		// is emitted as its own rubric verse immediately before the utterance
		// and on the same page, matching how standalone rubrics already read.
		if entry.hasRubric() {
			appendVerse(verses, stats, contentimport.Verse{
				TextGeez:   strings.TrimSpace(entry.RubricGeez),
				TextAm:     strings.TrimSpace(entry.RubricAmharic),
				TextEn:     strings.TrimSpace(entry.RubricEnglish),
				Role:       domain.RoleRubric,
				SourcePart: part + "-rubric",
				SourceKind: sourceKind + ":rubric",
			}, page.Number, entry)
		}

		if strings.TrimSpace(entry.TextGeez) == "" &&
			strings.TrimSpace(entry.TextAmharic) == "" &&
			strings.TrimSpace(entry.TextEnglish) == "" {
			continue
		}

		appendVerse(verses, stats, contentimport.Verse{
			TextGeez:   strings.TrimSpace(entry.TextGeez),
			TextAm:     strings.TrimSpace(entry.TextAmharic),
			TextEn:     strings.TrimSpace(entry.TextEnglish),
			Role:       role,
			SourcePart: part,
			SourceKind: sourceKind,
		}, page.Number, entry)
	}

	return nil
}

func appendVerse(
	verses *[]contentimport.Verse,
	stats *Stats,
	verse contentimport.Verse,
	pageNumber int,
	entry Entry,
) {
	page := pageNumber
	verse.Order = len(*verses) + 1
	verse.SourcePage = &page
	verse.SourceNote = joinNotes(entry.OCRNote, entry.CrossReference)

	// Unlike the older OCR dumps, these sources have been reviewed, so an
	// absent needs_review means reviewed rather than unknown.
	verse.SourceNeedsReview = entry.NeedsReview != nil && *entry.NeedsReview
	if verse.SourceNeedsReview {
		stats.ReviewSegments++
	}

	*verses = append(*verses, verse)
}

func lastGroup(entries []Entry) int {
	if len(entries) == 0 {
		return 0
	}

	return entries[len(entries)-1].group
}

func normalizeRole(role string) (string, error) {
	switch strings.TrimSpace(role) {
	case "priest":
		return domain.RolePriest, nil
	case "assistant_priest":
		return domain.RoleAssistantPriest, nil
	case "deacon":
		return domain.RoleDeacon, nil
	case "assistant_deacon":
		return domain.RoleAssistantDeacon, nil
	case "people", "congregation":
		return domain.RoleCongregation, nil
	case "rubric":
		return domain.RoleRubric, nil
	case "chanter":
		return domain.RoleChanter, nil
	case "reader":
		return domain.RoleReader, nil
	case "mixed":
		return domain.RoleMixed, nil
	default:
		return "", fmt.Errorf("unsupported role %q", role)
	}
}

func joinNotes(values ...string) string {
	notes := make([]string, 0, len(values))

	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			notes = append(notes, value)
		}
	}

	return strings.Join(notes, "; ")
}
