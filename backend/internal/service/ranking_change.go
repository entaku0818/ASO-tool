package service

import (
	"context"

	"github.com/entaku0818/aso-compass/backend/internal/notification"
	"github.com/entaku0818/aso-compass/backend/internal/repository"
)

// RankingChangeService detects significant ranking changes
type RankingChangeService struct {
	rankingRepo *repository.RankingRepository
	keywordRepo *repository.KeywordRepository
	appRepo     *repository.AppRepository
}

// NewRankingChangeService creates a new ranking change service
func NewRankingChangeService(
	rankingRepo *repository.RankingRepository,
	keywordRepo *repository.KeywordRepository,
	appRepo *repository.AppRepository,
) *RankingChangeService {
	return &RankingChangeService{
		rankingRepo: rankingRepo,
		keywordRepo: keywordRepo,
		appRepo:     appRepo,
	}
}

// SignificantChangeThreshold defines how many positions change is considered significant
const SignificantChangeThreshold = 5

// DetectChanges detects significant ranking changes for all apps
func (s *RankingChangeService) DetectChanges(ctx context.Context) ([]notification.RankingChange, error) {
	apps, err := s.appRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	var changes []notification.RankingChange

	for _, app := range apps {
		appChanges, err := s.detectChangesForApp(ctx, app.ID, app.Name)
		if err != nil {
			continue // Skip errors for individual apps
		}
		changes = append(changes, appChanges...)
	}

	return changes, nil
}

// detectChangesForApp detects ranking changes for a specific app
func (s *RankingChangeService) detectChangesForApp(ctx context.Context, appID, appName string) ([]notification.RankingChange, error) {
	keywords, err := s.keywordRepo.ListByApp(ctx, appID)
	if err != nil {
		return nil, err
	}

	var changes []notification.RankingChange

	for _, keyword := range keywords {
		// Get the last 2 rankings for this keyword
		rankings, err := s.rankingRepo.ListByKeyword(ctx, keyword.ID, 2)
		if err != nil {
			continue
		}

		if len(rankings) < 2 {
			// Not enough data to compare
			// Check if this is a new ranking entry
			if len(rankings) == 1 && rankings[0].Rank != nil {
				changes = append(changes, notification.RankingChange{
					AppName:    appName,
					Keyword:    keyword.Keyword,
					Country:    keyword.Country,
					OldRank:    nil,
					NewRank:    rankings[0].Rank,
					ChangeType: "new_ranking",
				})
			}
			continue
		}

		newRanking := rankings[0]
		oldRanking := rankings[1]

		change := s.analyzeChange(appName, keyword.Keyword, keyword.Country, oldRanking.Rank, newRanking.Rank)
		if change != nil {
			changes = append(changes, *change)
		}
	}

	return changes, nil
}

// analyzeChange analyzes ranking change and returns a change if significant
func (s *RankingChangeService) analyzeChange(appName, keyword, country string, oldRank, newRank *int) *notification.RankingChange {
	// Both nil - no change
	if oldRank == nil && newRank == nil {
		return nil
	}

	// New entry (was not ranked, now ranked)
	if oldRank == nil && newRank != nil {
		return &notification.RankingChange{
			AppName:    appName,
			Keyword:    keyword,
			Country:    country,
			OldRank:    nil,
			NewRank:    newRank,
			ChangeType: "new_ranking",
		}
	}

	// Dropped out (was ranked, now not)
	if oldRank != nil && newRank == nil {
		return &notification.RankingChange{
			AppName:    appName,
			Keyword:    keyword,
			Country:    country,
			OldRank:    oldRank,
			NewRank:    nil,
			ChangeType: "dropped_out",
		}
	}

	// Both ranked - check for significant change
	diff := *oldRank - *newRank // positive = improved (lower rank is better)

	if diff >= SignificantChangeThreshold {
		return &notification.RankingChange{
			AppName:    appName,
			Keyword:    keyword,
			Country:    country,
			OldRank:    oldRank,
			NewRank:    newRank,
			ChangeType: "improved",
		}
	}

	if diff <= -SignificantChangeThreshold {
		return &notification.RankingChange{
			AppName:    appName,
			Keyword:    keyword,
			Country:    country,
			OldRank:    oldRank,
			NewRank:    newRank,
			ChangeType: "declined",
		}
	}

	return nil
}
