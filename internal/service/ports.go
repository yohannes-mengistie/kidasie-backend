package service

import (
	"context"
	"github.com/yohannes/kidasie-backend/internal/domain"
)

type LiturgyRepository interface  {
	ListLiturgies(ctx context.Context) ([]domain.Liturgy, error)
	GetSection(ctx context.Context, id int64) (*domain.Section, error)
}