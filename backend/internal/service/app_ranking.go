package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/repository"
	"github.com/entaku0818/aso-tool/backend/internal/scraper"
)

type cachedRanking struct {
	entries   []scraper.AppRankingEntry
	fetchedAt time.Time
}

// defaultStoreRankingConfigs are the ranking combinations fetched by the daily batch job
var defaultStoreRankingConfigs = []struct{ country, rankingType, genreID string }{
	{"jp", "topfreeapplications", ""},
	{"jp", "toppaidapplications", ""},
	{"jp", "topgrossingapplications", ""},
	{"jp", "newfreeapplications", ""},
	{"jp", "newpaidapplications", ""},
	{"us", "topfreeapplications", ""},
	{"us", "toppaidapplications", ""},
	{"us", "topgrossingapplications", ""},
	{"kr", "topfreeapplications", ""},
	{"kr", "toppaidapplications", ""},
	{"kr", "topgrossingapplications", ""},
	{"gb", "topfreeapplications", ""},
	{"gb", "toppaidapplications", ""},
	{"gb", "topgrossingapplications", ""},
}

// AppRankingService fetches and caches App Store rankings with optional DB persistence
type AppRankingService struct {
	scraper   *scraper.ITunesRankingScraper
	repo      *repository.StoreRankingRepository
	mu        sync.RWMutex
	cache     map[string]cachedRanking
	memTTL    time.Duration // in-memory cache TTL
	dbStaleTTL time.Duration // treat DB data as stale after this duration
}

// NewAppRankingService creates a new AppRankingService.
// repo may be nil; in that case only in-memory caching is used (no DB persistence).
func NewAppRankingService(repo *repository.StoreRankingRepository) *AppRankingService {
	return &AppRankingService{
		scraper:    scraper.NewITunesRankingScraper(),
		repo:       repo,
		cache:      make(map[string]cachedRanking),
		memTTL:     30 * time.Minute,
		dbStaleTTL: 24 * time.Hour,
	}
}

// GetRanking returns app ranking entries.
// Priority: in-memory cache → DB (if configured, within 24h) → iTunes RSS API
func (s *AppRankingService) GetRanking(ctx context.Context, country, rankingType, genreID string, limit int) ([]scraper.AppRankingEntry, error) {
	key := fmt.Sprintf("%s:%s:%s:%d", country, rankingType, genreID, limit)

	// 1. In-memory cache
	s.mu.RLock()
	if cached, ok := s.cache[key]; ok && time.Since(cached.fetchedAt) < s.memTTL {
		s.mu.RUnlock()
		return cached.entries, nil
	}
	s.mu.RUnlock()

	// 2. DB cache (no limit filtering — return all stored entries)
	if s.repo != nil {
		entries, fetchedAt, err := s.repo.GetLatestRankings(ctx, country, rankingType, genreID)
		if err == nil && len(entries) > 0 && time.Since(fetchedAt) < s.dbStaleTTL {
			result := entries
			if limit > 0 && len(result) > limit {
				result = result[:limit]
			}
			s.mu.Lock()
			s.cache[key] = cachedRanking{entries: result, fetchedAt: fetchedAt}
			s.mu.Unlock()
			return result, nil
		}
	}

	// 3. Fetch from iTunes API
	entries, err := s.scraper.FetchRanking(ctx, country, rankingType, genreID, limit)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.cache[key] = cachedRanking{entries: entries, fetchedAt: time.Now()}
	s.mu.Unlock()

	return entries, nil
}

// GetRankingTrend returns the rank trend for a specific app over the last N days.
func (s *AppRankingService) GetRankingTrend(
	ctx context.Context, appID, country, rankingType string, days int,
) ([]repository.RankingTrendPoint, error) {
	return s.repo.GetRankingTrendForApp(ctx, appID, country, rankingType, days)
}

// GetCountryComparison returns the latest rank for an app across JP/US/KR/GB.
func (s *AppRankingService) GetCountryComparison(
	ctx context.Context, appID, rankingType string,
) ([]repository.CountryRankPoint, error) {
	countries := []string{"jp", "us", "kr", "gb"}
	return s.repo.GetRankingsAcrossCountries(ctx, appID, rankingType, countries)
}

// FetchAndSaveAll fetches the default set of rankings and persists them to the DB.
// Used by the daily batch job. Returns the total number of entries saved and any errors.
func (s *AppRankingService) FetchAndSaveAll(ctx context.Context) (int, []error) {
	if s.repo == nil {
		return 0, []error{fmt.Errorf("no repository configured")}
	}

	total := 0
	var errs []error

	for _, cfg := range defaultStoreRankingConfigs {
		log.Printf("Fetching store rankings: country=%s type=%s genre=%s", cfg.country, cfg.rankingType, cfg.genreID)
		entries, err := s.scraper.FetchRanking(ctx, cfg.country, cfg.rankingType, cfg.genreID, 100)
		if err != nil {
			errs = append(errs, fmt.Errorf("fetch %s/%s/%s: %w", cfg.country, cfg.rankingType, cfg.genreID, err))
			continue
		}
		if err := s.repo.SaveRankings(ctx, cfg.country, cfg.rankingType, cfg.genreID, entries); err != nil {
			errs = append(errs, fmt.Errorf("save %s/%s/%s: %w", cfg.country, cfg.rankingType, cfg.genreID, err))
			continue
		}
		total += len(entries)
		log.Printf("Saved %d entries for %s/%s/%s", len(entries), cfg.country, cfg.rankingType, cfg.genreID)
	}

	return total, errs
}
