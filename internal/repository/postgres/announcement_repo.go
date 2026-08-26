package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yohannes/kidasie-backend/internal/domain"
	"github.com/yohannes/kidasie-backend/internal/service"
)

type AnnouncementRepo struct {
	pool *pgxpool.Pool
}

var _ service.AnnouncementRepository = (*AnnouncementRepo)(nil)
var _ service.AnnouncementAdminRepository = (*AnnouncementRepo)(nil)

func NewAnnouncementRepo(pool *pgxpool.Pool) *AnnouncementRepo {
	return &AnnouncementRepo{pool: pool}
}

func (repository *AnnouncementRepo) ListPublishedAnnouncements(
	ctx context.Context,
) ([]domain.Announcement, error) {
	const query = `
		SELECT
			id,
			slug,
			version,
			title_am,
			title_en,
			body_am,
			body_en,
			kind,
			action_type,
			COALESCE(action_value, ''),
			priority,
			is_pinned,
			status,
			published_at,
			expires_at,
			updated_at
		FROM announcements
		WHERE status = 'published'
			AND published_at <= NOW()
			AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY is_pinned DESC, priority DESC, published_at DESC, id DESC
	`

	rows, err := repository.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query published announcements: %w", err)
	}
	defer rows.Close()

	announcements := make([]domain.Announcement, 0)
	for rows.Next() {
		var announcement domain.Announcement
		if err := scanAnnouncement(rows, &announcement); err != nil {
			return nil, fmt.Errorf("scan announcement: %w", err)
		}
		announcements = append(announcements, announcement)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate announcements: %w", err)
	}
	return announcements, nil
}

func (repository *AnnouncementRepo) UpsertAnnouncement(
	ctx context.Context,
	input domain.AnnouncementInput,
) (*domain.Announcement, error) {
	const query = `
		INSERT INTO announcements (
			slug,
			title_am,
			title_en,
			body_am,
			body_en,
			kind,
			action_type,
			action_value,
			priority,
			is_pinned,
			status,
			published_at,
			expires_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
			CASE WHEN $11 = 'published' THEN NOW() ELSE NULL END,
			$12
		)
		ON CONFLICT (slug) DO UPDATE
		SET
			title_am = EXCLUDED.title_am,
			title_en = EXCLUDED.title_en,
			body_am = EXCLUDED.body_am,
			body_en = EXCLUDED.body_en,
			kind = EXCLUDED.kind,
			action_type = EXCLUDED.action_type,
			action_value = EXCLUDED.action_value,
			priority = EXCLUDED.priority,
			is_pinned = EXCLUDED.is_pinned,
			status = EXCLUDED.status,
			published_at = CASE
				WHEN EXCLUDED.status = 'published'
					AND announcements.status <> 'published' THEN NOW()
				WHEN EXCLUDED.status = 'published' THEN announcements.published_at
				ELSE announcements.published_at
			END,
			expires_at = EXCLUDED.expires_at
		RETURNING
			id,
			slug,
			version,
			title_am,
			title_en,
			body_am,
			body_en,
			kind,
			action_type,
			COALESCE(action_value, ''),
			priority,
			is_pinned,
			status,
			published_at,
			expires_at,
			updated_at
	`

	var actionValue any
	if input.ActionValue != "" {
		actionValue = input.ActionValue
	}

	var announcement domain.Announcement
	if err := scanAnnouncement(
		repository.pool.QueryRow(
			ctx,
			query,
			input.Slug,
			input.TitleAm,
			input.TitleEn,
			input.BodyAm,
			input.BodyEn,
			input.Kind,
			input.ActionType,
			actionValue,
			input.Priority,
			input.IsPinned,
			input.Status,
			input.ExpiresAt,
		),
		&announcement,
	); err != nil {
		return nil, fmt.Errorf("upsert announcement: %w", err)
	}

	return &announcement, nil
}

func scanAnnouncement(row rowScanner, announcement *domain.Announcement) error {
	return row.Scan(
		&announcement.ID,
		&announcement.Slug,
		&announcement.Version,
		&announcement.TitleAm,
		&announcement.TitleEn,
		&announcement.BodyAm,
		&announcement.BodyEn,
		&announcement.Kind,
		&announcement.Action.Type,
		&announcement.Action.Value,
		&announcement.Priority,
		&announcement.IsPinned,
		&announcement.Status,
		&announcement.PublishedAt,
		&announcement.ExpiresAt,
		&announcement.UpdatedAt,
	)
}
