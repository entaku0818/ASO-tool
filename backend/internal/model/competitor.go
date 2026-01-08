package model

import (
	"time"
)

type Competitor struct {
	ID                 string    `json:"id"`
	AppID              string    `json:"app_id"`
	CompetitorBundleID string    `json:"competitor_bundle_id"`
	CompetitorName     string    `json:"competitor_name"`
	Platform           Platform  `json:"platform"`
	Notes              string    `json:"notes,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

type CompetitorRanking struct {
	ID           string    `json:"id"`
	CompetitorID string    `json:"competitor_id"`
	KeywordID    string    `json:"keyword_id"`
	Rank         *int      `json:"rank"`
	RecordedAt   time.Time `json:"recorded_at"`
}

type CreateCompetitorRequest struct {
	AppID              string   `json:"app_id"`
	CompetitorBundleID string   `json:"competitor_bundle_id"`
	CompetitorName     string   `json:"competitor_name"`
	Platform           Platform `json:"platform"`
	Notes              string   `json:"notes,omitempty"`
}

func (r *CreateCompetitorRequest) Validate() error {
	if r.AppID == "" {
		return ErrAppIDRequired
	}
	if r.CompetitorBundleID == "" {
		return ErrBundleIDRequired
	}
	if r.CompetitorName == "" {
		return ErrNameRequired
	}
	if r.Platform != PlatformIOS && r.Platform != PlatformAndroid {
		return ErrInvalidPlatform
	}
	return nil
}

type CreateCompetitorRankingRequest struct {
	CompetitorID string `json:"competitor_id"`
	KeywordID    string `json:"keyword_id"`
	Rank         *int   `json:"rank"`
}

func (r *CreateCompetitorRankingRequest) Validate() error {
	if r.CompetitorID == "" {
		return ErrCompetitorIDRequired
	}
	if r.KeywordID == "" {
		return ErrKeywordIDRequired
	}
	return nil
}

// CompetitorComparison is used for comparing app rankings with competitors
type CompetitorComparison struct {
	Keyword      string               `json:"keyword"`
	KeywordID    string               `json:"keyword_id"`
	AppRank      *int                 `json:"app_rank"`
	Competitors  []CompetitorRankData `json:"competitors"`
	RecordedAt   time.Time            `json:"recorded_at"`
}

type CompetitorRankData struct {
	CompetitorID   string `json:"competitor_id"`
	CompetitorName string `json:"competitor_name"`
	BundleID       string `json:"bundle_id"`
	Rank           *int   `json:"rank"`
}
