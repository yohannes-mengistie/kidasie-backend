package flatanaphora

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
		return nil, fmt.Errorf("decode flat anaphora document: %w", err)
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err != nil {
			return nil, fmt.Errorf("decode trailing flat anaphora content: %w", err)
		}

		return nil, fmt.Errorf(
			"flat anaphora file must contain exactly one JSON document",
		)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf(
			"flat anaphora document must contain at least one entry",
		)
	}

	previousPage := 0
	for index := range entries {
		entry := &entries[index]
		if entry.Page > 0 && entry.Number > 0 && entry.Page != entry.Number {
			return nil, fmt.Errorf(
				"entries[%d]: page and number must match when both are present",
				index,
			)
		}
		if entry.Page == 0 {
			entry.Page = entry.Number
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

		if entry.Page > previousPage+1 {
			return nil, fmt.Errorf(
				"entries[%d]: page must be %d or %d",
				index,
				previousPage,
				previousPage+1,
			)
		}

		if !hasEntryText(*entry) {
			return nil, fmt.Errorf(
				"entries[%d]: at least one text field is required",
				index,
			)
		}

		for partIndex, part := range entry.Parts {
			if !hasPartText(part) {
				return nil, fmt.Errorf(
					"entries[%d].parts[%d]: at least one text field is required",
					index,
					partIndex,
				)
			}
		}

		previousPage = entry.Page
	}

	return entries, nil
}

func hasEntryText(entry Entry) bool {
	return strings.TrimSpace(entry.EthiopicText) != "" ||
		strings.TrimSpace(entry.GeezText) != "" ||
		strings.TrimSpace(entry.TextGeez) != "" ||
		strings.TrimSpace(entry.AmharicText) != "" ||
		strings.TrimSpace(entry.TextAmharic) != "" ||
		strings.TrimSpace(entry.EnglishText) != "" ||
		strings.TrimSpace(entry.TextEnglish) != "" ||
		len(entry.Parts) > 0
}

func hasPartText(part Part) bool {
	return strings.TrimSpace(part.TextGeez) != "" ||
		strings.TrimSpace(part.TextAmharic) != "" ||
		strings.TrimSpace(part.TextEnglish) != ""
}
