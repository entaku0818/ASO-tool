package repository

import (
	"context"
	"time"

	"github.com/entaku0818/aso-compass/backend/internal/scraper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoreRankingRepository struct {
	pool *pgxpool.Pool
}

func NewStoreRankingRepository(pool *pgxpool.Pool) *StoreRankingRepository {
	return &StoreRankingRepository{pool: pool}
}

// RankingTrendPoint holds a single rank data point for a given date.
type RankingTrendPoint struct {
	Date string `json:"date"`
	Rank int    `json:"rank"`
}

// CountryRankPoint holds a rank for a specific country.
type CountryRankPoint struct {
	Country   string    `json:"country"`
	Rank      int       `json:"rank"`
	FetchedAt time.Time `json:"fetched_at"`
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

// GetRankingTrendForApp returns daily rank data for a specific app over the last N days.
// For each day, the rank from the most recent fetch of that day is returned.
func (r *StoreRankingRepository) GetRankingTrendForApp(
	ctx context.Context, appID, country, rankingType string, days int,
) ([]RankingTrendPoint, error) {
	since := time.Now().AddDate(0, 0, -days)
	rows, err := r.pool.Query(ctx, `
		WITH latest_per_day AS (
			SELECT
				DATE(fetched_at) AS day,
				MAX(fetched_at) AS max_fetched_at
			FROM store_rankings
			WHERE app_id = $1 AND country = $2 AND ranking_type = $3
			  AND fetched_at >= $4
			GROUP BY DATE(fetched_at)
		)
		SELECT
			TO_CHAR(lpd.day, 'YYYY-MM-DD') AS date,
			sr.rank
		FROM store_rankings sr
		JOIN latest_per_day lpd
		  ON DATE(sr.fetched_at) = lpd.day AND sr.fetched_at = lpd.max_fetched_at
		WHERE sr.app_id = $1 AND sr.country = $2 AND sr.ranking_type = $3
		ORDER BY lpd.day ASC
	`, appID, country, rankingType, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []RankingTrendPoint
	for rows.Next() {
		var p RankingTrendPoint
		if err := rows.Scan(&p.Date, &p.Rank); err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	return points, nil
}

// GetRankingsAcrossCountries returns the latest rank for a specific app in each given country.
func (r *StoreRankingRepository) GetRankingsAcrossCountries(
	ctx context.Context, appID, rankingType string, countries []string,
) ([]CountryRankPoint, error) {
	rows, err := r.pool.Query(ctx, `
		WITH latest_per_country AS (
			SELECT
				country,
				MAX(fetched_at) AS max_fetched_at
			FROM store_rankings
			WHERE app_id = $1 AND ranking_type = $2 AND country = ANY($3)
			GROUP BY country
		)
		SELECT sr.country, sr.rank, sr.fetched_at
		FROM store_rankings sr
		JOIN latest_per_country lpc
		  ON sr.country = lpc.country AND sr.fetched_at = lpc.max_fetched_at
		WHERE sr.app_id = $1 AND sr.ranking_type = $2
		ORDER BY sr.rank ASC
	`, appID, rankingType, countries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []CountryRankPoint
	for rows.Next() {
		var p CountryRankPoint
		if err := rows.Scan(&p.Country, &p.Rank, &p.FetchedAt); err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	return points, nil
}
