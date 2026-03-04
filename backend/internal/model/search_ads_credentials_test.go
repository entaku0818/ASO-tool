package model

import (
	"errors"
	"testing"
)

func TestCreateSearchAdsCredentialsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateSearchAdsCredentialsRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: CreateSearchAdsCredentialsRequest{
				ClientID:   "client-id",
				TeamID:     "team-id",
				KeyID:      "key-id",
				PrivateKey: "base64-encoded-key",
				AdamID:     123456789,
			},
			wantErr: nil,
		},
		{
			name: "valid request with org_id",
			req: CreateSearchAdsCredentialsRequest{
				ClientID:   "client-id",
				TeamID:     "team-id",
				KeyID:      "key-id",
				PrivateKey: "base64-encoded-key",
				OrgID:      "org-id",
				AdamID:     123456789,
			},
			wantErr: nil,
		},
		{
			name: "missing client_id",
			req: CreateSearchAdsCredentialsRequest{
				TeamID:     "team-id",
				KeyID:      "key-id",
				PrivateKey: "base64-encoded-key",
				AdamID:     123456789,
			},
			wantErr: ErrClientIDRequired,
		},
		{
			name: "empty client_id",
			req: CreateSearchAdsCredentialsRequest{
				ClientID:   "",
				TeamID:     "team-id",
				KeyID:      "key-id",
				PrivateKey: "base64-encoded-key",
				AdamID:     123456789,
			},
			wantErr: ErrClientIDRequired,
		},
		{
			name: "missing team_id",
			req: CreateSearchAdsCredentialsRequest{
				ClientID:   "client-id",
				KeyID:      "key-id",
				PrivateKey: "base64-encoded-key",
				AdamID:     123456789,
			},
			wantErr: ErrTeamIDRequired,
		},
		{
			name: "missing key_id",
			req: CreateSearchAdsCredentialsRequest{
				ClientID:   "client-id",
				TeamID:     "team-id",
				PrivateKey: "base64-encoded-key",
				AdamID:     123456789,
			},
			wantErr: ErrKeyIDRequired,
		},
		{
			name: "missing private_key",
			req: CreateSearchAdsCredentialsRequest{
				ClientID: "client-id",
				TeamID:   "team-id",
				KeyID:    "key-id",
				AdamID:   123456789,
			},
			wantErr: ErrPrivateKeyRequired,
		},
		{
			name: "adam_id is zero",
			req: CreateSearchAdsCredentialsRequest{
				ClientID:   "client-id",
				TeamID:     "team-id",
				KeyID:      "key-id",
				PrivateKey: "base64-encoded-key",
				AdamID:     0,
			},
			wantErr: ErrAdamIDRequired,
		},
		{
			name: "adam_id is negative",
			req: CreateSearchAdsCredentialsRequest{
				ClientID:   "client-id",
				TeamID:     "team-id",
				KeyID:      "key-id",
				PrivateKey: "base64-encoded-key",
				AdamID:     -1,
			},
			wantErr: ErrAdamIDRequired,
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

func TestSearchAdsCredentials_Fields(t *testing.T) {
	creds := SearchAdsCredentials{
		ID:       "cred-id",
		AppID:    "app-id",
		ClientID: "client-id",
		TeamID:   "team-id",
		KeyID:    "key-id",
		OrgID:    "org-id",
		AdamID:   123456789,
		IsValid:  true,
	}

	if creds.ID != "cred-id" {
		t.Errorf("ID = %v, want 'cred-id'", creds.ID)
	}
	if creds.AdamID != 123456789 {
		t.Errorf("AdamID = %v, want 123456789", creds.AdamID)
	}
	if !creds.IsValid {
		t.Error("IsValid = false, want true")
	}
}
