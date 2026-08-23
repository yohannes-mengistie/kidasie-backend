package liturgyguide

import (
	"fmt"
	"strings"

	"github.com/yohannes/kidasie-backend/internal/contentimport"
	"github.com/yohannes/kidasie-backend/internal/domain"
)

const mixedLanguageNote = "source labels this entry English but the text also contains Ethiopic characters"

func Convert(
	entries []Entry,
	options Options,
) (*contentimport.Document, Stats, error) {
	verses := make([]contentimport.Verse, 0, len(entries))
	stats := Stats{Entries: len(entries)}

	for index := range entries {
		entry := entries[index]
		value := strings.TrimSpace(entry.Text)
		textAmharic := ""
		textEnglish := ""
		note := "source_language=" + strings.TrimSpace(entry.Language)

		switch strings.TrimSpace(entry.Language) {
		case "amharic":
			textAmharic = value
		case "english":
			if containsEthiopic(value) {
				textAmharic = value
				note += "; " + mixedLanguageNote
				stats.MixedLanguage++
			} else {
				textEnglish = value
			}
		}

		role := domain.RoleReader
		if entry.Type == "heading" {
			role = domain.RoleRubric
		}

		page := entry.Page
		verses = append(verses, contentimport.Verse{
			Order:             index + 1,
			TextAm:            textAmharic,
			TextEn:            textEnglish,
			Role:              role,
			SourcePage:        &page,
			SourcePart:        fmt.Sprintf("entry-%d", entry.ID),
			SourceKind:        "liturgy-guide:" + entry.Type,
			SourceNote:        note,
			SourceNeedsReview: true,
		})
	}

	stats.ReviewSegments = len(verses)

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
			"validate converted liturgy guide: %w",
			err,
		)
	}

	return document, stats, nil
}

func containsEthiopic(value string) bool {
	for _, character := range value {
		if character >= '\u1200' && character <= '\u139F' {
			return true
		}
	}

	return false
}
