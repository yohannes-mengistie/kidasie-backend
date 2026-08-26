package service

import (
	"context"
	"fmt"

	"github.com/yohannes/kidasie-backend/internal/domain"
)

type AnnouncementService struct {
	repository AnnouncementRepository
}

func NewAnnouncementService(repository AnnouncementRepository) *AnnouncementService {
	return &AnnouncementService{repository: repository}
}

func (service *AnnouncementService) ListPublishedAnnouncements(
	ctx context.Context,
) ([]domain.Announcement, error) {
	return service.repository.ListPublishedAnnouncements(ctx)
}

type AnnouncementAdminService struct {
	repository AnnouncementAdminRepository
}

func NewAnnouncementAdminService(
	repository AnnouncementAdminRepository,
) *AnnouncementAdminService {
	return &AnnouncementAdminService{repository: repository}
}

func (service *AnnouncementAdminService) UpsertAnnouncement(
	ctx context.Context,
	input domain.AnnouncementInput,
) (*domain.Announcement, error) {
	if err := input.NormalizeAndValidate(); err != nil {
		return nil, fmt.Errorf("validate announcement: %w", err)
	}
	return service.repository.UpsertAnnouncement(ctx, input)
}
