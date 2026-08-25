package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yohannes/kidasie-backend/internal/domain"
	"github.com/yohannes/kidasie-backend/internal/service"
)

type LiturgyRepo struct {
	pool *pgxpool.Pool
}

var _ service.LiturgyRepository = (*LiturgyRepo)(nil)

func NewLiturgyRepo(pool *pgxpool.Pool) *LiturgyRepo {
	return &LiturgyRepo{
		pool: pool,
	}
}

func (r *LiturgyRepo) ListLiturgies(
	ctx context.Context,
) ([]domain.Liturgy, error) {
	const query = `
		SELECT
			l.id,
			l.slug,
			l.name,
			l.name_am,
			l.content_version,
			l.audio_url,
			l.audio_duration_ms,
			l.audio_size_bytes,
			l.audio_mime_type,
			l.audio_sha256,
			EXISTS (
				SELECT 1
				FROM sections AS s
				INNER JOIN verses AS v ON v.section_id = s.id
				WHERE s.liturgy_id = l.id
			) AS has_content
		FROM liturgies AS l
		WHERE l.status = 'published'
		ORDER BY l.id
  	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query liturgies: %w", err)
	}
	defer rows.Close()

	liturgies := make([]domain.Liturgy, 0)

	for rows.Next() {
		var liturgy domain.Liturgy

		if err := scanLiturgy(rows, &liturgy); err != nil {
			return nil, fmt.Errorf("scan liturgy: %w", err)
		}

		liturgies = append(liturgies, liturgy)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate liturgies: %w", err)
	}

	return liturgies, nil
}

func (r *LiturgyRepo) GetLiturgyBySlug(
	ctx context.Context,
	slug string,
) (*domain.Liturgy, error) {
	const query = `
		SELECT
			l.id,
			l.slug,
			l.name,
			l.name_am,
			l.content_version,
			l.audio_url,
			l.audio_duration_ms,
			l.audio_size_bytes,
			l.audio_mime_type,
			l.audio_sha256,
			EXISTS (
				SELECT 1
				FROM sections AS s
				INNER JOIN verses AS v ON v.section_id = s.id
				WHERE s.liturgy_id = l.id
			) AS has_content
		FROM liturgies AS l
		WHERE l.slug = $1
			AND l.status = 'published'
  	`

	var liturgy domain.Liturgy

	err := scanLiturgy(r.pool.QueryRow(ctx, query, slug), &liturgy)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query liturgy by slug: %w", err)
	}

	return &liturgy, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanLiturgy(row rowScanner, liturgy *domain.Liturgy) error {
	var (
		audioURL        *string
		audioDurationMs *int
		audioSizeBytes  *int64
		audioMIMEType   *string
		audioSHA256     *string
	)

	if err := row.Scan(
		&liturgy.ID,
		&liturgy.Slug,
		&liturgy.Name,
		&liturgy.NameAm,
		&liturgy.ContentVersion,
		&audioURL,
		&audioDurationMs,
		&audioSizeBytes,
		&audioMIMEType,
		&audioSHA256,
		&liturgy.HasContent,
	); err != nil {
		return err
	}

	if audioURL != nil {
		liturgy.Audio = &domain.Audio{
			URL:        *audioURL,
			DurationMs: valueOrZero(audioDurationMs),
			SizeBytes:  valueOrZero(audioSizeBytes),
			MIMEType:   valueOrEmpty(audioMIMEType),
			SHA256:     valueOrEmpty(audioSHA256),
		}
	}

	return nil
}

func valueOrZero[T int | int64](value *T) T {
	if value == nil {
		return 0
	}
	return *value
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
