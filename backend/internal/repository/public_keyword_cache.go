package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PublicKeywordCacheEntry struct {
	Keyword    string    `json:"keyword"`
	Country    string    `json:"country"`
	Genre      string    `json:"genre"`
	Popularity int       `json:"popularity"`
	FetchedAt  time.Time `json:"fetched_at"`
}

type PublicKeywordCacheRepository struct {
	pool *pgxpool.Pool
}

func NewPublicKeywordCacheRepository(pool *pgxpool.Pool) *PublicKeywordCacheRepository {
	return &PublicKeywordCacheRepository{pool: pool}
}

func (r *PublicKeywordCacheRepository) UpsertMany(ctx context.Context, entries []PublicKeywordCacheEntry) error {
	for _, e := range entries {
		_, err := r.pool.Exec(ctx, `
			INSERT INTO public_keyword_cache (keyword, country, genre, popularity, fetched_at)
			VALUES ($1, $2, $3, $4, NOW())
			ON CONFLICT (keyword, country, genre)
			DO UPDATE SET popularity = EXCLUDED.popularity, fetched_at = NOW()
		`, e.Keyword, e.Country, e.Genre, e.Popularity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PublicKeywordCacheRepository) Search(ctx context.Context, country, genre string, limit int) ([]PublicKeywordCacheEntry, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}

	var query string
	var args []any
	if genre == "" {
		query = `
			SELECT keyword, country, genre, popularity, fetched_at
			FROM public_keyword_cache
			WHERE country = $1
			ORDER BY popularity DESC
			LIMIT $2`
		args = []any{country, limit}
	} else {
		query = `
			SELECT keyword, country, genre, popularity, fetched_at
			FROM public_keyword_cache
			WHERE country = $1 AND genre = $2
			ORDER BY popularity DESC
			LIMIT $3`
		args = []any{country, genre, limit}
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []PublicKeywordCacheEntry
	for rows.Next() {
		var e PublicKeywordCacheEntry
		if err := rows.Scan(&e.Keyword, &e.Country, &e.Genre, &e.Popularity, &e.FetchedAt); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}
