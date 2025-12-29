package model

import (
	"errors"
	"testing"
)

func TestCreateKeywordRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateKeywordRequest
		wantErr error
	}{
		{
			name: "valid request with country",
			req: CreateKeywordRequest{
				AppID:   "app-123",
				Keyword: "calculator",
				Country: "jp",
			},
			wantErr: nil,
		},
		{
			name: "valid request without country (defaults to JP)",
			req: CreateKeywordRequest{
				AppID:   "app-123",
				Keyword: "calculator",
			},
			wantErr: nil,
		},
		{
			name: "missing app_id",
			req: CreateKeywordRequest{
				Keyword: "calculator",
				Country: "jp",
			},
			wantErr: ErrAppIDRequired,
		},
		{
			name: "empty app_id",
			req: CreateKeywordRequest{
				AppID:   "",
				Keyword: "calculator",
				Country: "jp",
			},
			wantErr: ErrAppIDRequired,
		},
		{
			name: "missing keyword",
			req: CreateKeywordRequest{
				AppID:   "app-123",
				Country: "jp",
			},
			wantErr: ErrKeywordRequired,
		},
		{
			name: "empty keyword",
			req: CreateKeywordRequest{
				AppID:   "app-123",
				Keyword: "",
				Country: "jp",
			},
			wantErr: ErrKeywordRequired,
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

func TestCreateKeywordRequest_DefaultCountry(t *testing.T) {
	req := &CreateKeywordRequest{
		AppID:   "app-123",
		Keyword: "test",
	}

	err := req.Validate()
	if err != nil {
		t.Errorf("Validate() unexpected error = %v", err)
	}

	if req.Country != "JP" {
		t.Errorf("Country = %v, want 'JP'", req.Country)
	}
}
