package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/entaku0818/aso-tool/backend/internal/repository"
	"github.com/entaku0818/aso-tool/backend/internal/scraper"
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

// GetKeywordSuggestions proxies Apple App Store search hints API.
// Query params: term (required), country (default: jp)
func (h *PublicHandler) GetKeywordSuggestions(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	if term == "" {
		http.Error(w, "term is required", http.StatusBadRequest)
		return
	}
	country := r.URL.Query().Get("country")
	if country == "" {
		country = "jp"
	}

	suggestions, err := scraper.FetchKeywordSuggestions(r.Context(), term, country)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if suggestions == nil {
		suggestions = []scraper.KeywordSuggestion{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(suggestions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// GetAppRankingTrend returns daily rank history for a specific app.
// Query params: app_id (required), country (default: jp), ranking_type (default: topfreeapplications), days (default: 30)
func (h *PublicHandler) GetAppRankingTrend(w http.ResponseWriter, r *http.Request) {
	appID := r.URL.Query().Get("app_id")
	country := r.URL.Query().Get("country")
	rankingType := r.URL.Query().Get("ranking_type")
	daysStr := r.URL.Query().Get("days")

	if appID == "" {
		http.Error(w, "app_id is required", http.StatusBadRequest)
		return
	}
	if country == "" {
		country = "jp"
	}
	if rankingType == "" {
		rankingType = "topfreeapplications"
	}
	days := 30
	if daysStr != "" {
		if parsed, err := strconv.Atoi(daysStr); err == nil && parsed > 0 {
			days = parsed
		}
	}

	points, err := h.appRankingService.GetRankingTrend(r.Context(), appID, country, rankingType, days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if points == nil {
		points = []repository.RankingTrendPoint{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(points); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetAppRankingCountries returns the latest rank for a specific app across multiple countries.
// Query params: app_id (required), ranking_type (default: topfreeapplications)
func (h *PublicHandler) GetAppRankingCountries(w http.ResponseWriter, r *http.Request) {
	appID := r.URL.Query().Get("app_id")
	rankingType := r.URL.Query().Get("ranking_type")

	if appID == "" {
		http.Error(w, "app_id is required", http.StatusBadRequest)
		return
	}
	if rankingType == "" {
		rankingType = "topfreeapplications"
	}

	points, err := h.appRankingService.GetCountryComparison(r.Context(), appID, rankingType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if points == nil {
		points = []repository.CountryRankPoint{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(points); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
