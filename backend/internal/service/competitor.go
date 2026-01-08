package service

import (
	"context"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
	"github.com/entaku0818/aso-tool/backend/internal/scraper"
)

type CompetitorService struct {
	repo              *repository.CompetitorRepository
	keywordRepo       *repository.KeywordRepository
	appStoreScraper   *scraper.AppStoreScraper
	googlePlayScraper *scraper.GooglePlayScraper
}

func NewCompetitorService(
	repo *repository.CompetitorRepository,
	keywordRepo *repository.KeywordRepository,
) *CompetitorService {
	return &CompetitorService{
		repo:              repo,
		keywordRepo:       keywordRepo,
		appStoreScraper:   scraper.NewAppStoreScraper(),
		googlePlayScraper: scraper.NewGooglePlayScraper(),
	}
}

func (s *CompetitorService) Create(ctx context.Context, req *model.CreateCompetitorRequest) (*model.Competitor, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, req)
}

func (s *CompetitorService) Get(ctx context.Context, id string) (*model.Competitor, error) {
	return s.repo.Get(ctx, id)
}

func (s *CompetitorService) ListByApp(ctx context.Context, appID string) ([]*model.Competitor, error) {
	return s.repo.ListByApp(ctx, appID)
}

func (s *CompetitorService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// UpdateCompetitorRankings fetches and stores rankings for all competitors of an app
func (s *CompetitorService) UpdateCompetitorRankings(ctx context.Context, appID string) (int, error) {
	competitors, err := s.repo.ListByApp(ctx, appID)
	if err != nil {
		return 0, err
	}

	keywords, err := s.keywordRepo.ListByApp(ctx, appID)
	if err != nil {
		return 0, err
	}

	updated := 0
	for _, competitor := range competitors {
		for _, keyword := range keywords {
			var rank *int
			var fetchErr error

			switch competitor.Platform {
			case model.PlatformIOS:
				rank, fetchErr = s.appStoreScraper.GetAppRanking(ctx, competitor.CompetitorBundleID, keyword.Keyword, keyword.Country)
			case model.PlatformAndroid:
				rank, fetchErr = s.googlePlayScraper.GetAppRanking(ctx, competitor.CompetitorBundleID, keyword.Keyword, keyword.Country)
			}

			if fetchErr != nil {
				continue
			}

			_, err := s.repo.CreateRanking(ctx, &model.CreateCompetitorRankingRequest{
				CompetitorID: competitor.ID,
				KeywordID:    keyword.ID,
				Rank:         rank,
			})
			if err != nil {
				continue
			}
			updated++
		}
	}

	return updated, nil
}

func (s *CompetitorService) GetComparison(ctx context.Context, appID string, keywordID string) (*model.CompetitorComparison, error) {
	return s.repo.GetComparisonData(ctx, appID, keywordID)
}

// GetLatestRanking returns the latest ranking for a competitor-keyword pair
func (s *CompetitorService) GetLatestRanking(ctx context.Context, competitorID string, keywordID string) (*model.CompetitorRanking, error) {
	return s.repo.GetLatestRankings(ctx, competitorID, keywordID)
}
