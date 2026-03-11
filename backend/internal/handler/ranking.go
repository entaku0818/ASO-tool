package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/entaku0818/aso-tool/backend/internal/middleware"
	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type RankingHandler struct {
	service    *service.RankingService
	appService *service.AppService
}

func NewRankingHandler(service *service.RankingService, appService *service.AppService) *RankingHandler {
	return &RankingHandler{service: service, appService: appService}
}

func (h *RankingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateRankingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ranking, err := h.service.Create(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, ranking)
}

func (h *RankingHandler) ListByKeyword(w http.ResponseWriter, r *http.Request) {
	keywordID := chi.URLParam(r, "keywordID")

	limit := 30
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	rankings, err := h.service.ListByKeyword(r.Context(), keywordID, limit)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if rankings == nil {
		rankings = []*model.RankingHistory{}
	}

	respondJSON(w, http.StatusOK, rankings)
}

func (h *RankingHandler) ListByApp(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")

	days := 30
	if d := r.URL.Query().Get("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	rankings, err := h.service.ListByApp(r.Context(), appID, days)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if rankings == nil {
		rankings = []*model.RankingWithKeyword{}
	}

	respondJSON(w, http.StatusOK, rankings)
}

func (h *RankingHandler) GetLatest(w http.ResponseWriter, r *http.Request) {
	keywordID := chi.URLParam(r, "keywordID")

	ranking, err := h.service.GetLatestByKeyword(r.Context(), keywordID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, ranking)
}

// CompareWithYesterday returns keyword rankings for today vs yesterday for an app specified by name.
// Query param: name (app name, partial match supported)
func (h *RankingHandler) CompareWithYesterday(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	name := r.URL.Query().Get("name")
	if name == "" {
		respondError(w, http.StatusBadRequest, "name query parameter is required")
		return
	}

	apps, err := h.appService.FindByName(r.Context(), userID, name)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	if len(apps) == 0 {
		respondError(w, http.StatusNotFound, "no app found with given name")
		return
	}

	comparisons, err := h.service.CompareWithYesterday(r.Context(), apps[0].ID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if comparisons == nil {
		comparisons = []repository.DailyRankingComparison{}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"app":         apps[0],
		"comparisons": comparisons,
	})
}

func (h *RankingHandler) GetRisingKeywords(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")

	days := 7
	if d := r.URL.Query().Get("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	keywords, err := h.service.GetRisingKeywords(r.Context(), appID, days)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if keywords == nil {
		keywords = []repository.RisingKeyword{}
	}

	respondJSON(w, http.StatusOK, keywords)
}
