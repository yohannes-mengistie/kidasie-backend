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
	"syscall"
	"time"

	"github.com/yohannes/kidasie-backend/internal/slideextract"
)

func main() {
	if err := run(); err != nil {
		slog.Error("slide extraction failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var options slideextract.Options
	var outputPath string

	flag.StringVar(
		&options.PDFPath,
		"pdf",
		"source-material/liturgy.pdf",
		"path to the permitted liturgy PDF",
	)
	flag.StringVar(
		&options.AudioPath,
		"audio",
		"",
		"optional path to the matching audio file",
	)
	flag.StringVar(
		&outputPath,
		"out",
		"content/generated/slides.json",
		"path for the reviewable draft JSON",
	)
	flag.IntVar(&options.StartPage, "start", 1, "first PDF page")
	flag.IntVar(&options.EndPage, "end", 0, "last PDF page; 0 means all")
	flag.IntVar(&options.DPI, "dpi", 180, "OCR rendering resolution")
	flag.IntVar(&options.Workers, "workers", 2, "parallel OCR workers")
	flag.Parse()

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
				"extracting slides",
				"completed",
				done,
				"total",
				total,
			)
		}
	}

	draft, err := slideextract.Extract(ctx, options)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(
		filepath.Dir(outputPath),
		0o755,
	); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}

	encoder := json.NewEncoder(outputFile)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(draft); err != nil {
		_ = outputFile.Close()
		return fmt.Errorf("encode slide draft: %w", err)
	}

	if err := outputFile.Close(); err != nil {
		return fmt.Errorf("close output file: %w", err)
	}

	slog.Info(
		"slide extraction completed",
		"pages",
		len(draft.Pages),
		"output",
		outputPath,
		"elapsed",
		time.Since(startedAt).Round(time.Second),
	)

	return nil
}
