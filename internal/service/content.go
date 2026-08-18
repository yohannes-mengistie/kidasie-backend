package service

import (
	"context"
	"github.com/yohannes/kidasie-backend/internal/domain"
)


type ContentService struct{
	repository ContentRepository
}

func NewContentService(repository ContentRepository) *ContentService{
	return &ContentService{
		repository:repository,
	}
}

func (s *ContentService) GetLiturgyContentBySlug(ctx context.Context, slug string)(*domain.LiturgyContent,error){
	return s.repository.GetLiturgyContentBySlug(ctx,slug)
}