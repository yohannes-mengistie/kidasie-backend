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
  			id,
  			sort_order,
  			text_geez,
  			text_am,
  			COALESCE(text_en, ''),
  			role,
  			start_ms,
  			end_ms
  		FROM verses
  		WHERE section_id = $1
  		ORDER BY sort_order
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
