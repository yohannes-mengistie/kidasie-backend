package flatanaphora

import (
	"fmt"
	"strings"

	"github.com/yohannes/kidasie-backend/internal/contentimport"
	"github.com/yohannes/kidasie-backend/internal/domain"
)

const ambiguousTextNote = "source does not distinguish or separate Ge'ez and Amharic text"

func Convert(
	entries []Entry,
	options Options,
) (*contentimport.Document, Stats, error) {
	verses := make([]contentimport.Verse, 0, len(entries))
	stats := Stats{Entries: len(entries)}

	for index := range entries {
		entry := entries[index]
		for partIndex, part := range entry.Parts {
			role, err := normalizeRole(part.Role)
			if err != nil {
				return nil, Stats{}, fmt.Errorf(
					"entries[%d].parts[%d]: %w",
					index,
					partIndex,
					err,
				)
			}

			appendVerse(
				&verses,
				&stats,
				entry,
				role,
				part.TextGeez,
				part.TextAmharic,
				part.TextEnglish,
				strings.TrimSpace(entry.RoleSource),
				fmt.Sprintf("entry-%d-part-%d", index+1, partIndex+1),
			)
		}

		if !hasDirectEntryText(entry) {
			continue
		}

		textGeez := joinUnique(entry.GeezText, entry.TextGeez)
		textAmharic := joinUnique(entry.AmharicText, entry.TextAmharic)
		textEnglish := joinUnique(entry.EnglishText, entry.TextEnglish)
		ethiopicText := strings.TrimSpace(entry.EthiopicText)
		note := strings.TrimSpace(entry.RoleSource)

		if ethiopicText != "" {
			textAmharic = joinUnique(textAmharic, ethiopicText)
			note = joinNote(note, ambiguousTextNote)
			stats.AmbiguousText++
		}

		role, err := normalizeRole(entry.Role)
		if err != nil {
			return nil, Stats{}, fmt.Errorf(
				"entries[%d]: %w",
				index,
				err,
			)
		}

		appendVerse(
			&verses,
			&stats,
			entry,
			role,
			textGeez,
			textAmharic,
			textEnglish,
			note,
			fmt.Sprintf("entry-%d", index+1),
		)
	}

	stats.UntimedSegments = len(verses)

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
		return nil, Stats{}, fmt.Errorf(
			"validate converted anaphora: %w",
			err,
		)
	}

	return document, stats, nil
}

func hasDirectEntryText(entry Entry) bool {
	return strings.TrimSpace(entry.EthiopicText) != "" ||
		strings.TrimSpace(entry.GeezText) != "" ||
		strings.TrimSpace(entry.TextGeez) != "" ||
		strings.TrimSpace(entry.AmharicText) != "" ||
		strings.TrimSpace(entry.TextAmharic) != "" ||
		strings.TrimSpace(entry.EnglishText) != "" ||
		strings.TrimSpace(entry.TextEnglish) != ""
}

func appendVerse(
	verses *[]contentimport.Verse,
	stats *Stats,
	entry Entry,
	role string,
	textGeez string,
	textAmharic string,
	textEnglish string,
	note string,
	sourcePart string,
) {
	page := entry.Page
	sourceKind := "flat-anaphora"
	if kind := strings.TrimSpace(entry.Kind); kind != "" {
		sourceKind += ":" + kind
	}
	requiresReview := entry.NeedsReview == nil || *entry.NeedsReview
	if requiresReview {
		stats.ReviewSegments++
	}

	*verses = append(*verses, contentimport.Verse{
		Order:             len(*verses) + 1,
		TextGeez:          strings.TrimSpace(textGeez),
		TextAm:            strings.TrimSpace(textAmharic),
		TextEn:            strings.TrimSpace(textEnglish),
		Role:              role,
		SourcePage:        &page,
		SourcePart:        sourcePart,
		SourceKind:        sourceKind,
		SourceNote:        strings.TrimSpace(note),
		SourceNeedsReview: requiresReview,
	})
}

func normalizeRole(role string) (string, error) {
	switch strings.TrimSpace(role) {
	case "", "mixed":
		return domain.RoleMixed, nil
	case "priest",
		"ካህን/Priest",
		"ካህናት/Priests",
		"የመፈተት ጸሎት/Prayer of Fraction":
		return domain.RolePriest, nil
	case "assistant_priest", "ንፍቅ ካህን/Assistant Priest":
		return domain.RoleAssistantPriest, nil
	case "deacon", "ዲያቆን/Deacon", "ዲያቆናት/Deacons":
		return domain.RoleDeacon, nil
	case "assistant_deacon",
		"ንፍቅ ዲያቆን/Assistant Deacon",
		"ንፍቅ ዲያቆናት/Assistant Deacons":
		return domain.RoleAssistantDeacon, nil
	case "congregation", "people", "ሕዝብ/People":
		return domain.RoleCongregation, nil
	case "rubric", "መመሪያ/Instruction":
		return domain.RoleRubric, nil
	default:
		return "", fmt.Errorf("unsupported role %q", role)
	}
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

func joinNote(values ...string) string {
	return strings.Join(nonEmpty(values), "; ")
}

func nonEmpty(values []string) []string {
	result := make([]string, 0, len(values))

	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			result = append(result, value)
		}
	}

	return result
}
