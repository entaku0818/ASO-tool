package model

import "time"

type ASCAnalytics struct {
	ID             string    `json:"id"`
	AppID          string    `json:"app_id"`
	Date           time.Time `json:"date"`
	Impressions    int       `json:"impressions"`
	Downloads      int       `json:"downloads"`
	PageViews      int       `json:"page_views"`
	ConversionRate *float64  `json:"conversion_rate,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type AppVersion struct {
	ID           string    `json:"id"`
	AppID        string    `json:"app_id"`
	Version      string    `json:"version"`
	ReleaseDate  time.Time `json:"release_date"`
	ReleaseNotes string    `json:"release_notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// AnalyticsWithCorrelation includes analytics data with correlation indicators
type AnalyticsWithCorrelation struct {
	ASCAnalytics
	RankingChange       *int     `json:"ranking_change,omitempty"`
	NewVersion          *string  `json:"new_version,omitempty"`
	ReviewCountChange   *int     `json:"review_count_change,omitempty"`
	AverageRatingChange *float64 `json:"average_rating_change,omitempty"`
}

type AnalyticsSummary struct {
	TotalImpressions        int64   `json:"total_impressions"`
	TotalDownloads          int64   `json:"total_downloads"`
	AverageConversion       float64 `json:"average_conversion"`
	ImpressionsChangePercent float64 `json:"impressions_change_percent"`
	DownloadsChangePercent   float64 `json:"downloads_change_percent"`
}

type CreateAppVersionRequest struct {
	Version      string `json:"version"`
	ReleaseDate  string `json:"release_date"` // YYYY-MM-DD format
	ReleaseNotes string `json:"release_notes,omitempty"`
}

func (r *CreateAppVersionRequest) Validate() error {
	if r.Version == "" {
		return ErrVersionRequired
	}
	if r.ReleaseDate == "" {
		return ErrReleaseDateRequired
	}
	return nil
}
