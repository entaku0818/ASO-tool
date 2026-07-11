package service

import (
	"context"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/plan"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
)

type AppService struct {
	repo     *repository.AppRepository
	userRepo *repository.UserRepository
}

func NewAppService(repo *repository.AppRepository) *AppService {
	return &AppService{repo: repo}
}

func NewAppServiceWithBilling(repo *repository.AppRepository, userRepo *repository.UserRepository) *AppService {
	return &AppService{repo: repo, userRepo: userRepo}
}

func (s *AppService) Create(ctx context.Context, userID string, req *model.CreateAppRequest) (*model.App, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Enforce free plan app limit atomically to prevent TOCTOU races.
	if s.userRepo != nil && userID != "" {
		user, err := s.userRepo.GetByID(ctx, userID)
		if err == nil && !user.IsPro() {
			return s.repo.CreateWithLimit(ctx, userID, req, plan.FreeAppLimit)
		}
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
