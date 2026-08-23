package liturgyguide

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func Decode(reader io.Reader) ([]Entry, error) {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	var entries []Entry
	if err := decoder.Decode(&entries); err != nil {
		return nil, fmt.Errorf("decode liturgy guide: %w", err)
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err != nil {
			return nil, fmt.Errorf("decode trailing liturgy guide content: %w", err)
		}

		return nil, fmt.Errorf(
			"liturgy guide file must contain exactly one JSON document",
		)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("liturgy guide must contain at least one entry")
	}

	previousPage := 0
	for index := range entries {
		entry := entries[index]
		expectedID := index + 1

		if entry.ID != expectedID {
			return nil, fmt.Errorf(
				"entries[%d]: id must be %d",
				index,
				expectedID,
			)
		}
		if entry.Page <= 0 {
			return nil, fmt.Errorf(
				"entries[%d]: page must be positive",
				index,
			)
		}
		if entry.Page < previousPage {
			return nil, fmt.Errorf(
				"entries[%d]: page %d is before previous page %d",
				index,
				entry.Page,
				previousPage,
			)
		}
		if !validType(entry.Type) {
			return nil, fmt.Errorf(
				"entries[%d]: unsupported type %q",
				index,
				entry.Type,
			)
		}
		if !validLanguage(entry.Language) {
			return nil, fmt.Errorf(
				"entries[%d]: unsupported language %q",
				index,
				entry.Language,
			)
		}
		if strings.TrimSpace(entry.Text) == "" {
			return nil, fmt.Errorf(
				"entries[%d]: text is required",
				index,
			)
		}

		previousPage = entry.Page
	}

	return entries, nil
}

func validType(value string) bool {
	switch strings.TrimSpace(value) {
	case "heading", "paragraph", "list_item":
		return true
	default:
		return false
	}
}

func validLanguage(value string) bool {
	switch strings.TrimSpace(value) {
	case "amharic", "english":
		return true
	default:
		return false
	}
}
