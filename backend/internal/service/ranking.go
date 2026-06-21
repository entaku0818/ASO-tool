package service

import (
	"context"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
)

type RankingService struct {
	repo *repository.RankingRepository
}

func NewRankingService(repo *repository.RankingRepository) *RankingService {
	return &RankingService{repo: repo}
}

func (s *RankingService) Create(ctx context.Context, req *model.CreateRankingRequest) (*model.RankingHistory, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, req)
}

func (s *RankingService) ListByKeyword(ctx context.Context, keywordID string, limit int) ([]*model.RankingHistory, error) {
	if limit <= 0 {
		limit = 30
	}
	return s.repo.ListByKeyword(ctx, keywordID, limit)
}

func (s *RankingService) ListByKeywordDays(ctx context.Context, keywordID string, days int) ([]*model.RankingHistory, error) {
	if days <= 0 {
		days = 30
	}
	return s.repo.ListByKeywordDays(ctx, keywordID, days)
}

func (s *RankingService) ListByApp(ctx context.Context, appID string, days int) ([]*model.RankingWithKeyword, error) {
	if days <= 0 {
		days = 30
	}
	to := time.Now()
	from := to.AddDate(0, 0, -days)
	return s.repo.ListByAppWithKeyword(ctx, appID, from, to)
}

func (s *RankingService) GetLatestByKeyword(ctx context.Context, keywordID string) (*model.RankingHistory, error) {
	return s.repo.GetLatestByKeyword(ctx, keywordID)
}

func (s *RankingService) GetRisingKeywords(ctx context.Context, appID string, days int) ([]repository.RisingKeyword, error) {
	if days <= 0 {
		days = 7
	}
	return s.repo.GetRisingKeywords(ctx, appID, days)
}

// GetAllKeywordRanks returns current/previous ranks for all keywords of an app.
func (s *RankingService) GetAllKeywordRanks(ctx context.Context, appID string) ([]repository.KeywordRankSummary, error) {
	return s.repo.GetAllKeywordRanks(ctx, appID)
}
