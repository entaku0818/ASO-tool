package service

import (
	"context"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
)

type KeywordService struct {
	repo            *repository.KeywordRepository
	appRepo         *repository.AppRepository
	trackedKeyRepo  *repository.TrackedKeywordRepository
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

func (s *KeywordService) Create(ctx context.Context, req *model.CreateKeywordRequest) (*model.Keyword, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	keyword, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	// Auto-create tracked keyword if tracking repo is available
	if s.trackedKeyRepo != nil && s.appRepo != nil {
		app, err := s.appRepo.GetByID(ctx, req.AppID)
		if err == nil {
			s.trackedKeyRepo.Create(ctx, app.UserID, req.Keyword, req.Country, string(app.Platform))
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
