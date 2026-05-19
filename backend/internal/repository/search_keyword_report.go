package repository

import (
	"context"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/appstoreconnect"
	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchKeywordReportRepository struct {
	pool *pgxpool.Pool
}

func NewSearchKeywordReportRepository(pool *pgxpool.Pool) *SearchKeywordReportRepository {
	return &SearchKeywordReportRepository{pool: pool}
}

func (r *SearchKeywordReportRepository) CreateReportRequest(ctx context.Context, appID, requestID string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO asc_report_requests (id, app_id, request_id, status)
		VALUES ($1, $2, $3, 'pending')
		ON CONFLICT (app_id, request_id) DO NOTHING
	`, uuid.New().String(), appID, requestID)
	return err
}

func (r *SearchKeywordReportRepository) GetLatestReportRequest(ctx context.Context, appID string) (*model.ASCReportRequest, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, app_id, request_id, status, created_at, updated_at
		FROM asc_report_requests
		WHERE app_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, appID)

	var req model.ASCReportRequest
	err := row.Scan(&req.ID, &req.AppID, &req.RequestID, &req.Status, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *SearchKeywordReportRepository) UpdateReportRequestStatus(ctx context.Context, appID, requestID, status string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE asc_report_requests SET status = $1, updated_at = NOW()
		WHERE app_id = $2 AND request_id = $3
	`, status, appID, requestID)
	return err
}

func (r *SearchKeywordReportRepository) SaveSearchKeywords(ctx context.Context, appID string, entries []appstoreconnect.SearchTermEntry) error {
	for _, e := range entries {
		date, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			continue
		}
		_, err = r.pool.Exec(ctx, `
			INSERT INTO asc_search_keywords (id, app_id, date, keyword, impressions, page_views, installs)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (app_id, date, keyword) DO UPDATE SET
				impressions = EXCLUDED.impressions,
				page_views  = EXCLUDED.page_views,
				installs    = EXCLUDED.installs,
				fetched_at  = NOW()
		`, uuid.New().String(), appID, date, e.Keyword, e.Impressions, e.PageViews, e.Installs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SearchKeywordReportRepository) GetSearchKeywords(ctx context.Context, appID string, days int) ([]*model.ASCSearchKeyword, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, app_id, date, keyword, impressions, page_views, installs, fetched_at
		FROM asc_search_keywords
		WHERE app_id = $1 AND date >= NOW() - INTERVAL '1 day' * $2
		ORDER BY impressions DESC, date DESC
	`, appID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.ASCSearchKeyword
	for rows.Next() {
		k := &model.ASCSearchKeyword{}
		if err := rows.Scan(&k.ID, &k.AppID, &k.Date, &k.Keyword, &k.Impressions, &k.PageViews, &k.Installs, &k.FetchedAt); err != nil {
			return nil, err
		}
		result = append(result, k)
	}
	return result, nil
}
