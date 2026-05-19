package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/middleware"
	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type AnalyticsHandler struct {
	service *service.AnalyticsService
}

func NewAnalyticsHandler(service *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

// SetCredentials handles POST /apps/{appID}/asc-credentials
func (h *AnalyticsHandler) SetCredentials(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	var req model.CreateASCCredentialsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.SetCredentials(r.Context(), userID, appID, &req); err != nil {
		handleAnalyticsError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "credentials saved"})
}

// GetCredentials handles GET /apps/{appID}/asc-credentials
func (h *AnalyticsHandler) GetCredentials(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	creds, err := h.service.GetCredentials(r.Context(), userID, appID)
	if err != nil {
		handleAnalyticsError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, creds)
}

// DeleteCredentials handles DELETE /apps/{appID}/asc-credentials
func (h *AnalyticsHandler) DeleteCredentials(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	if err := h.service.DeleteCredentials(r.Context(), userID, appID); err != nil {
		handleAnalyticsError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ValidateCredentials handles POST /apps/{appID}/asc-credentials/validate
func (h *AnalyticsHandler) ValidateCredentials(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	if err := h.service.ValidateCredentials(r.Context(), userID, appID); err != nil {
		if errors.Is(err, model.ErrInvalidASCCredentials) {
			respondJSON(w, http.StatusOK, map[string]interface{}{
				"valid":   false,
				"message": "Invalid credentials",
			})
			return
		}
		handleAnalyticsError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"valid":   true,
		"message": "Credentials are valid",
	})
}

// GetAnalytics handles GET /apps/{appID}/analytics
func (h *AnalyticsHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	days := 30
	if d := r.URL.Query().Get("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	analytics, err := h.service.GetAnalytics(r.Context(), userID, appID, days)
	if err != nil {
		handleAnalyticsError(w, err)
		return
	}

	if analytics == nil {
		analytics = []*model.ASCAnalytics{}
	}

	respondJSON(w, http.StatusOK, analytics)
}

// GetAnalyticsWithCorrelation handles GET /apps/{appID}/analytics/correlation
func (h *AnalyticsHandler) GetAnalyticsWithCorrelation(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	days := 30
	if d := r.URL.Query().Get("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	analytics, err := h.service.GetAnalyticsWithCorrelation(r.Context(), userID, appID, days)
	if err != nil {
		handleAnalyticsError(w, err)
		return
	}

	if analytics == nil {
		analytics = []*model.AnalyticsWithCorrelation{}
	}

	respondJSON(w, http.StatusOK, analytics)
}

// GetSummary handles GET /apps/{appID}/analytics/summary
func (h *AnalyticsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	days := 30
	if d := r.URL.Query().Get("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	summary, err := h.service.GetSummary(r.Context(), userID, appID, days)
	if err != nil {
		handleAnalyticsError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, summary)
}

// TriggerFetch handles POST /apps/{appID}/analytics/fetch
func (h *AnalyticsHandler) TriggerFetch(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	// Verify app belongs to user
	_, err := h.service.GetCredentials(r.Context(), userID, appID)
	if err != nil {
		handleAnalyticsError(w, err)
		return
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	count, err := h.service.FetchAnalytics(r.Context(), appID, startDate, endDate)
	if err != nil {
		handleAnalyticsError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":     "Analytics fetch completed",
		"days_stored": count,
	})
}

// FetchSearchKeywords handles POST /apps/{appID}/analytics/search-keywords/fetch
func (h *AnalyticsHandler) FetchSearchKeywords(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	req, err := h.service.FetchSearchKeywords(r.Context(), userID, appID)
	if err != nil {
		handleAnalyticsError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, req)
}

// PollSearchKeywords handles POST /apps/{appID}/analytics/search-keywords/poll
func (h *AnalyticsHandler) PollSearchKeywords(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	req, err := h.service.PollSearchKeywords(r.Context(), userID, appID)
	if err != nil {
		handleAnalyticsError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, req)
}

// GetSearchKeywords handles GET /apps/{appID}/analytics/search-keywords
func (h *AnalyticsHandler) GetSearchKeywords(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	days := 90
	if d := r.URL.Query().Get("days"); d != "" {
		if n, err := strconv.Atoi(d); err == nil && n > 0 {
			days = n
		}
	}

	keywords, err := h.service.GetSearchKeywords(r.Context(), userID, appID, days)
	if err != nil {
		handleAnalyticsError(w, err)
		return
	}

	if keywords == nil {
		keywords = []*model.ASCSearchKeyword{}
	}
	respondJSON(w, http.StatusOK, keywords)
}

func handleAnalyticsError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		respondError(w, http.StatusNotFound, "app not found")
	case errors.Is(err, model.ErrCredentialsNotFound):
		respondError(w, http.StatusNotFound, "credentials not found")
	case errors.Is(err, model.ErrIOSOnly):
		respondError(w, http.StatusBadRequest, "analytics only available for iOS apps")
	case errors.Is(err, model.ErrInvalidPrivateKey):
		respondError(w, http.StatusBadRequest, "invalid private key format")
	case errors.Is(err, model.ErrIssuerIDRequired),
		errors.Is(err, model.ErrKeyIDRequired),
		errors.Is(err, model.ErrPrivateKeyRequired):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, model.ErrInvalidASCCredentials):
		respondError(w, http.StatusUnauthorized, "invalid App Store Connect credentials")
	default:
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
