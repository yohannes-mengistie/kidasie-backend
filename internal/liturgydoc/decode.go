package liturgydoc

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Decode reads a normalized liturgy source. The shared beginning is a bare
// object with only "entries"; the anaphoras add a "title". A top-level array
// is accepted so older hand-edited exports still load.
func Decode(reader io.Reader) (*Document, error) {
	payload, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read liturgy document: %w", err)
	}

	document, err := unmarshalDocument(payload)
	if err != nil {
		return nil, err
	}

	if len(document.Entries) == 0 {
		return nil, fmt.Errorf(
			"liturgy document must contain at least one entry",
		)
	}

	if err := resolveGroups(document.Entries); err != nil {
		return nil, err
	}

	return document, nil
}

func unmarshalDocument(payload []byte) (*Document, error) {
	for _, character := range payload {
		switch character {
		case ' ', '\t', '\r', '\n':
			continue
		case '[':
			var entries []Entry
			if err := json.Unmarshal(payload, &entries); err != nil {
				return nil, fmt.Errorf("decode liturgy entries: %w", err)
			}

			return &Document{Entries: entries}, nil
		}

		break
	}

	var document Document
	if err := json.Unmarshal(payload, &document); err != nil {
		return nil, fmt.Errorf("decode liturgy document: %w", err)
	}

	return &document, nil
}

// resolveGroups parses every entry number and checks that the distinct
// numbers form the sequence 1..N. A gap or a repeat after the group has moved
// on means the source lost or duplicated an entry, which pagination would
// otherwise hide.
func resolveGroups(entries []Entry) error {
	previous := 0

	for index := range entries {
		entry := &entries[index]

		group, err := parseNumber(entry.Number)
		if err != nil {
			return fmt.Errorf("entries[%d]: %w", index, err)
		}

		switch {
		case group == previous:
		case group == previous+1:
			previous = group
		default:
			return fmt.Errorf(
				"entries[%d]: number %q resolves to %d, expected %d or %d",
				index,
				entry.Number,
				group,
				previous,
				previous+1,
			)
		}

		entry.group = group

		if strings.TrimSpace(entry.Role) == "" {
			return fmt.Errorf("entries[%d]: role is required", index)
		}

		if !entry.hasText() {
			return fmt.Errorf(
				"entries[%d]: at least one text field is required",
				index,
			)
		}
	}

	return nil
}

func (e Entry) hasText() bool {
	return strings.TrimSpace(e.TextGeez) != "" ||
		strings.TrimSpace(e.TextAmharic) != "" ||
		strings.TrimSpace(e.TextEnglish) != "" ||
		e.hasRubric()
}

func (e Entry) hasRubric() bool {
	return strings.TrimSpace(e.RubricGeez) != "" ||
		strings.TrimSpace(e.RubricAmharic) != "" ||
		strings.TrimSpace(e.RubricEnglish) != ""
}
