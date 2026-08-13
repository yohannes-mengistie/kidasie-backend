package memory

import (
	"context"
	"github.com/yohannes/kidasie-backend/internal/domain"
)

type LiturgyRepo struct {
	liturgies []domain.Liturgy
	sections  map[int64]*domain.Section
}

func NewLiturgyRepo(liturgies []domain.Liturgy, sections []domain.Section) *LiturgyRepo {
	sectionMap := make(map[int64]*domain.Section, len(sections))

	for i := range sections {
		section := sections[i]
		sectionMap[section.ID] = &section
	}

	return &LiturgyRepo{
		liturgies: liturgies,
		sections:  sectionMap,
	}
}

func (r *LiturgyRepo) ListLiturgies(
	ctx context.Context,
) ([]domain.Liturgy, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	liturgies := make([]domain.Liturgy, len(r.liturgies))
	copy(liturgies, r.liturgies)

	return liturgies, nil

}

func (r *LiturgyRepo) GetSection(
  	ctx context.Context,
  	id int64,
  ) (*domain.Section, error) {
  	if err := ctx.Err(); err != nil {
  		return nil, err
  	}

  	section, exists := r.sections[id]
  	if !exists {
  		return nil, domain.ErrNotFound
  	}

  	result := *section
  	result.Verses = make([]domain.Verse, len(section.Verses))
  	copy(result.Verses, section.Verses)

  	return &result, nil
  }

