package service

import (
	"context"

	"github.com/entaku0818/aso-tool/backend/internal/repository"
)

type PopularKeywordsService struct {
	trackedKeywordRepo *repository.TrackedKeywordRepository
}

func NewPopularKeywordsService(trackedKeywordRepo *repository.TrackedKeywordRepository) *PopularKeywordsService {
	return &PopularKeywordsService{
		trackedKeywordRepo: trackedKeywordRepo,
	}
}

// GetPopularKeywords returns popular keywords filtered by country and platform
func (s *PopularKeywordsService) GetPopularKeywords(ctx context.Context, country, platform string, limit int) ([]*repository.KeywordStats, error) {
	// Get all tracked keywords filtered by country and platform
	keywords, err := s.trackedKeywordRepo.ListByCountryAndPlatform(ctx, country, platform, limit)
	if err != nil {
		return nil, err
	}

	return keywords, nil
}
