package service

import (
	"context"

	"github.com/yohannes/kidasie-backend/internal/domain"
)

type VerseService struct {
	repository VerseRepository
}

func NewVerseService(repository VerseRepository) *VerseService {
	return &VerseService{
		repository: repository,
	}
}

func (s *VerseService) ListVersesBySectionID(
	ctx context.Context,
	sectionID int64,
) ([]domain.Verse, error) {
	return s.repository.ListVersesBySectionID(ctx, sectionID)
}
