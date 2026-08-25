package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/yohannes/kidasie-backend/internal/config"
	"github.com/yohannes/kidasie-backend/internal/contentimport"
	"github.com/yohannes/kidasie-backend/internal/database"
)

func main() {
	if err := run(); err != nil {
		slog.Error("set liturgy audio failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		slug       string
		filePath   string
		publicURL  string
		durationMs int
		confirm    string
	)

	flag.StringVar(&slug, "slug", "", "liturgy slug")
	flag.StringVar(&filePath, "file", "", "path to the complete recording")
	flag.StringVar(&publicURL, "url", "", "public HTTPS audio URL")
	flag.IntVar(&durationMs, "duration-ms", 0, "complete recording duration in milliseconds")
	flag.StringVar(&confirm, "confirm", "", "must exactly match the liturgy slug")
	flag.Parse()

	slug = strings.TrimSpace(slug)
	filePath = strings.TrimSpace(filePath)
	publicURL = strings.TrimSpace(publicURL)

	if slug == "" || filePath == "" || publicURL == "" {
		return fmt.Errorf("-slug, -file, and -url are required")
	}
	if confirm != slug {
		return fmt.Errorf("-confirm must exactly match -slug")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open audio file: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat audio file: %w", err)
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("fingerprint audio file: %w", err)
	}

	mimeType := mime.TypeByExtension(strings.ToLower(filepath.Ext(filePath)))
	if separator := strings.IndexByte(mimeType, ';'); separator >= 0 {
		mimeType = mimeType[:separator]
	}
	if mimeType == "" && strings.EqualFold(filepath.Ext(filePath), ".mp3") {
		mimeType = "audio/mpeg"
	}

	audio := contentimport.Audio{
		URL:        publicURL,
		DurationMs: durationMs,
		SizeBytes:  info.Size(),
		MIMEType:   mimeType,
		SHA256:     fmt.Sprintf("%x", hasher.Sum(nil)),
	}
	if err := audio.Validate(); err != nil {
		return fmt.Errorf("validate audio metadata: %w", err)
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

	updateCtx, cancelUpdate := context.WithTimeout(ctx, 30*time.Second)
	defer cancelUpdate()

	const query = `
		UPDATE liturgies
		SET
			audio_url = $2,
			audio_duration_ms = $3,
			audio_size_bytes = $4,
			audio_mime_type = $5,
			audio_sha256 = $6,
			updated_at = NOW()
		WHERE slug = $1
		RETURNING status, content_version
	`

	var (
		status         string
		contentVersion int64
	)
	if err := pool.QueryRow(
		updateCtx,
		query,
		slug,
		audio.URL,
		audio.DurationMs,
		audio.SizeBytes,
		audio.MIMEType,
		audio.SHA256,
	).Scan(&status, &contentVersion); err != nil {
		return fmt.Errorf("update liturgy audio: %w", err)
	}

	slog.Info(
		"liturgy audio metadata updated",
		"slug",
		slug,
		"status",
		status,
		"content_version",
		contentVersion,
		"duration_ms",
		audio.DurationMs,
		"size_bytes",
		audio.SizeBytes,
		"sha256",
		audio.SHA256,
	)

	return nil
}
