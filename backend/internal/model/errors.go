package model

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrNameRequired     = errors.New("name is required")
	ErrBundleIDRequired = errors.New("bundle_id is required")
	ErrInvalidPlatform  = errors.New("platform must be 'ios' or 'android'")
	ErrAlreadyExists    = errors.New("app already exists")
)
