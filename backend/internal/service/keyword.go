package service

import (
	"context"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
)

const freePlanKeywordLimit = 10

type KeywordService struct {
	repo           *repository.KeywordRepository
	appRepo        *repository.AppRepository
	trackedKeyRepo *repository.TrackedKeywordRepository
	userRepo       *repository.UserRepository
}

func NewKeywordService(repo *repository.KeywordRepository) *KeywordService {
	return &KeywordService{repo: repo}
}

func NewKeywordServiceWithTracking(repo *repository.KeywordRepository, appRepo *repository.AppRepository, trackedKeyRepo *repository.TrackedKeywordRepository) *KeywordService {
	return &KeywordService{
		repo:           repo,
		appRepo:        appRepo,
		trackedKeyRepo: trackedKeyRepo,
	}
}

func NewKeywordServiceWithBilling(repo *repository.KeywordRepository, appRepo *repository.AppRepository, trackedKeyRepo *repository.TrackedKeywordRepository, userRepo *repository.UserRepository) *KeywordService {
	return &KeywordService{
		repo:           repo,
		appRepo:        appRepo,
		trackedKeyRepo: trackedKeyRepo,
		userRepo:       userRepo,
	}
}

func (s *KeywordService) Create(ctx context.Context, userID string, req *model.CreateKeywordRequest) (*model.Keyword, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Enforce free plan keyword limit
	if s.userRepo != nil && userID != "" {
		user, err := s.userRepo.GetByID(ctx, userID)
		if err == nil && !user.IsPro() {
			count, err := s.repo.CountByApp(ctx, req.AppID)
			if err != nil {
				return nil, err
			}
			if count >= freePlanKeywordLimit {
				return nil, model.ErrPlanLimitExceeded
			}
		}
	}

	keyword, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	// Auto-create tracked keyword if tracking repo is available
	if s.trackedKeyRepo != nil && s.appRepo != nil {
		app, err := s.appRepo.GetByID(ctx, req.AppID)
		if err == nil {
			_, _ = s.trackedKeyRepo.Create(ctx, app.UserID, req.Keyword, req.Country, string(app.Platform))
		}
	}

	return keyword, nil
}

func (s *KeywordService) Get(ctx context.Context, id string) (*model.Keyword, error) {
	return s.repo.Get(ctx, id)
}

func (s *KeywordService) ListByApp(ctx context.Context, appID string) ([]*model.Keyword, error) {
	return s.repo.ListByApp(ctx, appID)
}

func (s *KeywordService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
