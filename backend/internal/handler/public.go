package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/entaku0818/aso-tool/backend/internal/service"
)

type PublicHandler struct {
	popularKeywordsService *service.PopularKeywordsService
}

func NewPublicHandler(popularKeywordsService *service.PopularKeywordsService) *PublicHandler {
	return &PublicHandler{
		popularKeywordsService: popularKeywordsService,
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
