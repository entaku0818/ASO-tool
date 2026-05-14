package model

import "time"

type AppMetadataVersion struct {
	ID              string    `json:"id"`
	AppID           string    `json:"app_id"`
	Locale          string    `json:"locale"`
	VersionTag      string    `json:"version_tag"`
	Title           *string   `json:"title"`
	Subtitle        *string   `json:"subtitle"`
	Description     *string   `json:"description"`
	Keywords        *string   `json:"keywords"`
	PromotionalText *string   `json:"promotional_text"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UpsertMetadataRequest struct {
	AppID           string  `json:"app_id"`
	Locale          string  `json:"locale"`
	VersionTag      string  `json:"version_tag"`
	Title           *string `json:"title"`
	Subtitle        *string `json:"subtitle"`
	Description     *string `json:"description"`
	Keywords        *string `json:"keywords"`
	PromotionalText *string `json:"promotional_text"`
}

func (r *UpsertMetadataRequest) Validate() error {
	if r.AppID == "" {
		return ErrAppIDRequired
	}
	if r.Locale == "" {
		r.Locale = "ja"
	}
	if r.VersionTag == "" {
		r.VersionTag = "draft"
	}
	return nil
}
