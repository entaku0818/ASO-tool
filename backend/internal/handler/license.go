package handler

import (
	"encoding/json"
	"net/http"

	"github.com/entaku0818/aso-compass/backend/internal/middleware"
	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/entaku0818/aso-compass/backend/internal/service"
)

type LicenseHandler struct {
	service *service.LicenseService
}

func NewLicenseHandler(service *service.LicenseService) *LicenseHandler {
	return &LicenseHandler{service: service}
}

// Activate validates a license key and returns a JWT token.
// POST /api/licenses/activate
func (h *LicenseHandler) Activate(w http.ResponseWriter, r *http.Request) {
	var req model.ActivateLicenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.service.Activate(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// Generate creates a new license key (admin only or called from Stripe webhook).
// POST /api/licenses/generate
func (h *LicenseHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var req model.GenerateLicenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	lk, err := h.service.Generate(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, lk)
}

// GetMyLicense returns the active license for the authenticated user.
// GET /api/licenses/me
func (h *LicenseHandler) GetMyLicense(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	lk, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, lk)
}

// CreateCheckout starts a Stripe Checkout session for a one-time license purchase.
// POST /api/licenses/checkout
func (h *LicenseHandler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Email == "" {
		respondError(w, http.StatusBadRequest, "email is required")
		return
	}

	url, err := h.service.CreateCheckoutSession(r.Context(), body.Email)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"url": url})
}

// ListAll returns all license keys (admin only).
// GET /api/admin/licenses
func (h *LicenseHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	keys, err := h.service.List(r.Context())
	if err != nil {
		handleServiceError(w, err)
		return
	}
	if keys == nil {
		keys = []*model.LicenseKey{}
	}
	respondJSON(w, http.StatusOK, keys)
}
