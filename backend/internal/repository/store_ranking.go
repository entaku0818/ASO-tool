package repository

import (
	"context"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/scraper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoreRankingRepository struct {
	pool *pgxpool.Pool
}

func NewStoreRankingRepository(pool *pgxpool.Pool) *StoreRankingRepository {
	return &StoreRankingRepository{pool: pool}
}

// SaveRankings inserts a batch of ranking entries into the DB
func (r *StoreRankingRepository) SaveRankings(ctx context.Context, country, rankingType, genreID string, entries []scraper.AppRankingEntry) error {
	fetchedAt := time.Now()
	for _, entry := range entries {
		_, err := r.pool.Exec(ctx, `
			INSERT INTO store_rankings
				(id, country, ranking_type, genre_id, rank, app_id, app_name, developer, icon_url, category, store_url, price, release_date, fetched_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		`,
			uuid.New().String(),
			country, rankingType, genreID,
			entry.Rank, entry.AppID, entry.Name, entry.Developer,
			entry.IconURL, entry.Category, entry.StoreURL, entry.Price, entry.ReleaseDate,
			fetchedAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetLatestRankings returns the most recently fetched batch for the given filter
func (r *StoreRankingRepository) GetLatestRankings(ctx context.Context, country, rankingType, genreID string) ([]scraper.AppRankingEntry, time.Time, error) {
	// Find the latest fetched_at for this combination
	var latestAt time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT fetched_at FROM store_rankings
		WHERE country=$1 AND ranking_type=$2 AND genre_id=$3
		ORDER BY fetched_at DESC LIMIT 1
	`, country, rankingType, genreID).Scan(&latestAt)
	if err != nil {
		return nil, time.Time{}, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT rank, app_id, app_name, developer, icon_url, category, store_url, price, release_date
		FROM store_rankings
		WHERE country=$1 AND ranking_type=$2 AND genre_id=$3 AND fetched_at=$4
		ORDER BY rank ASC
	`, country, rankingType, genreID, latestAt)
	if err != nil {
		return nil, time.Time{}, err
	}
	defer rows.Close()

	var entries []scraper.AppRankingEntry
	for rows.Next() {
		var e scraper.AppRankingEntry
		if err := rows.Scan(&e.Rank, &e.AppID, &e.Name, &e.Developer, &e.IconURL, &e.Category, &e.StoreURL, &e.Price, &e.ReleaseDate); err != nil {
			return nil, time.Time{}, err
		}
		entries = append(entries, e)
	}
	return entries, latestAt, nil
}
