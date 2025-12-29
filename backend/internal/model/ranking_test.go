package model

import (
	"errors"
	"testing"
)

func TestCreateRankingRequest_Validate(t *testing.T) {
	rank := 5

	tests := []struct {
		name    string
		req     CreateRankingRequest
		wantErr error
	}{
		{
			name: "valid request with rank",
			req: CreateRankingRequest{
				KeywordID: "keyword-123",
				Rank:      &rank,
			},
			wantErr: nil,
		},
		{
			name: "valid request without rank (not found)",
			req: CreateRankingRequest{
				KeywordID: "keyword-123",
				Rank:      nil,
			},
			wantErr: nil,
		},
		{
			name: "missing keyword_id",
			req: CreateRankingRequest{
				Rank: &rank,
			},
			wantErr: ErrKeywordIDRequired,
		},
		{
			name: "empty keyword_id",
			req: CreateRankingRequest{
				KeywordID: "",
				Rank:      &rank,
			},
			wantErr: ErrKeywordIDRequired,
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
