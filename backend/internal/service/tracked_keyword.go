package service

import (
	"context"

	"github.com/entaku0818/aso-tool/backend/internal/repository"
)

type TrackedKeywordService struct {
	repo *repository.TrackedKeywordRepository
}

func NewTrackedKeywordService(repo *repository.TrackedKeywordRepository) *TrackedKeywordService {
	return &TrackedKeywordService{repo: repo}
}

func (s *TrackedKeywordService) List(ctx context.Context) ([]*repository.TrackedKeyword, error) {
	return s.repo.List(ctx)
}

func (s *TrackedKeywordService) Get(ctx context.Context, id string) (*repository.TrackedKeyword, error) {
	return s.repo.Get(ctx, id)
}

func (s *TrackedKeywordService) GetSearchResults(ctx context.Context, id string) ([]*repository.SearchResult, error) {
	return s.repo.GetLatestSearchResults(ctx, id)
}

func (s *TrackedKeywordService) SaveSearchResults(ctx context.Context, id string, results []repository.SearchResult) error {
	return s.repo.SaveSearchResults(ctx, id, results)
}
