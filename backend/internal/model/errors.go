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

	// Competitor errors
	ErrCompetitorIDRequired = errors.New("competitor_id is required")

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

	// Screenshot errors
	ErrURLRequired         = errors.New("url is required")
	ErrInvalidDeviceType   = errors.New("invalid device type")
	ErrInvalidDisplayOrder = errors.New("display_order must be non-negative")

	// ASC Credentials errors
	ErrIssuerIDRequired      = errors.New("issuer_id is required")
	ErrKeyIDRequired         = errors.New("key_id is required")
	ErrPrivateKeyRequired    = errors.New("private_key is required")
	ErrInvalidPrivateKey     = errors.New("invalid private key format")
	ErrCredentialsNotFound   = errors.New("credentials not found")
	ErrInvalidASCCredentials = errors.New("invalid App Store Connect credentials")

	// Analytics errors
	ErrAnalyticsNotAvailable = errors.New("analytics not available for this app")
	ErrIOSOnly               = errors.New("analytics is only available for iOS apps")
	ErrVersionRequired       = errors.New("version is required")
	ErrReleaseDateRequired   = errors.New("release_date is required")
)
