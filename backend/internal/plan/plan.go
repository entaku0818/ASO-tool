// Package plan is the single source of truth for free-plan usage limits.
// Actual prices/Price IDs live in Stripe (see service.BillingService), which
// is configured via STRIPE_PRICE_MONTHLY/YEARLY/LICENSE env vars, not here.
package plan

const (
	// FreeAppLimit is the number of apps a free-plan user may register.
	FreeAppLimit = 1

	// FreeKeywordLimit is the number of keywords a free-plan user may track per app.
	FreeKeywordLimit = 10
)
