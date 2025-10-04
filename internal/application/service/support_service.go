package service

import (
	"context"

	"telecomx-support-service/internal/domain/model"
	"telecomx-support-service/internal/infrastructure/adapter/repository"
)

type SupportService struct {
	repo *repository.MongoRepository
}

func NewSupportService(repo *repository.MongoRepository) *SupportService {
	return &SupportService{repo: repo}
}

func (s *SupportService) Create(ctx context.Context, p *model.Support) error {
	return s.repo.Create(ctx, p)
}

func (s *SupportService) UpdateStatus(ctx context.Context, userID, status string) error {
	return s.repo.UpdateStatus(ctx, userID, status)
}

func (s *SupportService) Delete(ctx context.Context, userID string) error {
	return s.repo.DeleteByUserID(ctx, userID)
}

func (s *SupportService) GetAll(ctx context.Context) ([]model.Support, error) {
	return s.repo.GetAll(ctx)
}
