package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/yohannes/kidasie-backend/internal/flatanaphora"
)

func main() {
	if err := run(); err != nil {
		slog.Error("flat anaphora conversion failed", "error", err)
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

	flag.StringVar(
		&sourcePath,
		"file",
		"",
		"path to flat anaphora JSON",
	)
	flag.StringVar(
		&outputPath,
		"out",
		"",
		"path for importer JSON",
	)
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
		return fmt.Errorf("open flat anaphora file: %w", err)
	}
	defer sourceFile.Close()

	entries, err := flatanaphora.Decode(sourceFile)
	if err != nil {
		return err
	}

	document, stats, err := flatanaphora.Convert(
		entries,
		flatanaphora.Options{
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
		"flat anaphora converted",
		"output",
		outputPath,
		"entries",
		stats.Entries,
		"ambiguous_text",
		stats.AmbiguousText,
		"review_segments",
		stats.ReviewSegments,
		"untimed_segments",
		stats.UntimedSegments,
	)

	return nil
}

func writeJSONAtomically(path string, value any) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	tempFile, err := os.CreateTemp(directory, ".anaphora-*.json")
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
