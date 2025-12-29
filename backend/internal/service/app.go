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

func (s *AppService) Create(ctx context.Context, req *model.CreateAppRequest) (*model.App, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, req)
}

func (s *AppService) Get(ctx context.Context, id string) (*model.App, error) {
	return s.repo.Get(ctx, id)
}

func (s *AppService) List(ctx context.Context) ([]*model.App, error) {
	return s.repo.List(ctx)
}

func (s *AppService) Update(ctx context.Context, id string, req *model.UpdateAppRequest) (*model.App, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *AppService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
