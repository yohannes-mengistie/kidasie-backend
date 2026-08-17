package service

import (
	"context"
	"github.com/yohannes/kidasie-backend/internal/domain"
)

type LiturgyService struct{
	repository LiturgyRepository
}

func NewLiturgyService(repository LiturgyRepository) *LiturgyService {
	return &LiturgyService{
		repository:repository,
	}
}

func (s *LiturgyService) ListLiturgies(ctx context.Context) ([]domain.Liturgy, error) {
	return s.repository.ListLiturgies(ctx)
}