package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackedKeyword struct {
	ID        string    `json:"id"`
	Keyword   string    `json:"keyword"`
	Country   string    `json:"country"`
	Platform  string    `json:"platform"`
	UserID    string    `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type SearchResult struct {
	ID               string    `json:"id"`
	TrackedKeywordID string    `json:"tracked_keyword_id"`
	Rank             int       `json:"rank"`
	AppName          string    `json:"app_name"`
	BundleID         string    `json:"bundle_id"`
	Developer        string    `json:"developer"`
	RecordedAt       time.Time `json:"recorded_at"`
}

type TrackedKeywordRepository struct {
	pool *pgxpool.Pool
}

func NewTrackedKeywordRepository(pool *pgxpool.Pool) *TrackedKeywordRepository {
	return &TrackedKeywordRepository{pool: pool}
}

func (r *TrackedKeywordRepository) Create(ctx context.Context, userID, keyword, country, platform string) (*TrackedKeyword, error) {
	id := uuid.New().String()
	tk := &TrackedKeyword{
		ID:        id,
		Keyword:   keyword,
		Country:   country,
		Platform:  platform,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	_, err := r.pool.Exec(ctx,
		`INSERT INTO tracked_keywords (id, keyword, country, platform, user_id, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (keyword, country, platform, user_id) DO NOTHING`,
		tk.ID, tk.Keyword, tk.Country, tk.Platform, tk.UserID, tk.CreatedAt)
	if err != nil {
		return nil, err
	}

	return tk, nil
}

func (r *TrackedKeywordRepository) GetByKeyword(ctx context.Context, userID, keyword, country, platform string) (*TrackedKeyword, error) {
	tk := &TrackedKeyword{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, keyword, country, platform, user_id, created_at
		 FROM tracked_keywords
		 WHERE keyword = $1 AND country = $2 AND platform = $3 AND user_id = $4`,
		keyword, country, platform, userID).Scan(&tk.ID, &tk.Keyword, &tk.Country, &tk.Platform, &tk.UserID, &tk.CreatedAt)
	if err != nil {
		return nil, err
	}
	return tk, nil
}

func (r *TrackedKeywordRepository) List(ctx context.Context, userID string) ([]*TrackedKeyword, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, keyword, country, platform, user_id, created_at
		 FROM tracked_keywords
		 WHERE user_id = $1
		 ORDER BY keyword`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []*TrackedKeyword
	for rows.Next() {
		tk := &TrackedKeyword{}
		if err := rows.Scan(&tk.ID, &tk.Keyword, &tk.Country, &tk.Platform, &tk.UserID, &tk.CreatedAt); err != nil {
			return nil, err
		}
		keywords = append(keywords, tk)
	}
	return keywords, nil
}

// ListAll returns all tracked keywords without user filtering (for internal use like scraper)
func (r *TrackedKeywordRepository) ListAll(ctx context.Context) ([]*TrackedKeyword, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, keyword, country, platform, user_id, created_at
		 FROM tracked_keywords
		 ORDER BY keyword`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []*TrackedKeyword
	for rows.Next() {
		tk := &TrackedKeyword{}
		if err := rows.Scan(&tk.ID, &tk.Keyword, &tk.Country, &tk.Platform, &tk.UserID, &tk.CreatedAt); err != nil {
			return nil, err
		}
		keywords = append(keywords, tk)
	}
	return keywords, nil
}

func (r *TrackedKeywordRepository) Get(ctx context.Context, userID, id string) (*TrackedKeyword, error) {
	tk := &TrackedKeyword{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, keyword, country, platform, user_id, created_at
		 FROM tracked_keywords
		 WHERE id = $1 AND user_id = $2`,
		id, userID).Scan(&tk.ID, &tk.Keyword, &tk.Country, &tk.Platform, &tk.UserID, &tk.CreatedAt)
	if err != nil {
		return nil, err
	}
	return tk, nil
}

// GetByID returns a tracked keyword without user filtering (for internal use like scraper)
func (r *TrackedKeywordRepository) GetByID(ctx context.Context, id string) (*TrackedKeyword, error) {
	tk := &TrackedKeyword{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, keyword, country, platform, user_id, created_at
		 FROM tracked_keywords
		 WHERE id = $1`,
		id).Scan(&tk.ID, &tk.Keyword, &tk.Country, &tk.Platform, &tk.UserID, &tk.CreatedAt)
	if err != nil {
		return nil, err
	}
	return tk, nil
}

func (r *TrackedKeywordRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM tracked_keywords WHERE id = $1`, id)
	return err
}

// Search Results

func (r *TrackedKeywordRepository) SaveSearchResults(ctx context.Context, trackedKeywordID string, results []SearchResult) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Delete old results from today
	_, err = tx.Exec(ctx,
		`DELETE FROM search_results
		 WHERE tracked_keyword_id = $1
		 AND recorded_at::date = CURRENT_DATE`,
		trackedKeywordID)
	if err != nil {
		return err
	}

	// Insert new results
	for _, result := range results {
		id := uuid.New().String()
		_, err = tx.Exec(ctx,
			`INSERT INTO search_results (id, tracked_keyword_id, rank, app_name, bundle_id, developer, recorded_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			id, trackedKeywordID, result.Rank, result.AppName, result.BundleID, result.Developer, time.Now())
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *TrackedKeywordRepository) GetLatestSearchResults(ctx context.Context, trackedKeywordID string) ([]*SearchResult, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT ON (rank) id, tracked_keyword_id, rank, app_name, bundle_id, developer, recorded_at
		 FROM search_results
		 WHERE tracked_keyword_id = $1
		 ORDER BY rank, recorded_at DESC`,
		trackedKeywordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*SearchResult
	for rows.Next() {
		sr := &SearchResult{}
		if err := rows.Scan(&sr.ID, &sr.TrackedKeywordID, &sr.Rank, &sr.AppName, &sr.BundleID, &sr.Developer, &sr.RecordedAt); err != nil {
			return nil, err
		}
		results = append(results, sr)
	}
	return results, nil
}
