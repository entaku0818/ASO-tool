package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/entaku0818/aso-tool/backend/internal/service"
)

type PublicHandler struct {
	popularKeywordsService *service.PopularKeywordsService
	appRankingService      *service.AppRankingService
}

func NewPublicHandler(popularKeywordsService *service.PopularKeywordsService, appRankingService *service.AppRankingService) *PublicHandler {
	return &PublicHandler{
		popularKeywordsService: popularKeywordsService,
		appRankingService:      appRankingService,
	}
}

// GetPopularKeywords returns a list of popular keywords without authentication
func (h *PublicHandler) GetPopularKeywords(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	platform := r.URL.Query().Get("platform")
	limitStr := r.URL.Query().Get("limit")

	// Default values
	if country == "" {
		country = "jp"
	}
	if platform == "" {
		platform = "ios"
	}

	limit := 100
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	keywords, err := h.popularKeywordsService.GetPopularKeywords(r.Context(), country, platform, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(keywords); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetAppRanking returns App Store rankings from iTunes RSS Feed
func (h *PublicHandler) GetAppRanking(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	rankingType := r.URL.Query().Get("ranking_type")
	genreID := r.URL.Query().Get("genre_id")
	limitStr := r.URL.Query().Get("limit")

	if country == "" {
		country = "jp"
	}
	if rankingType == "" {
		rankingType = "topfreeapplications"
	}

	limit := 100
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 200 {
			limit = parsedLimit
		}
	}

	entries, err := h.appRankingService.GetRanking(r.Context(), country, rankingType, genreID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
