package slidecontent

import (
	"fmt"
	"strings"

	"github.com/yohannes/kidasie-backend/internal/contentimport"
	"github.com/yohannes/kidasie-backend/internal/domain"
)

type Options struct {
	Slug           string
	Name           string
	NameAm         string
	SectionTitle   string
	SectionTitleAm string
}

type Stats struct {
	Pages            int
	Segments         int
	PagesWithoutText int
	ReviewSegments   int
	UntimedSegments  int
}

type text struct {
	geez    string
	amharic string
	english string
}

func Convert(
	source *Document,
	options Options,
) (*contentimport.Document, Stats, error) {
	if source == nil {
		return nil, Stats{}, fmt.Errorf("slide document is required")
	}

	verses := make([]contentimport.Verse, 0, len(source.Pages))
	stats := Stats{Pages: len(source.Pages)}

	for pageIndex := range source.Pages {
		page := source.Pages[pageIndex]
		before := len(verses)

		appendPageSegments(&verses, page)

		if len(verses) == before {
			stats.PagesWithoutText++
		}
	}

	for index := range verses {
		verses[index].Order = index + 1
		if verses[index].SourceNeedsReview {
			stats.ReviewSegments++
		}
		if verses[index].StartMs == nil {
			stats.UntimedSegments++
		}
	}

	stats.Segments = len(verses)

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
		return nil, Stats{}, fmt.Errorf("validate converted content: %w", err)
	}

	return document, stats, nil
}

func appendPageSegments(
	verses *[]contentimport.Verse,
	page Page,
) {
	main := text{
		geez:    page.TextGeez,
		amharic: page.TextAmharic,
		english: page.TextEnglish,
	}
	title := text{
		geez:    page.TitleGeez,
		amharic: page.TitleAmharic,
		english: firstNonEmpty(page.TitleEnglish, page.Title),
	}

	if hasText(main) {
		appendVerse(
			verses,
			page,
			"main",
			page.Kind,
			normalizeRole(page.Role, page.Kind),
			main,
		)
	} else if hasText(title) {
		appendVerse(
			verses,
			page,
			"title",
			page.Kind+":title",
			domain.RoleRubric,
			title,
		)
	}

	for index := range page.Parts {
		part := page.Parts[index]
		appendVerse(
			verses,
			page,
			fmt.Sprintf("part-%d", index+1),
			page.Kind+":part",
			normalizeRole(part.Role, page.Kind),
			text{
				geez:    part.TextGeez,
				amharic: part.TextAmharic,
				english: part.TextEnglish,
			},
		)
	}

	instruction := text{
		amharic: page.InstructionAmharic,
		english: page.Instruction,
	}
	if hasText(instruction) {
		appendVerse(
			verses,
			page,
			"instruction",
			page.Kind+":instruction",
			domain.RoleRubric,
			instruction,
		)
	}

	if strings.TrimSpace(page.Reference) != "" {
		appendVerse(
			verses,
			page,
			"reference",
			page.Kind+":reference",
			domain.RoleReader,
			text{english: page.Reference},
		)
	}

	deaconInstruction := text{
		geez:    page.DeaconInstructionGeez,
		amharic: page.DeaconInstructionAmharic,
		english: page.DeaconInstructionEnglish,
	}
	if hasText(deaconInstruction) {
		appendVerse(
			verses,
			page,
			"deacon-instruction",
			page.Kind+":deacon-instruction",
			domain.RoleDeacon,
			deaconInstruction,
		)
	}

	response := text{
		geez: joinUnique(
			page.ResponsePeopleGeez,
			page.ResponsePeople,
			page.TextGeezPeople,
		),
		amharic: joinUnique(
			page.ResponsePeopleAmharic,
			page.ResponseAmharic,
			page.TextAmharicPeople,
		),
		english: joinUnique(
			page.ResponsePeopleEnglish,
			page.ResponseEnglish,
			page.TextEnglishPeople,
		),
	}
	if hasText(response) {
		appendVerse(
			verses,
			page,
			"people-response",
			page.Kind+":people-response",
			domain.RoleCongregation,
			response,
		)
	}

	for index := range page.ResponseMixed {
		part := page.ResponseMixed[index]
		appendVerse(
			verses,
			page,
			fmt.Sprintf("mixed-response-%d", index+1),
			page.Kind+":mixed-response",
			normalizeRole(part.Role, page.Kind),
			text{
				geez:    part.TextGeez,
				amharic: part.TextAmharic,
				english: part.TextEnglish,
			},
		)
	}
}

func appendVerse(
	verses *[]contentimport.Verse,
	page Page,
	sourcePart string,
	sourceKind string,
	role string,
	value text,
) {
	value = cleanText(value)
	if !hasText(value) {
		return
	}

	pageNumber := page.Number

	*verses = append(*verses, contentimport.Verse{
		TextGeez:          value.geez,
		TextAm:            value.amharic,
		TextEn:            value.english,
		Role:              role,
		SourcePage:        &pageNumber,
		SourcePart:        sourcePart,
		SourceKind:        sourceKind,
		SourceNote:        strings.TrimSpace(page.Note),
		SourceNeedsReview: page.NeedsReview == nil || *page.NeedsReview,
	})
}

func normalizeRole(role string, kind string) string {
	role = strings.TrimSpace(role)

	if role == "people" {
		return domain.RoleCongregation
	}

	if domain.IsValidRole(role) {
		return role
	}

	switch kind {
	case "header", "anaphora_header", "instruction", "prayer_header":
		return domain.RoleRubric
	case "reading_announcement", "reading_ref", "scripture", "gospel_intro":
		return domain.RoleReader
	default:
		return domain.RoleMixed
	}
}

func cleanText(value text) text {
	return text{
		geez:    strings.TrimSpace(value.geez),
		amharic: strings.TrimSpace(value.amharic),
		english: strings.TrimSpace(value.english),
	}
}

func hasText(value text) bool {
	return strings.TrimSpace(value.geez) != "" ||
		strings.TrimSpace(value.amharic) != "" ||
		strings.TrimSpace(value.english) != ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}

func joinUnique(values ...string) string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))

	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}

		seen[value] = struct{}{}
		result = append(result, value)
	}

	return strings.Join(result, "\n\n")
}
