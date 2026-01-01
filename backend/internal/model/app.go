package model

import "time"

type Platform string

const (
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
)

type App struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	BundleID  string    `json:"bundle_id"`
	Platform  Platform  `json:"platform"`
	StoreURL  string    `json:"store_url,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAppRequest struct {
	Name     string   `json:"name"`
	BundleID string   `json:"bundle_id"`
	Platform Platform `json:"platform"`
	StoreURL string   `json:"store_url,omitempty"`
}

type UpdateAppRequest struct {
	Name     string `json:"name,omitempty"`
	StoreURL string `json:"store_url,omitempty"`
}

func (r *CreateAppRequest) Validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}
	if r.BundleID == "" {
		return ErrBundleIDRequired
	}
	if r.Platform != PlatformIOS && r.Platform != PlatformAndroid {
		return ErrInvalidPlatform
	}
	return nil
}
