package model

import "time"

// RankingChange represents a detected change in keyword ranking
type RankingChange struct {
	ID           string    `json:"id"`
	KeywordID    string    `json:"keyword_id"`
	Keyword      string    `json:"keyword"`
	Country      string    `json:"country"`
	AppID        string    `json:"app_id"`
	AppName      string    `json:"app_name"`
	OldRank      *int      `json:"old_rank"`
	NewRank      *int      `json:"new_rank"`
	Change       int       `json:"change"` // Positive = improved, Negative = dropped
	DetectedAt   time.Time `json:"detected_at"`
	IsSignificant bool     `json:"is_significant"` // True if change >= threshold
}

// RankingChangeType classifies the type of ranking change
type RankingChangeType string

const (
	RankingChangeTypeImproved    RankingChangeType = "improved"
	RankingChangeTypeDropped     RankingChangeType = "dropped"
	RankingChangeTypeNew         RankingChangeType = "new"         // First time ranking
	RankingChangeTypeLost        RankingChangeType = "lost"        // Dropped out of rankings
	RankingChangeTypeNoChange    RankingChangeType = "no_change"
)

// GetChangeType determines the type of ranking change
func (rc *RankingChange) GetChangeType() RankingChangeType {
	if rc.OldRank == nil && rc.NewRank != nil {
		return RankingChangeTypeNew
	}
	if rc.OldRank != nil && rc.NewRank == nil {
		return RankingChangeTypeLost
	}
	if rc.Change > 0 {
		return RankingChangeTypeImproved
	}
	if rc.Change < 0 {
		return RankingChangeTypeDropped
	}
	return RankingChangeTypeNoChange
}

// BatchResult represents the result of a batch job execution
type BatchResult struct {
	JobName       string           `json:"job_name"`
	StartedAt     time.Time        `json:"started_at"`
	CompletedAt   time.Time        `json:"completed_at"`
	Duration      time.Duration    `json:"duration"`
	AppsProcessed int              `json:"apps_processed"`
	KeywordsUpdated int            `json:"keywords_updated"`
	Changes       []*RankingChange `json:"changes"`
	Errors        []string         `json:"errors"`
	Success       bool             `json:"success"`
}
