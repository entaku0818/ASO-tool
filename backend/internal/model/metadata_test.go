package model

import (
	"errors"
	"testing"
)

func TestUpsertMetadataRequest_Validate(t *testing.T) {
	tests := []struct {
		name        string
		req         UpsertMetadataRequest
		wantErr     error
		wantLocale  string
		wantVersion string
	}{
		{
			name:        "valid full request",
			req:         UpsertMetadataRequest{AppID: "app-1", Locale: "ja", VersionTag: "v1.0"},
			wantErr:     nil,
			wantLocale:  "ja",
			wantVersion: "v1.0",
		},
		{
			name:        "defaults locale and version_tag",
			req:         UpsertMetadataRequest{AppID: "app-1"},
			wantErr:     nil,
			wantLocale:  "ja",
			wantVersion: "draft",
		},
		{
			name:    "missing app_id",
			req:     UpsertMetadataRequest{Locale: "ja"},
			wantErr: ErrAppIDRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				if tt.req.Locale != tt.wantLocale {
					t.Errorf("Locale = %q, want %q", tt.req.Locale, tt.wantLocale)
				}
				if tt.req.VersionTag != tt.wantVersion {
					t.Errorf("VersionTag = %q, want %q", tt.req.VersionTag, tt.wantVersion)
				}
			}
		})
	}
}
