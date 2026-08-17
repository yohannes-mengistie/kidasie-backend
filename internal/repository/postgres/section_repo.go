package postgres

  import (
  	"context"
  	"fmt"

  	"github.com/jackc/pgx/v5/pgxpool"

  	"github.com/yohannes/kidasie-backend/internal/domain"
  	"github.com/yohannes/kidasie-backend/internal/service"
  )

  type SectionRepo struct {
  	pool *pgxpool.Pool
  }

  var _ service.SectionRepository = (*SectionRepo)(nil)

  func NewSectionRepo(pool *pgxpool.Pool) *SectionRepo {
  	return &SectionRepo{
  		pool: pool,
  	}
  }

  func (r *SectionRepo) ListSectionsByLiturgySlug(
  	ctx context.Context,
  	slug string,
  ) ([]domain.Section, error) {
  	const query = `
  		SELECT
  			s.id,
  			s.liturgy_id,
  			s.sort_order,
  			s.title,
  			s.title_am,
  			COALESCE(s.audio_url, '')
  		FROM sections AS s
  		INNER JOIN liturgies AS l
  			ON l.id = s.liturgy_id
  		WHERE l.slug = $1
  		ORDER BY s.sort_order
  	`

  	rows, err := r.pool.Query(ctx, query, slug)
  	if err != nil {
  		return nil, fmt.Errorf("query sections: %w", err)
  	}
  	defer rows.Close()

  	sections := make([]domain.Section, 0)

  	for rows.Next() {
  		var section domain.Section

  		if err := rows.Scan(
  			&section.ID,
  			&section.LiturgyID,
  			&section.Order,
  			&section.Title,
  			&section.TitleAm,
  			&section.AudioURL,
  		); err != nil {
  			return nil, fmt.Errorf("scan section: %w", err)
  		}

  		sections = append(sections, section)
  	}

  	if err := rows.Err(); err != nil {
  		return nil, fmt.Errorf("iterate sections: %w", err)
  	}

  	return sections, nil
  }
