package model

import "time"

type SearchAdsCredentials struct {
	ID                  string     `json:"id"`
	AppID               string     `json:"app_id"`
	ClientID            string     `json:"client_id"`
	TeamID              string     `json:"team_id"`
	KeyID               string     `json:"key_id"`
	PrivateKeyEncrypted []byte     `json:"-"`
	EncryptionNonce     []byte     `json:"-"`
	OrgID               string     `json:"org_id,omitempty"`
	AdamID              int64      `json:"adam_id"`
	IsValid             bool       `json:"is_valid"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type CreateSearchAdsCredentialsRequest struct {
	ClientID   string `json:"client_id"`
	TeamID     string `json:"team_id"`
	KeyID      string `json:"key_id"`
	PrivateKey string `json:"private_key"` // Base64 encoded .p8 content
	OrgID      string `json:"org_id,omitempty"`
	AdamID     int64  `json:"adam_id"`
}

func (r *CreateSearchAdsCredentialsRequest) Validate() error {
	if r.ClientID == "" {
		return ErrClientIDRequired
	}
	if r.TeamID == "" {
		return ErrTeamIDRequired
	}
	if r.KeyID == "" {
		return ErrKeyIDRequired
	}
	if r.PrivateKey == "" {
		return ErrPrivateKeyRequired
	}
	if r.AdamID <= 0 {
		return ErrAdamIDRequired
	}
	return nil
}
