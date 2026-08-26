package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/yohannes/kidasie-backend/internal/config"
	"github.com/yohannes/kidasie-backend/internal/database"
	"github.com/yohannes/kidasie-backend/internal/domain"
	"github.com/yohannes/kidasie-backend/internal/repository/postgres"
	"github.com/yohannes/kidasie-backend/internal/service"
)

func main() {
	if err := run(); err != nil {
		slog.Error("announcement upsert failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var filePath string
	var confirmation string
	flag.StringVar(&filePath, "file", "", "path to announcement JSON")
	flag.StringVar(&confirmation, "confirm", "", "must exactly match the announcement slug")
	flag.Parse()

	if strings.TrimSpace(filePath) == "" {
		return fmt.Errorf("-file is required")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open announcement file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	var input domain.AnnouncementInput
	if err := decoder.Decode(&input); err != nil {
		return fmt.Errorf("decode announcement: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("announcement file must contain one JSON object")
	}

	if input.Kind == "" {
		input.Kind = "general"
	}
	if input.ActionType == "" {
		input.ActionType = "none"
	}
	if input.Status == "" {
		input.Status = domain.AnnouncementStatusDraft
	}
	if strings.TrimSpace(confirmation) != strings.TrimSpace(input.Slug) {
		return fmt.Errorf("-confirm must exactly match the announcement slug")
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

	adminService := service.NewAnnouncementAdminService(
		postgres.NewAnnouncementRepo(pool),
	)
	upsertCtx, cancelUpsert := context.WithTimeout(ctx, 30*time.Second)
	defer cancelUpsert()

	announcement, err := adminService.UpsertAnnouncement(upsertCtx, input)
	if err != nil {
		return err
	}

	slog.Info(
		"announcement saved",
		"slug", announcement.Slug,
		"version", announcement.Version,
		"status", announcement.Status,
		"published_at", announcement.PublishedAt,
	)
	return nil
}
