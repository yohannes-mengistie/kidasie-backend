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

	"github.com/jackc/pgx/v5"

	"github.com/yohannes/kidasie-backend/internal/config"
	"github.com/yohannes/kidasie-backend/internal/database"
)

type publicationSummary struct {
	ID             int64
	Status         string
	Sections       int
	Verses         int
	ReviewRequired int
}

func main() {
	if err := run(); err != nil {
		slog.Error("content publication failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		slug                string
		confirmation        string
		allowReviewRequired bool
	)

	flag.StringVar(&slug, "slug", "", "liturgy slug to publish")
	flag.StringVar(
		&confirmation,
		"confirm",
		"",
		"must exactly match the liturgy slug",
	)
	flag.BoolVar(
		&allowReviewRequired,
		"allow-review-required",
		false,
		"publish even when imported source segments require review",
	)
	flag.Parse()

	slug = strings.TrimSpace(slug)
	if slug == "" {
		return fmt.Errorf("-slug is required")
	}
	if confirmation != slug {
		return fmt.Errorf("-confirm must exactly match -slug")
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

	publishCtx, cancelPublish := context.WithTimeout(ctx, 30*time.Second)
	defer cancelPublish()

	summary, err := publish(
		publishCtx,
		pool,
		slug,
		allowReviewRequired,
	)
	if err != nil {
		return err
	}

	slog.Info(
		"content published",
		"slug",
		slug,
		"sections",
		summary.Sections,
		"verses",
		summary.Verses,
		"review_required",
		summary.ReviewRequired,
	)

	return nil
}

type databasePool interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}

func publish(
	ctx context.Context,
	pool databasePool,
	slug string,
	allowReviewRequired bool,
) (publicationSummary, error) {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return publicationSummary{}, fmt.Errorf(
			"begin publication transaction: %w",
			err,
		)
	}
	defer tx.Rollback(ctx)

	const lockQuery = `
		SELECT id, status
		FROM liturgies
		WHERE slug = $1
		FOR UPDATE
	`

	var summary publicationSummary
	err = tx.QueryRow(ctx, lockQuery, slug).Scan(
		&summary.ID,
		&summary.Status,
	)
	if err == pgx.ErrNoRows {
		return publicationSummary{}, fmt.Errorf(
			"liturgy %q does not exist",
			slug,
		)
	}
	if err != nil {
		return publicationSummary{}, fmt.Errorf(
			"lock liturgy for publication: %w",
			err,
		)
	}

	const summaryQuery = `
		SELECT
			COUNT(DISTINCT s.id),
			COUNT(v.id),
			COUNT(v.id) FILTER (WHERE v.source_needs_review)
		FROM sections AS s
		LEFT JOIN verses AS v
			ON v.section_id = s.id
		WHERE s.liturgy_id = $1
	`
	if err := tx.QueryRow(ctx, summaryQuery, summary.ID).Scan(
		&summary.Sections,
		&summary.Verses,
		&summary.ReviewRequired,
	); err != nil {
		return publicationSummary{}, fmt.Errorf(
			"load publication summary: %w",
			err,
		)
	}

	if summary.Sections == 0 || summary.Verses == 0 {
		return publicationSummary{}, fmt.Errorf(
			"liturgy %q has no importable content",
			slug,
		)
	}

	if summary.ReviewRequired > 0 && !allowReviewRequired {
		return publicationSummary{}, fmt.Errorf(
			"%d verse segments still require review; "+
				"use -allow-review-required only for an explicit preview publication",
			summary.ReviewRequired,
		)
	}

	const invalidAudioTimingQuery = `
		SELECT COUNT(*)
		FROM verses AS v
		INNER JOIN sections AS s
			ON s.id = v.section_id
		WHERE s.liturgy_id = $1
			AND s.audio_url IS NOT NULL
			AND (v.start_ms IS NULL OR v.end_ms IS NULL)
	`

	var invalidAudioTimings int
	if err := tx.QueryRow(
		ctx,
		invalidAudioTimingQuery,
		summary.ID,
	).Scan(&invalidAudioTimings); err != nil {
		return publicationSummary{}, fmt.Errorf(
			"validate audio timing: %w",
			err,
		)
	}
	if invalidAudioTimings > 0 {
		return publicationSummary{}, fmt.Errorf(
			"%d verses are missing timing for an audio section",
			invalidAudioTimings,
		)
	}

	const publishQuery = `
		UPDATE liturgies
		SET
			status = 'published',
			published_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
	`
	if _, err := tx.Exec(ctx, publishQuery, summary.ID); err != nil {
		return publicationSummary{}, fmt.Errorf("publish liturgy: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return publicationSummary{}, fmt.Errorf(
			"commit publication: %w",
			err,
		)
	}

	return summary, nil
}
