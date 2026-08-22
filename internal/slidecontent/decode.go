package slidecontent

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func Decode(reader io.Reader) (*Document, error) {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	var document Document
	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("decode slide document: %w", err)
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err != nil {
			return nil, fmt.Errorf("decode trailing slide content: %w", err)
		}

		return nil, fmt.Errorf("slide file must contain exactly one JSON document")
	}

	if document.SchemaVersion != 1 {
		return nil, fmt.Errorf(
			"unsupported slide schema_version %d",
			document.SchemaVersion,
		)
	}

	if len(document.Pages) == 0 {
		return nil, fmt.Errorf("slide document must contain at least one page")
	}

	for index := range document.Pages {
		page := document.Pages[index]
		expectedNumber := index + 1

		if page.Number != expectedNumber {
			return nil, fmt.Errorf(
				"pages[%d]: number must be %d",
				index,
				expectedNumber,
			)
		}

		if strings.TrimSpace(page.Kind) == "" {
			return nil, fmt.Errorf("pages[%d]: kind is required", index)
		}
	}

	return &document, nil
}
