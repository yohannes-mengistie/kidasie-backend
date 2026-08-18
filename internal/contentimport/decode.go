package contentimport

import (
	"encoding/json"
	"fmt"
	"io"
)

func Decode(reader io.Reader) (*Document, error) {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	var document Document

	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("decode content document: %w", err)
	}

	var extraValue any

	err := decoder.Decode(&extraValue)
	if err != io.EOF {
		if err != nil {
			return nil, fmt.Errorf("decode trailing content: %w", err)
		}

		return nil, fmt.Errorf(
			"content file must contain exactly one JSON document",
		)
	}

	if err := document.Validate(); err != nil {
		return nil, fmt.Errorf("validate content document: %w", err)
	}

	return &document, nil
}
