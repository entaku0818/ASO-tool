package service

import (
	"context"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/entaku0818/aso-compass/backend/internal/repository"
)

type ReviewService struct {
	repo *repository.ReviewRepository
}

func NewReviewService(repo *repository.ReviewRepository) *ReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) Create(ctx context.Context, req *model.CreateReviewRequest) (*model.Review, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, req)
}

func (s *ReviewService) Get(ctx context.Context, id string) (*model.Review, error) {
	return s.repo.Get(ctx, id)
}

func (s *ReviewService) ListByApp(ctx context.Context, appID string, limit, offset int) ([]*model.Review, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.ListByApp(ctx, appID, limit, offset)
}

func (s *ReviewService) GetStats(ctx context.Context, appID string) (*model.ReviewStats, error) {
	return s.repo.GetStats(ctx, appID)
}

func (s *ReviewService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
