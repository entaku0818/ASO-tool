package repository

import (
	"context"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsRepository struct {
	pool *pgxpool.Pool
}

func NewAnalyticsRepository(pool *pgxpool.Pool) *AnalyticsRepository {
	return &AnalyticsRepository{pool: pool}
}

func (r *AnalyticsRepository) Upsert(ctx context.Context, analytics *model.ASCAnalytics) error {
	if analytics.ID == "" {
		analytics.ID = uuid.New().String()
	}

	query := `
		INSERT INTO asc_analytics (id, app_id, date, impressions, downloads, page_views, conversion_rate, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (app_id, date) DO UPDATE SET
			impressions = EXCLUDED.impressions,
			downloads = EXCLUDED.downloads,
			page_views = EXCLUDED.page_views,
			conversion_rate = EXCLUDED.conversion_rate
	`

	_, err := r.pool.Exec(ctx, query,
		analytics.ID, analytics.AppID, analytics.Date,
		analytics.Impressions, analytics.Downloads, analytics.PageViews, analytics.ConversionRate,
	)

	return err
}

func (r *AnalyticsRepository) GetByDateRange(ctx context.Context, appID string, start, end time.Time) ([]*model.ASCAnalytics, error) {
	query := `
		SELECT id, app_id, date, impressions, downloads, page_views, conversion_rate, created_at
		FROM asc_analytics
		WHERE app_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
	`

	rows, err := r.pool.Query(ctx, query, appID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []*model.ASCAnalytics
	for rows.Next() {
		a := &model.ASCAnalytics{}
		err := rows.Scan(
			&a.ID, &a.AppID, &a.Date,
			&a.Impressions, &a.Downloads, &a.PageViews, &a.ConversionRate,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		analytics = append(analytics, a)
	}

	return analytics, nil
}

func (r *AnalyticsRepository) GetLatest(ctx context.Context, appID string, days int) ([]*model.ASCAnalytics, error) {
	query := `
		SELECT id, app_id, date, impressions, downloads, page_views, conversion_rate, created_at
		FROM asc_analytics
		WHERE app_id = $1 AND date >= NOW() - INTERVAL '1 day' * $2
		ORDER BY date DESC
	`

	rows, err := r.pool.Query(ctx, query, appID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []*model.ASCAnalytics
	for rows.Next() {
		a := &model.ASCAnalytics{}
		err := rows.Scan(
			&a.ID, &a.AppID, &a.Date,
			&a.Impressions, &a.Downloads, &a.PageViews, &a.ConversionRate,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		analytics = append(analytics, a)
	}

	return analytics, nil
}

func (r *AnalyticsRepository) GetSummary(ctx context.Context, appID string, start, end time.Time) (*model.AnalyticsSummary, error) {
	// Get current period totals
	currentQuery := `
		SELECT COALESCE(SUM(impressions), 0), COALESCE(SUM(downloads), 0), COALESCE(AVG(conversion_rate), 0)
		FROM asc_analytics
		WHERE app_id = $1 AND date >= $2 AND date <= $3
	`

	var totalImpressions, totalDownloads int64
	var avgConversion float64
	err := r.pool.QueryRow(ctx, currentQuery, appID, start, end).Scan(
		&totalImpressions, &totalDownloads, &avgConversion,
	)
	if err != nil {
		return nil, err
	}

	// Calculate period length for previous period comparison
	periodDays := int(end.Sub(start).Hours() / 24)
	prevStart := start.AddDate(0, 0, -periodDays)
	prevEnd := start.AddDate(0, 0, -1)

	// Get previous period totals
	prevQuery := `
		SELECT COALESCE(SUM(impressions), 0), COALESCE(SUM(downloads), 0)
		FROM asc_analytics
		WHERE app_id = $1 AND date >= $2 AND date <= $3
	`

	var prevImpressions, prevDownloads int64
	err = r.pool.QueryRow(ctx, prevQuery, appID, prevStart, prevEnd).Scan(
		&prevImpressions, &prevDownloads,
	)
	if err != nil {
		return nil, err
	}

	// Calculate percentage changes
	var impressionsChange, downloadsChange float64
	if prevImpressions > 0 {
		impressionsChange = float64(totalImpressions-prevImpressions) / float64(prevImpressions) * 100
	}
	if prevDownloads > 0 {
		downloadsChange = float64(totalDownloads-prevDownloads) / float64(prevDownloads) * 100
	}

	return &model.AnalyticsSummary{
		TotalImpressions:         totalImpressions,
		TotalDownloads:           totalDownloads,
		AverageConversion:        avgConversion,
		ImpressionsChangePercent: impressionsChange,
		DownloadsChangePercent:   downloadsChange,
	}, nil
}
