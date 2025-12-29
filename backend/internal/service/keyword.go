package service

import (
	"context"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
)

type KeywordService struct {
	repo *repository.KeywordRepository
}

func NewKeywordService(repo *repository.KeywordRepository) *KeywordService {
	return &KeywordService{repo: repo}
}

func (s *KeywordService) Create(ctx context.Context, req *model.CreateKeywordRequest) (*model.Keyword, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, req)
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
