package service

import (
	"context"
	"github.com/yohannes/kidasie-backend/internal/domain"
)

type LiturgyRepository interface  {
	ListLiturgies(ctx context.Context) ([]domain.Liturgy, error)
}

type SectionRepository interface {
	GetSection(ctx context.Context, id int64) (*domain.Section, error)
}

