package contentimport

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/yohannes/kidasie-backend/internal/domain"
)

var (
	slugPattern   = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	sha256Pattern = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

func (d Document) Validate() error {
	if !slugPattern.MatchString(d.Slug) {
		return fmt.Errorf(
			"slug must contain lowercase letters, numbers, and hyphens",
		)
	}

	if strings.TrimSpace(d.Name) == "" {
		return fmt.Errorf("name is required")
	}

	if strings.TrimSpace(d.NameAm) == "" {
		return fmt.Errorf("name_am is required")
	}

	if d.Audio != nil {
		if err := d.Audio.Validate(); err != nil {
			return fmt.Errorf("audio: %w", err)
		}
	}

	for i := range d.Sections {
		expectedOrder := i + 1

		if d.Sections[i].Order != expectedOrder {
			return fmt.Errorf(
				"sections[%d]: order must be %d",
				i,
				expectedOrder,
			)
		}

		if err := d.Sections[i].validate(); err != nil {
			return fmt.Errorf("sections[%d]: %w", i, err)
		}
	}

	return nil
}

func (s Section) validate() error {
	if strings.TrimSpace(s.Title) == "" {
		return fmt.Errorf("title is required")
	}

	if strings.TrimSpace(s.TitleAm) == "" {
		return fmt.Errorf("title_am is required")
	}

	if s.Audio != nil {
		if err := s.Audio.Validate(); err != nil {
			return fmt.Errorf("audio: %w", err)
		}
	}

	if len(s.Verses) == 0 {
		return fmt.Errorf("at least one verse is required")
	}

	for i := range s.Verses {
		expectedOrder := i + 1
		verse := s.Verses[i]

		if verse.Order != expectedOrder {
			return fmt.Errorf(
				"verses[%d]: order must be %d",
				i,
				expectedOrder,
			)
		}

		if err := verse.validate(); err != nil {
			return fmt.Errorf("verses[%d]: %w", i, err)
		}

		if s.Audio != nil &&
			(verse.StartMs == nil || verse.EndMs == nil) {
			return fmt.Errorf(
				"verses[%d]: timing is required when section audio is present",
				i,
			)
		}

		if i > 0 {
			previous := s.Verses[i-1]

			if verse.StartMs != nil &&
				previous.EndMs != nil &&
				*verse.StartMs < *previous.EndMs {
				return fmt.Errorf(
					"verses[%d]: timing overlaps the previous verse",
					i,
				)
			}
		}
	}

	return nil
}

func (a Audio) Validate() error {
	parsedURL, err := url.Parse(a.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "http" &&
		parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must use http or https")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL host is required")
	}

	if a.DurationMs <= 0 {
		return fmt.Errorf("duration_ms must be positive")
	}

	if a.SizeBytes <= 0 {
		return fmt.Errorf("size_bytes must be positive")
	}

	if !strings.HasPrefix(a.MIMEType, "audio/") {
		return fmt.Errorf("mime_type must start with audio/")
	}

	if !sha256Pattern.MatchString(a.SHA256) {
		return fmt.Errorf(
			"sha256 must contain exactly 64 lowercase hexadecimal characters",
		)
	}

	return nil
}

func (v Verse) validate() error {
	if strings.TrimSpace(v.TextGeez) == "" &&
		strings.TrimSpace(v.TextAm) == "" &&
		strings.TrimSpace(v.TextEn) == "" {
		return fmt.Errorf("at least one language text is required")
	}

	if !domain.IsValidRole(v.Role) {
		return fmt.Errorf("unsupported role %q", v.Role)
	}

	if (v.StartMs == nil) != (v.EndMs == nil) {
		return fmt.Errorf("start_ms and end_ms must both be set or both be omitted")
	}

	if v.StartMs == nil {
		if v.SourcePage != nil && *v.SourcePage <= 0 {
			return fmt.Errorf("source_page must be positive")
		}

		return nil
	}

	if *v.StartMs < 0 {
		return fmt.Errorf("start_ms must not be negative")
	}

	if *v.EndMs <= *v.StartMs {
		return fmt.Errorf("end_ms must be greater than start_ms")
	}

	if v.SourcePage != nil && *v.SourcePage <= 0 {
		return fmt.Errorf("source_page must be positive")
	}

	return nil
}
