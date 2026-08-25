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

	"github.com/yohannes/kidasie-backend/internal/config"
	"github.com/yohannes/kidasie-backend/internal/database"
	"github.com/yohannes/kidasie-backend/internal/repository/postgres"
)

func main() {
	if err := run(); err != nil {
		slog.Error("offline export failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var outputDirectory string
	flag.StringVar(
		&outputDirectory,
		"out",
		"",
		"Flutter offline asset directory",
	)
	flag.Parse()

	outputDirectory = strings.TrimSpace(outputDirectory)
	if outputDirectory == "" {
		return fmt.Errorf("-out is required")
	}

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	connectCtx, cancelConnect := context.WithTimeout(ctx, time.Minute)
	pool, err := database.OpenPostgres(connectCtx, cfg.DatabaseURL)
	cancelConnect()
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	liturgies, err := postgres.NewLiturgyRepo(pool).ListLiturgies(ctx)
	if err != nil {
		return err
	}

	if err := writeJSON(
		filepath.Join(outputDirectory, "liturgies.json"),
		map[string]any{"data": liturgies},
	); err != nil {
		return err
	}

	contentRepository := postgres.NewContentRepo(pool)
	for _, liturgy := range liturgies {
		content, err := contentRepository.GetLiturgyContentBySlug(
			ctx,
			liturgy.Slug,
		)
		if err != nil {
			return fmt.Errorf("export %s: %w", liturgy.Slug, err)
		}

		if err := writeJSON(
			filepath.Join(outputDirectory, liturgy.Slug+".json"),
			map[string]any{"data": content},
		); err != nil {
			return err
		}
	}

	slog.Info(
		"offline assets exported",
		"directory",
		outputDirectory,
		"liturgies",
		len(liturgies),
	)
	return nil
}

func writeJSON(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	temporary, err := os.CreateTemp(filepath.Dir(path), ".offline-*.json")
	if err != nil {
		return fmt.Errorf("create temporary export: %w", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	encoder := json.NewEncoder(temporary)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		temporary.Close()
		return fmt.Errorf("encode %s: %w", path, err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close %s: %w", path, err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("replace %s: %w", path, err)
	}
	return nil
}
