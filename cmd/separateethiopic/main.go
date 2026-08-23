package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/yohannes/kidasie-backend/internal/ethiopicsplit"
	"github.com/yohannes/kidasie-backend/internal/flatanaphora"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Ethiopic separation failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var filePath string
	var pdfPath string
	var outputPath string
	var reportPath string
	var options ethiopicsplit.Options

	flag.StringVar(&filePath, "file", "", "path to flat anaphora JSON")
	flag.StringVar(&pdfPath, "pdf", "", "path to the matching source PDF")
	flag.StringVar(&outputPath, "out", "", "path for separated JSON")
	flag.StringVar(&reportPath, "report", "", "path for the separation report")
	flag.IntVar(&options.DPI, "dpi", 220, "PDF rendering resolution")
	flag.IntVar(&options.Workers, "workers", 4, "parallel OCR workers")
	flag.BoolVar(
		&options.MetadataOnly,
		"metadata-only",
		false,
		"skip OCR fallback and retain unresolved Ethiopic text",
	)
	flag.Float64Var(
		&options.MinConfidence,
		"min-confidence",
		0.72,
		"minimum automatic split confidence",
	)
	flag.Parse()

	if strings.TrimSpace(filePath) == "" || strings.TrimSpace(pdfPath) == "" {
		return fmt.Errorf("-file and -pdf are required")
	}
	if outputPath == "" {
		outputPath = strings.TrimSuffix(
			filePath,
			filepath.Ext(filePath),
		) + "-separated.json"
	}
	if reportPath == "" {
		reportPath = strings.TrimSuffix(
			outputPath,
			filepath.Ext(outputPath),
		) + "-report.json"
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open input JSON: %w", err)
	}
	entries, err := flatanaphora.Decode(file)
	closeErr := file.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return fmt.Errorf("close input JSON: %w", closeErr)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	startedAt := time.Now()
	options.Progress = func(done int, total int) {
		if done == 1 || done == total || done%10 == 0 {
			slog.Info(
				"checking PDF pages",
				"completed",
				done,
				"total",
				total,
			)
		}
	}

	separated, report, err := ethiopicsplit.Separate(
		ctx,
		pdfPath,
		entries,
		options,
	)
	if err != nil {
		return err
	}
	if err := writeJSON(outputPath, separated); err != nil {
		return err
	}
	if err := writeJSON(reportPath, report); err != nil {
		return err
	}

	slog.Info(
		"Ethiopic separation completed",
		"total",
		report.Total,
		"separated",
		report.Separated,
		"unresolved",
		report.Unresolved,
		"output",
		outputPath,
		"report",
		reportPath,
		"elapsed",
		time.Since(startedAt).Round(time.Second),
	)
	return nil
}

func writeJSON(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		_ = file.Close()
		return fmt.Errorf("encode %s: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close %s: %w", path, err)
	}
	return nil
}
