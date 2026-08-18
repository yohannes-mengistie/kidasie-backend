package service

import (
	"context"
	"github.com/yohannes/kidasie-backend/internal/domain"
)

type LiturgyRepository interface  {
	ListLiturgies(ctx context.Context) ([]domain.Liturgy, error)
	GetLiturgyBySlug(ctx context.Context , slug string)(*domain.Liturgy,error)
}

type SectionRepository interface {
  	ListSectionsByLiturgySlug(
  		ctx context.Context,
  		slug string,
  	) ([]domain.Section, error)
  }

type VerseRepository interface {
  	ListVersesBySectionID(
  		ctx context.Context,
  		sectionID int64,
  	) ([]domain.Verse, error)
  }


