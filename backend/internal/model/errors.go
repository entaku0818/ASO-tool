package model

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrNameRequired     = errors.New("name is required")
	ErrBundleIDRequired = errors.New("bundle_id is required")
	ErrInvalidPlatform  = errors.New("platform must be 'ios' or 'android'")
	ErrAlreadyExists    = errors.New("already exists")

	// Keyword errors
	ErrAppIDRequired     = errors.New("app_id is required")
	ErrKeywordRequired   = errors.New("keyword is required")
	ErrKeywordIDRequired = errors.New("keyword_id is required")

	// Review errors
	ErrReviewIDRequired = errors.New("review_id is required")
	ErrInvalidRating    = errors.New("rating must be between 1 and 5")

	// Auth errors
	ErrEmailRequired     = errors.New("email is required")
	ErrPasswordRequired  = errors.New("password is required")
	ErrPasswordTooShort  = errors.New("password must be at least 8 characters")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
)
