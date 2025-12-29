package model

import "time"

type Review struct {
	ID         string     `json:"id"`
	AppID      string     `json:"app_id"`
	ReviewID   string     `json:"review_id"`
	Author     string     `json:"author,omitempty"`
	Rating     int        `json:"rating"`
	Title      string     `json:"title,omitempty"`
	Content    string     `json:"content,omitempty"`
	Version    string     `json:"version,omitempty"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type CreateReviewRequest struct {
	AppID      string     `json:"app_id"`
	ReviewID   string     `json:"review_id"`
	Author     string     `json:"author,omitempty"`
	Rating     int        `json:"rating"`
	Title      string     `json:"title,omitempty"`
	Content    string     `json:"content,omitempty"`
	Version    string     `json:"version,omitempty"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
}

func (r *CreateReviewRequest) Validate() error {
	if r.AppID == "" {
		return ErrAppIDRequired
	}
	if r.ReviewID == "" {
		return ErrReviewIDRequired
	}
	if r.Rating < 1 || r.Rating > 5 {
		return ErrInvalidRating
	}
	return nil
}

type ReviewStats struct {
	AppID        string  `json:"app_id"`
	TotalCount   int     `json:"total_count"`
	AverageRating float64 `json:"average_rating"`
	Rating1Count int     `json:"rating_1_count"`
	Rating2Count int     `json:"rating_2_count"`
	Rating3Count int     `json:"rating_3_count"`
	Rating4Count int     `json:"rating_4_count"`
	Rating5Count int     `json:"rating_5_count"`
}
