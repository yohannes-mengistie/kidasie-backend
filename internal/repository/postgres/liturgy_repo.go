package postgres

import (
	"context"
	"fmt"
	"errors"
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
			SELECT id,slug,name, name_am, content_version
			FROM liturgies
			ORDER BY id
  	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query liturgies: %w", err)
	}
	defer rows.Close()

	liturgies := make([]domain.Liturgy, 0)

	for rows.Next() {
		var liturgy domain.Liturgy

		if err := rows.Scan(
			&liturgy.ID,
			&liturgy.Slug,
			&liturgy.Name,
			&liturgy.NameAm,
			&liturgy.ContentVersion,
		); err != nil {
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
			SELECT id, slug, name, name_am, content_version
			FROM liturgies
			WHERE slug = $1
  	`

  	var liturgy domain.Liturgy

  	err := r.pool.QueryRow(ctx, query, slug).Scan(
  		&liturgy.ID,
  		&liturgy.Slug,
  		&liturgy.Name,
  		&liturgy.NameAm,
		&liturgy.ContentVersion,
  	)
  	if errors.Is(err, pgx.ErrNoRows) {
  		return nil, domain.ErrNotFound
  	}
  	if err != nil {
  		return nil, fmt.Errorf("query liturgy by slug: %w", err)
  	}

  	return &liturgy, nil
  }

