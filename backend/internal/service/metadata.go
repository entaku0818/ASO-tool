package service

import (
	"context"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/entaku0818/aso-compass/backend/internal/repository"
)

type MetadataService struct {
	repo *repository.MetadataRepository
}

func NewMetadataService(repo *repository.MetadataRepository) *MetadataService {
	return &MetadataService{repo: repo}
}

func (s *MetadataService) List(ctx context.Context, appID string) ([]*model.AppMetadataVersion, error) {
	return s.repo.List(ctx, appID)
}

func (s *MetadataService) Get(ctx context.Context, id string) (*model.AppMetadataVersion, error) {
	return s.repo.Get(ctx, id)
}

func (s *MetadataService) Upsert(ctx context.Context, req *model.UpsertMetadataRequest) (*model.AppMetadataVersion, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return s.repo.Upsert(ctx, req)
}

func (s *MetadataService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
