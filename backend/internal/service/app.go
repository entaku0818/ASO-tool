package service

import (
	"context"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
)

type AppService struct {
	repo *repository.AppRepository
}

func NewAppService(repo *repository.AppRepository) *AppService {
	return &AppService{repo: repo}
}

func (s *AppService) Create(ctx context.Context, userID string, req *model.CreateAppRequest) (*model.App, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, userID, req)
}

func (s *AppService) Get(ctx context.Context, userID, id string) (*model.App, error) {
	return s.repo.Get(ctx, userID, id)
}

func (s *AppService) List(ctx context.Context, userID string) ([]*model.App, error) {
	return s.repo.List(ctx, userID)
}

func (s *AppService) Update(ctx context.Context, userID, id string, req *model.UpdateAppRequest) (*model.App, error) {
	return s.repo.Update(ctx, userID, id, req)
}

func (s *AppService) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, userID, id)
}

func (s *AppService) FindByName(ctx context.Context, userID, name string) ([]*model.App, error) {
	return s.repo.FindByName(ctx, userID, name)
}
