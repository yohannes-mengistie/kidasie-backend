package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/yohannes/kidasie-backend/internal/liturgyguide"
)

func main() {
	if err := run(); err != nil {
		slog.Error("liturgy guide conversion failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		sourcePath     string
		outputPath     string
		slug           string
		name           string
		nameAm         string
		sectionTitle   string
		sectionTitleAm string
	)

	flag.StringVar(&sourcePath, "file", "", "path to liturgy guide JSON")
	flag.StringVar(&outputPath, "out", "", "path for importer JSON")
	flag.StringVar(&slug, "slug", "", "liturgy slug")
	flag.StringVar(&name, "name", "", "English liturgy name")
	flag.StringVar(&nameAm, "name-am", "", "Amharic liturgy name")
	flag.StringVar(
		&sectionTitle,
		"section-title",
		"",
		"English section title",
	)
	flag.StringVar(
		&sectionTitleAm,
		"section-title-am",
		"",
		"Amharic section title",
	)
	flag.Parse()

	for label, value := range map[string]string{
		"-file":             sourcePath,
		"-out":              outputPath,
		"-slug":             slug,
		"-name":             name,
		"-name-am":          nameAm,
		"-section-title":    sectionTitle,
		"-section-title-am": sectionTitleAm,
	} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required", label)
		}
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open liturgy guide file: %w", err)
	}
	defer sourceFile.Close()

	entries, err := liturgyguide.Decode(sourceFile)
	if err != nil {
		return err
	}

	document, stats, err := liturgyguide.Convert(
		entries,
		liturgyguide.Options{
			Slug:           slug,
			Name:           name,
			NameAm:         nameAm,
			SectionTitle:   sectionTitle,
			SectionTitleAm: sectionTitleAm,
		},
	)
	if err != nil {
		return err
	}

	if err := writeJSONAtomically(outputPath, document); err != nil {
		return err
	}

	slog.Info(
		"liturgy guide converted",
		"output",
		outputPath,
		"entries",
		stats.Entries,
		"mixed_language",
		stats.MixedLanguage,
		"review_segments",
		stats.ReviewSegments,
	)

	return nil
}

func writeJSONAtomically(path string, value any) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	tempFile, err := os.CreateTemp(directory, ".guide-*.json")
	if err != nil {
		return fmt.Errorf("create temporary output: %w", err)
	}

	tempPath := tempFile.Name()
	defer os.Remove(tempPath)

	encoder := json.NewEncoder(tempFile)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(value); err != nil {
		tempFile.Close()
		return fmt.Errorf("encode importer JSON: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close importer JSON: %w", err)
	}

	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("replace importer JSON: %w", err)
	}

	return nil
}
