package model

import (
	"time"
)

type DeviceType string

const (
	DeviceTypeIPhone    DeviceType = "iphone"
	DeviceTypeIPad      DeviceType = "ipad"
	DeviceTypeAndroid   DeviceType = "android"
	DeviceTypeTablet    DeviceType = "tablet"
)

func (d DeviceType) IsValid() bool {
	switch d {
	case DeviceTypeIPhone, DeviceTypeIPad, DeviceTypeAndroid, DeviceTypeTablet:
		return true
	}
	return false
}

type Screenshot struct {
	ID           string     `json:"id"`
	AppID        string     `json:"app_id"`
	URL          string     `json:"url"`
	DisplayOrder int        `json:"display_order"`
	DeviceType   DeviceType `json:"device_type"`
	Locale       string     `json:"locale"`
	Description  string     `json:"description,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CreateScreenshotRequest struct {
	AppID        string     `json:"app_id"`
	URL          string     `json:"url"`
	DisplayOrder int        `json:"display_order"`
	DeviceType   DeviceType `json:"device_type"`
	Locale       string     `json:"locale"`
	Description  string     `json:"description,omitempty"`
}

func (r *CreateScreenshotRequest) Validate() error {
	if r.AppID == "" {
		return ErrAppIDRequired
	}
	if r.URL == "" {
		return ErrURLRequired
	}
	if r.DisplayOrder < 0 {
		return ErrInvalidDisplayOrder
	}
	if r.DeviceType != "" && !r.DeviceType.IsValid() {
		return ErrInvalidDeviceType
	}
	return nil
}

type UpdateScreenshotRequest struct {
	URL          *string     `json:"url,omitempty"`
	DisplayOrder *int        `json:"display_order,omitempty"`
	DeviceType   *DeviceType `json:"device_type,omitempty"`
	Locale       *string     `json:"locale,omitempty"`
	Description  *string     `json:"description,omitempty"`
}

func (r *UpdateScreenshotRequest) Validate() error {
	if r.DisplayOrder != nil && *r.DisplayOrder < 0 {
		return ErrInvalidDisplayOrder
	}
	if r.DeviceType != nil && !r.DeviceType.IsValid() {
		return ErrInvalidDeviceType
	}
	return nil
}

type ReorderScreenshotsRequest struct {
	ScreenshotIDs []string `json:"screenshot_ids"`
}

func (r *ReorderScreenshotsRequest) Validate() error {
	if len(r.ScreenshotIDs) == 0 {
		return ErrURLRequired
	}
	return nil
}
