package model

import "time"

type ASCCredentials struct {
	ID                  string     `json:"id"`
	AppID               string     `json:"app_id"`
	IssuerID            string     `json:"issuer_id"`
	KeyID               string     `json:"key_id"`
	PrivateKeyEncrypted []byte     `json:"-"`
	EncryptionNonce     []byte     `json:"-"`
	IsValid             bool       `json:"is_valid"`
	LastValidatedAt     *time.Time `json:"last_validated_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type CreateASCCredentialsRequest struct {
	IssuerID   string `json:"issuer_id"`
	KeyID      string `json:"key_id"`
	PrivateKey string `json:"private_key"` // Base64 encoded .p8 content
}

func (r *CreateASCCredentialsRequest) Validate() error {
	if r.IssuerID == "" {
		return ErrIssuerIDRequired
	}
	if r.KeyID == "" {
		return ErrKeyIDRequired
	}
	if r.PrivateKey == "" {
		return ErrPrivateKeyRequired
	}
	return nil
}
