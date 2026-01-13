package service

import (
	"context"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
	"github.com/google/uuid"
)

const (
	// DefaultSignificantChangeThreshold is the default minimum rank change to be considered significant
	DefaultSignificantChangeThreshold = 5
)

// RankingChangeService detects ranking changes
type RankingChangeService struct {
	keywordRepo *repository.KeywordRepository
	rankingRepo *repository.RankingRepository
	appRepo     *repository.AppRepository
	threshold   int
}

// NewRankingChangeService creates a new ranking change service
func NewRankingChangeService(
	keywordRepo *repository.KeywordRepository,
	rankingRepo *repository.RankingRepository,
	appRepo *repository.AppRepository,
) *RankingChangeService {
	return &RankingChangeService{
		keywordRepo: keywordRepo,
		rankingRepo: rankingRepo,
		appRepo:     appRepo,
		threshold:   DefaultSignificantChangeThreshold,
	}
}

// SetThreshold sets the significance threshold
func (s *RankingChangeService) SetThreshold(threshold int) {
	s.threshold = threshold
}

// DetectChangesForApp detects ranking changes for a specific app
func (s *RankingChangeService) DetectChangesForApp(ctx context.Context, appID string) ([]*model.RankingChange, error) {
	app, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return nil, err
	}

	keywords, err := s.keywordRepo.ListByApp(ctx, appID)
	if err != nil {
		return nil, err
	}

	var changes []*model.RankingChange

	for _, keyword := range keywords {
		// Get the last two rankings for this keyword
		rankings, err := s.rankingRepo.ListByKeyword(ctx, keyword.ID, 2)
		if err != nil {
			continue
		}

		if len(rankings) < 1 {
			continue
		}

		var change *model.RankingChange
		now := time.Now()

		if len(rankings) == 1 {
			// First ranking - check if it's a new entry
			if rankings[0].Rank != nil {
				change = &model.RankingChange{
					ID:           uuid.New().String(),
					KeywordID:    keyword.ID,
					Keyword:      keyword.Keyword,
					Country:      keyword.Country,
					AppID:        appID,
					AppName:      app.Name,
					OldRank:      nil,
					NewRank:      rankings[0].Rank,
					Change:       0,
					DetectedAt:   now,
					IsSignificant: true, // New rankings are always significant
				}
			}
		} else if len(rankings) >= 2 {
			// Compare latest with previous
			current := rankings[0]
			previous := rankings[1]

			// Skip if both are the same
			if equalRanks(current.Rank, previous.Rank) {
				continue
			}

			change = s.calculateChange(keyword, app, previous.Rank, current.Rank, now)
		}

		if change != nil {
			changes = append(changes, change)
		}
	}

	return changes, nil
}

// DetectChangesForAllApps detects ranking changes for all apps
func (s *RankingChangeService) DetectChangesForAllApps(ctx context.Context) ([]*model.RankingChange, error) {
	apps, err := s.appRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	var allChanges []*model.RankingChange
	for _, app := range apps {
		changes, err := s.DetectChangesForApp(ctx, app.ID)
		if err != nil {
			continue
		}
		allChanges = append(allChanges, changes...)
	}

	return allChanges, nil
}

func (s *RankingChangeService) calculateChange(
	keyword *model.Keyword,
	app *model.App,
	oldRank, newRank *int,
	detectedAt time.Time,
) *model.RankingChange {
	change := &model.RankingChange{
		ID:         uuid.New().String(),
		KeywordID:  keyword.ID,
		Keyword:    keyword.Keyword,
		Country:    keyword.Country,
		AppID:      app.ID,
		AppName:    app.Name,
		OldRank:    oldRank,
		NewRank:    newRank,
		DetectedAt: detectedAt,
	}

	// Calculate the change value
	// Positive = improved (rank went down, e.g., 10 -> 5)
	// Negative = dropped (rank went up, e.g., 5 -> 10)
	if oldRank == nil && newRank != nil {
		// New entry - treat as significant improvement
		change.Change = *newRank
		change.IsSignificant = true
	} else if oldRank != nil && newRank == nil {
		// Lost from rankings
		change.Change = -(*oldRank)
		change.IsSignificant = true
	} else if oldRank != nil && newRank != nil {
		// Regular change: improvement = old - new (positive is better)
		change.Change = *oldRank - *newRank
		absChange := change.Change
		if absChange < 0 {
			absChange = -absChange
		}
		change.IsSignificant = absChange >= s.threshold
	}

	return change
}

func equalRanks(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
