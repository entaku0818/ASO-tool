package service

import (
	"context"

	"github.com/entaku0818/aso-compass/backend/internal/repository"
)

type TrackedKeywordService struct {
	repo *repository.TrackedKeywordRepository
}

func NewTrackedKeywordService(repo *repository.TrackedKeywordRepository) *TrackedKeywordService {
	return &TrackedKeywordService{repo: repo}
}

func (s *TrackedKeywordService) List(ctx context.Context, userID string) ([]*repository.TrackedKeyword, error) {
	return s.repo.List(ctx, userID)
}

func (s *TrackedKeywordService) Get(ctx context.Context, userID, id string) (*repository.TrackedKeyword, error) {
	return s.repo.Get(ctx, userID, id)
}

func (s *TrackedKeywordService) GetSearchResults(ctx context.Context, userID, id string) ([]*repository.SearchResult, error) {
	// Verify ownership first
	_, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetLatestSearchResults(ctx, id)
}

func (s *TrackedKeywordService) SaveSearchResults(ctx context.Context, id string, results []repository.SearchResult) error {
	return s.repo.SaveSearchResults(ctx, id, results)
}
