package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/yohannes/kidasie-backend/internal/config"
	"github.com/yohannes/kidasie-backend/internal/contentimport"
	"github.com/yohannes/kidasie-backend/internal/database"
)

func main() {
	if err := run(); err != nil {
		slog.Error(
			"content import failed",
			"error",
			err,
		)
		os.Exit(1)
	}
}

func run() error {
	var filePath string

	flag.StringVar(
		&filePath,
		"file",
		"",
		"path to the content JSON file",
	)
	flag.Parse()

	if strings.TrimSpace(filePath) == "" {
		return fmt.Errorf("-file is required")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open content file: %w", err)
	}
	defer file.Close()

	document, err := contentimport.Decode(file)
	if err != nil {
		return err
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

	startupCtx, cancelStartup := context.WithTimeout(
		ctx,
		time.Minute,
	)

	pool, err := database.OpenPostgres(
		startupCtx,
		cfg.DatabaseURL,
	)
	cancelStartup()

	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	importCtx, cancelImport := context.WithTimeout(
		ctx,
		5*time.Minute,
	)
	defer cancelImport()

	importer := contentimport.NewImporter(pool)

	if err := importer.Import(
		importCtx,
		document,
	); err != nil {
		return err
	}

	slog.Info(
		"content imported successfully",
		"slug",
		document.Slug,
		"sections",
		len(document.Sections),
		"status",
		"draft",
	)

	return nil
}
