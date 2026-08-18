package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yohannes/kidasie-backend/internal/domain"
	"github.com/yohannes/kidasie-backend/internal/service"
)

type ContentRepo struct {
	liturgies *LiturgyRepo
	sections  *SectionRepo
	verses    *VerseRepo
}

var _ service.ContentRepository = (*ContentRepo)(nil)

func NewContentRepo(pool *pgxpool.Pool) *ContentRepo {
	return &ContentRepo{
		liturgies: NewLiturgyRepo(pool),
		sections:  NewSectionRepo(pool),
		verses:    NewVerseRepo(pool),
	}
}

func (r *ContentRepo) GetLiturgyContentBySlug(
	ctx context.Context,
	slug string,
) (*domain.LiturgyContent, error) {
	liturgy, err := r.liturgies.GetLiturgyBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("get liturgy: %w", err)
	}

	sections, err := r.sections.ListSectionsByLiturgySlug(
		ctx,
		slug,
	)
	if err != nil {
		return nil, fmt.Errorf("list sections: %w", err)
	}

	for i := range sections {
		verses, err := r.verses.ListVersesBySectionID(
			ctx,
			sections[i].ID,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"list verses for section %d: %w",
				sections[i].ID,
				err,
			)
		}

		sections[i].Verses = verses
	}

	return &domain.LiturgyContent{
		Liturgy:  *liturgy,
		Sections: sections,
	}, nil
}
