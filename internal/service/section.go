package service

import (
	"context"
	"github.com/yohannes/kidasie-backend/internal/domain"
)

type SectionService struct {
	repository SectionRepository
}

func NewSectionService(repository SectionRepository) *SectionService {
	return &SectionService{
		repository: repository,
	}
}

func (s *SectionService) ListSectionsByLiturgySlug(ctx context.Context, slug string) ([]domain.Section, error) {
	return s.repository.ListSectionsByLiturgySlug(ctx, slug)
}
