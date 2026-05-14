package model

import "time"

type LicenseKey struct {
	ID              string     `json:"id"`
	Key             string     `json:"key"`
	Email           string     `json:"email"`
	UserID          *string    `json:"user_id,omitempty"`
	IsActive        bool       `json:"is_active"`
	ActivatedAt     *time.Time `json:"activated_at,omitempty"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	StripeSessionID *string    `json:"stripe_session_id,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type ActivateLicenseRequest struct {
	Key   string `json:"key"`
	Email string `json:"email"`
}

func (r *ActivateLicenseRequest) Validate() error {
	if r.Key == "" {
		return ErrLicenseKeyRequired
	}
	if r.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

type ActivateLicenseResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
	Key   string `json:"key"`
}

type GenerateLicenseRequest struct {
	Email           string  `json:"email"`
	StripeSessionID *string `json:"stripe_session_id,omitempty"`
}
