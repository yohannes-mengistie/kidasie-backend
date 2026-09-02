package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/yohannes/kidasie-backend/internal/liturgydoc"
)

func main() {
	if err := run(); err != nil {
		slog.Error("liturgy conversion failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		beginningPath  string
		manifestPath   string
		sourceDir      string
		outputDir      string
		anaphoraPath   string
		outputPath     string
		slug           string
		name           string
		nameAm         string
		sectionTitle   string
		sectionTitleAm string
		targetRunes    int
	)

	flag.StringVar(
		&beginningPath,
		"beginning",
		"content/updated/Qidase_serate.json",
		"path to the shared liturgy beginning JSON",
	)
	flag.StringVar(
		&manifestPath,
		"manifest",
		"",
		"tab-separated manifest; converts every liturgy it lists",
	)
	flag.StringVar(
		&sourceDir,
		"source-dir",
		"content/updated",
		"directory holding the manifest's source files",
	)
	flag.StringVar(
		&outputDir,
		"out-dir",
		"",
		"directory for the manifest's importer JSON",
	)
	flag.StringVar(&anaphoraPath, "file", "", "path to a single anaphora JSON")
	flag.StringVar(&outputPath, "out", "", "path for a single importer JSON")
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
	flag.IntVar(
		&targetRunes,
		"target-runes",
		liturgydoc.DefaultTargetRunes,
		"per-page text budget across all languages",
	)
	flag.Parse()

	if strings.TrimSpace(beginningPath) == "" {
		return fmt.Errorf("-beginning is required")
	}

	beginning, err := decodeFile(beginningPath)
	if err != nil {
		return fmt.Errorf("beginning %s: %w", beginningPath, err)
	}

	budget := liturgydoc.Budget{TargetRunes: targetRunes}

	if strings.TrimSpace(manifestPath) != "" {
		return convertManifest(
			beginning,
			manifestPath,
			sourceDir,
			outputDir,
			budget,
		)
	}

	for label, value := range map[string]string{
		"-file":             anaphoraPath,
		"-out":              outputPath,
		"-slug":             slug,
		"-name":             name,
		"-name-am":          nameAm,
		"-section-title":    sectionTitle,
		"-section-title-am": sectionTitleAm,
	} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required without -manifest", label)
		}
	}

	entry := liturgydoc.ManifestEntry{
		Slug:           slug,
		Name:           name,
		NameAm:         nameAm,
		SectionTitle:   sectionTitle,
		SectionTitleAm: sectionTitleAm,
	}

	return convertOne(beginning, anaphoraPath, outputPath, entry, budget)
}

func convertManifest(
	beginning *liturgydoc.Document,
	manifestPath string,
	sourceDir string,
	outputDir string,
	budget liturgydoc.Budget,
) error {
	if strings.TrimSpace(outputDir) == "" {
		return fmt.Errorf("-out-dir is required with -manifest")
	}

	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		return fmt.Errorf("open liturgy manifest: %w", err)
	}
	defer manifestFile.Close()

	entries, err := liturgydoc.DecodeManifest(manifestFile)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(sourceDir, entry.SourceFile)
		targetPath := filepath.Join(outputDir, entry.Slug+"-import.json")

		if err := convertOne(
			beginning,
			sourcePath,
			targetPath,
			entry,
			budget,
		); err != nil {
			return fmt.Errorf("%s: %w", entry.Slug, err)
		}
	}

	slog.Info("liturgy manifest converted", "liturgies", len(entries))

	return nil
}

func convertOne(
	beginning *liturgydoc.Document,
	sourcePath string,
	outputPath string,
	entry liturgydoc.ManifestEntry,
	budget liturgydoc.Budget,
) error {
	anaphora, err := decodeFile(sourcePath)
	if err != nil {
		return fmt.Errorf("anaphora %s: %w", sourcePath, err)
	}

	document, stats, err := liturgydoc.Convert(
		beginning,
		anaphora,
		entry.Options(budget),
	)
	if err != nil {
		return err
	}

	if err := writeJSONAtomically(outputPath, document); err != nil {
		return err
	}

	slog.Info(
		"liturgy converted",
		"slug", entry.Slug,
		"output", outputPath,
		"beginning_entries", stats.BeginningEntries,
		"anaphora_entries", stats.AnaphoraEntries,
		"source_groups", stats.SourceGroups,
		"pages", stats.Pages,
		"verses", stats.Verses,
		"review_segments", stats.ReviewSegments,
		"oversize_pages", stats.OversizePages,
		"largest_page_runes", stats.LargestPageRunes,
	)

	return nil
}

func decodeFile(path string) (*liturgydoc.Document, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open liturgy file: %w", err)
	}
	defer file.Close()

	return liturgydoc.Decode(file)
}

func writeJSONAtomically(path string, value any) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	tempFile, err := os.CreateTemp(directory, ".liturgy-*.json")
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
