package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/entaku0818/aso-tool/backend/internal/middleware"
	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/service"
)

type BillingHandler struct {
	service *service.BillingService
}

func NewBillingHandler(service *service.BillingService) *BillingHandler {
	return &BillingHandler{service: service}
}

// CreateCheckout handles POST /api/billing/checkout
// Body: {"plan_type": "monthly" | "yearly"}
func (h *BillingHandler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req struct {
		PlanType string `json:"plan_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.PlanType != "monthly" && req.PlanType != "yearly" {
		respondError(w, http.StatusBadRequest, "plan_type must be 'monthly' or 'yearly'")
		return
	}

	url, err := h.service.CreateCheckoutSession(r.Context(), userID, req.PlanType)
	if err != nil {
		if errors.Is(err, model.ErrStripeKeyRequired) {
			respondError(w, http.StatusServiceUnavailable, "billing is not configured")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to create checkout session")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"url": url})
}

// Webhook handles POST /api/billing/webhook (public — no JWT, Stripe signature verified)
// Signature validation is delegated to the service layer, which skips it when
// STRIPE_WEBHOOK_SECRET is not set (local dev only).
func (h *BillingHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	signature := r.Header.Get("Stripe-Signature")

	if err := h.service.HandleWebhook(r.Context(), r.Body, signature); err != nil {
		// Return 400 so Stripe retries on signature/parse errors, 200 on unknown events
		respondError(w, http.StatusBadRequest, "webhook processing failed")
		return
	}

	w.WriteHeader(http.StatusOK)
}
