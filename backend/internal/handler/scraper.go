package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type ScraperHandler struct {
	service *service.ScraperService
}

func NewScraperHandler(service *service.ScraperService) *ScraperHandler {
	return &ScraperHandler{service: service}
}

// FetchAppInfo fetches app information from store
func (h *ScraperHandler) FetchAppInfo(w http.ResponseWriter, r *http.Request) {
	bundleID := r.URL.Query().Get("bundle_id")
	platform := model.Platform(r.URL.Query().Get("platform"))
	country := r.URL.Query().Get("country")

	if bundleID == "" {
		respondError(w, http.StatusBadRequest, "bundle_id is required")
		return
	}

	if platform != model.PlatformIOS && platform != model.PlatformAndroid {
		respondError(w, http.StatusBadRequest, "platform must be 'ios' or 'android'")
		return
	}

	appInfo, err := h.service.FetchAppInfo(r.Context(), bundleID, platform, country)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, appInfo)
}

// UpdateRankings triggers ranking update for an app
func (h *ScraperHandler) UpdateRankings(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")

	count, err := h.service.UpdateKeywordRankings(r.Context(), appID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":         "rankings updated",
		"keywords_updated": count,
	})
}

// FetchReviews triggers review fetch for an app
func (h *ScraperHandler) FetchReviews(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")
	itunesID := r.URL.Query().Get("itunes_id")

	count, err := h.service.FetchReviews(r.Context(), appID, itunesID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":        "reviews fetched",
		"reviews_stored": count,
	})
}

// SearchApps searches for apps in the store
func (h *ScraperHandler) SearchApps(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	platform := model.Platform(r.URL.Query().Get("platform"))
	country := r.URL.Query().Get("country")
	limitStr := r.URL.Query().Get("limit")

	if keyword == "" {
		respondError(w, http.StatusBadRequest, "keyword is required")
		return
	}

	if platform != model.PlatformIOS && platform != model.PlatformAndroid {
		respondError(w, http.StatusBadRequest, "platform must be 'ios' or 'android'")
		return
	}

	limit := 50
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	results, err := h.service.SearchApps(r.Context(), keyword, platform, country, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, results)
}

// TriggerAllUpdates triggers updates for all apps (for scheduler)
func (h *ScraperHandler) TriggerAllUpdates(w http.ResponseWriter, r *http.Request) {
	// This endpoint can be called by Cloud Scheduler
	// In a production environment, you'd want to add authentication

	var req struct {
		UpdateRankings bool   `json:"update_rankings"`
		FetchReviews   bool   `json:"fetch_reviews"`
		ITunesIDs      map[string]string `json:"itunes_ids"` // app_id -> itunes_id mapping
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// If no body, update rankings by default
		req.UpdateRankings = true
	}

	// Note: In a real implementation, you'd iterate through all apps
	// and update rankings/reviews for each one

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "update triggered",
	})
}
