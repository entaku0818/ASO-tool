package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/entaku0818/aso-tool/backend/internal/middleware"
	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/scraper"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type KeywordPopularityHandler struct {
	service *service.KeywordPopularityService
}

func NewKeywordPopularityHandler(svc *service.KeywordPopularityService) *KeywordPopularityHandler {
	return &KeywordPopularityHandler{service: svc}
}

// SetCredentials handles POST /apps/{appID}/search-ads-credentials
func (h *KeywordPopularityHandler) SetCredentials(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	var req model.CreateSearchAdsCredentialsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.SetCredentials(r.Context(), userID, appID, &req); err != nil {
		handleSearchAdsError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "credentials saved"})
}

// GetCredentials handles GET /apps/{appID}/search-ads-credentials
func (h *KeywordPopularityHandler) GetCredentials(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	creds, err := h.service.GetCredentials(r.Context(), userID, appID)
	if err != nil {
		handleSearchAdsError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, creds)
}

// DeleteCredentials handles DELETE /apps/{appID}/search-ads-credentials
func (h *KeywordPopularityHandler) DeleteCredentials(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	if err := h.service.DeleteCredentials(r.Context(), userID, appID); err != nil {
		handleSearchAdsError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RefreshScores handles POST /apps/{appID}/keywords/refresh-popularity
func (h *KeywordPopularityHandler) RefreshScores(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	updated, err := h.service.RefreshScores(r.Context(), userID, appID)
	if err != nil {
		handleSearchAdsError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"updated": updated})
}

// GetSuggestions handles GET /apps/{appID}/keywords/suggestions
func (h *KeywordPopularityHandler) GetSuggestions(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	limit := 25
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	suggestions, err := h.service.GetSuggestions(r.Context(), userID, appID, limit)
	if err != nil {
		handleSearchAdsError(w, err)
		return
	}

	if suggestions == nil {
		suggestions = []scraper.KeywordPopularity{}
	}

	respondJSON(w, http.StatusOK, suggestions)
}

func handleSearchAdsError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		respondError(w, http.StatusNotFound, "app not found")
	case errors.Is(err, model.ErrSearchAdsCredentialsNotFound):
		respondError(w, http.StatusNotFound, "search ads credentials not found")
	case errors.Is(err, model.ErrInvalidPrivateKey):
		respondError(w, http.StatusBadRequest, "invalid private key format")
	case errors.Is(err, model.ErrClientIDRequired),
		errors.Is(err, model.ErrTeamIDRequired),
		errors.Is(err, model.ErrKeyIDRequired),
		errors.Is(err, model.ErrPrivateKeyRequired),
		errors.Is(err, model.ErrAdamIDRequired):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, model.ErrInvalidSearchAdsCredentials):
		respondError(w, http.StatusUnauthorized, "invalid Search Ads credentials")
	default:
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
