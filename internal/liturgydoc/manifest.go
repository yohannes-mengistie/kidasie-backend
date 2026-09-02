package liturgydoc

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ManifestEntry names one liturgy to build from the normalized sources.
type ManifestEntry struct {
	Slug           string
	SourceFile     string
	Name           string
	NameAm         string
	SectionTitle   string
	SectionTitleAm string
}

const manifestFields = 6

// DecodeManifest reads the tab-separated liturgy manifest. Tabs rather than
// spaces because both the source file names and the liturgy titles contain
// spaces. Blank lines and lines opening with # are comments.
func DecodeManifest(reader io.Reader) ([]ManifestEntry, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	entries := make([]ManifestEntry, 0, 16)
	seen := make(map[string]int)
	line := 0

	for scanner.Scan() {
		line++

		text := strings.TrimRight(scanner.Text(), "\r")
		if trimmed := strings.TrimSpace(text); trimmed == "" ||
			strings.HasPrefix(trimmed, "#") {
			continue
		}

		fields := strings.Split(text, "\t")
		if len(fields) != manifestFields {
			return nil, fmt.Errorf(
				"manifest line %d: expected %d tab-separated fields, found %d",
				line,
				manifestFields,
				len(fields),
			)
		}

		entry := ManifestEntry{
			Slug:           strings.TrimSpace(fields[0]),
			SourceFile:     strings.TrimSpace(fields[1]),
			Name:           strings.TrimSpace(fields[2]),
			NameAm:         strings.TrimSpace(fields[3]),
			SectionTitle:   strings.TrimSpace(fields[4]),
			SectionTitleAm: strings.TrimSpace(fields[5]),
		}

		for label, value := range map[string]string{
			"slug":             entry.Slug,
			"source file":      entry.SourceFile,
			"name":             entry.Name,
			"name_am":          entry.NameAm,
			"section title":    entry.SectionTitle,
			"section title_am": entry.SectionTitleAm,
		} {
			if value == "" {
				return nil, fmt.Errorf(
					"manifest line %d: %s is required",
					line,
					label,
				)
			}
		}

		if previous, duplicate := seen[entry.Slug]; duplicate {
			return nil, fmt.Errorf(
				"manifest line %d: slug %q already used on line %d",
				line,
				entry.Slug,
				previous,
			)
		}

		seen[entry.Slug] = line
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read liturgy manifest: %w", err)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("liturgy manifest is empty")
	}

	return entries, nil
}

// Options builds the conversion options for a manifest entry.
func (m ManifestEntry) Options(pagination Budget) Options {
	return Options{
		Slug:           m.Slug,
		Name:           m.Name,
		NameAm:         m.NameAm,
		SectionTitle:   m.SectionTitle,
		SectionTitleAm: m.SectionTitleAm,
		Pagination:     pagination,
	}
}
