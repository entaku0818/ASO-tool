package scraper

import (
	"context"
	"time"
)

// AppInfo represents basic app information from store
type AppInfo struct {
	BundleID    string
	Name        string
	Developer   string
	Price       float64
	Currency    string
	Rating      float64
	RatingCount int
	Version     string
	StoreURL    string
	IconURL     string
	Description string
	ReleaseDate time.Time
}

// SearchResult represents a search result with ranking
type SearchResult struct {
	Rank    int
	AppInfo AppInfo
}

// Review represents a user review
type Review struct {
	ID         string
	Author     string
	Rating     int
	Title      string
	Content    string
	Version    string
	ReviewedAt time.Time
}

// Scraper interface for app store data fetching
type Scraper interface {
	// GetAppInfo fetches app information by bundle ID
	GetAppInfo(ctx context.Context, bundleID string, country string) (*AppInfo, error)

	// SearchKeyword searches for apps with a keyword and returns ranked results
	SearchKeyword(ctx context.Context, keyword string, country string, limit int) ([]SearchResult, error)

	// GetReviews fetches recent reviews for an app
	GetReviews(ctx context.Context, appID string, country string, page int) ([]Review, error)

	// GetAppRanking finds the rank of an app for a specific keyword
	GetAppRanking(ctx context.Context, bundleID string, keyword string, country string) (*int, error)
}
