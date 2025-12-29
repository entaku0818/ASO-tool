package model

import (
	"errors"
	"testing"
	"time"
)

func TestCreateReviewRequest_Validate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		req     CreateReviewRequest
		wantErr error
	}{
		{
			name: "valid request with all fields",
			req: CreateReviewRequest{
				AppID:      "app-123",
				ReviewID:   "review-456",
				Author:     "Test User",
				Rating:     5,
				Title:      "Great App",
				Content:    "Really love this app!",
				Version:    "1.0.0",
				ReviewedAt: &now,
			},
			wantErr: nil,
		},
		{
			name: "valid request with minimum fields",
			req: CreateReviewRequest{
				AppID:    "app-123",
				ReviewID: "review-456",
				Rating:   3,
			},
			wantErr: nil,
		},
		{
			name: "valid rating 1",
			req: CreateReviewRequest{
				AppID:    "app-123",
				ReviewID: "review-456",
				Rating:   1,
			},
			wantErr: nil,
		},
		{
			name: "valid rating 5",
			req: CreateReviewRequest{
				AppID:    "app-123",
				ReviewID: "review-456",
				Rating:   5,
			},
			wantErr: nil,
		},
		{
			name: "missing app_id",
			req: CreateReviewRequest{
				ReviewID: "review-456",
				Rating:   5,
			},
			wantErr: ErrAppIDRequired,
		},
		{
			name: "empty app_id",
			req: CreateReviewRequest{
				AppID:    "",
				ReviewID: "review-456",
				Rating:   5,
			},
			wantErr: ErrAppIDRequired,
		},
		{
			name: "missing review_id",
			req: CreateReviewRequest{
				AppID:  "app-123",
				Rating: 5,
			},
			wantErr: ErrReviewIDRequired,
		},
		{
			name: "empty review_id",
			req: CreateReviewRequest{
				AppID:    "app-123",
				ReviewID: "",
				Rating:   5,
			},
			wantErr: ErrReviewIDRequired,
		},
		{
			name: "rating too low",
			req: CreateReviewRequest{
				AppID:    "app-123",
				ReviewID: "review-456",
				Rating:   0,
			},
			wantErr: ErrInvalidRating,
		},
		{
			name: "rating too high",
			req: CreateReviewRequest{
				AppID:    "app-123",
				ReviewID: "review-456",
				Rating:   6,
			},
			wantErr: ErrInvalidRating,
		},
		{
			name: "negative rating",
			req: CreateReviewRequest{
				AppID:    "app-123",
				ReviewID: "review-456",
				Rating:   -1,
			},
			wantErr: ErrInvalidRating,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
