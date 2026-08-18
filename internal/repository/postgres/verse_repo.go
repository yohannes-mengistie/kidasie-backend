package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yohannes/kidasie-backend/internal/domain"
	"github.com/yohannes/kidasie-backend/internal/service"
)

type VerseRepo struct {
	pool *pgxpool.Pool
}

var _ service.VerseRepository = (*VerseRepo)(nil)

func NewVerseRepo(pool *pgxpool.Pool) *VerseRepo {
	return &VerseRepo{
		pool: pool,
	}
}

func (r *VerseRepo) ListVersesBySectionID(
	ctx context.Context,
	sectionID int64,
) ([]domain.Verse, error) {
	const query = `
		SELECT
			v.id,
			v.sort_order,
			v.text_geez,
			v.text_am,
			COALESCE(v.text_en, ''),
			v.role,
			v.start_ms,
			v.end_ms
		FROM verses AS v
		INNER JOIN sections AS s
			ON s.id = v.section_id
		INNER JOIN liturgies AS l
			ON l.id = s.liturgy_id
		WHERE v.section_id = $1
			AND l.status = 'published'
		ORDER BY v.sort_order
	`

	rows, err := r.pool.Query(ctx, query, sectionID)
	if err != nil {
		return nil, fmt.Errorf("query verses: %w", err)
	}
	defer rows.Close()

	verses := make([]domain.Verse, 0)

	for rows.Next() {
		var verse domain.Verse

		if err := rows.Scan(
			&verse.ID,
			&verse.Order,
			&verse.TextGeez,
			&verse.TextAm,
			&verse.TextEn,
			&verse.Role,
			&verse.StartMs,
			&verse.EndMs,
		); err != nil {
			return nil, fmt.Errorf("scan verse: %w", err)
		}

		verses = append(verses, verse)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate verses: %w", err)
	}

	return verses, nil
}
