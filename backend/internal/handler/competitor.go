package handler

import (
	"encoding/json"
	"net/http"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type CompetitorHandler struct {
	service *service.CompetitorService
}

func NewCompetitorHandler(service *service.CompetitorService) *CompetitorHandler {
	return &CompetitorHandler{service: service}
}

func (h *CompetitorHandler) Create(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")

	var req model.CreateCompetitorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.AppID = appID

	competitor, err := h.service.Create(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, competitor)
}

func (h *CompetitorHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "competitorID")

	competitor, err := h.service.Get(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, competitor)
}

func (h *CompetitorHandler) List(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")

	competitors, err := h.service.ListByApp(r.Context(), appID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if competitors == nil {
		competitors = []*model.Competitor{}
	}

	respondJSON(w, http.StatusOK, competitors)
}

func (h *CompetitorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "competitorID")

	if err := h.service.Delete(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CompetitorHandler) UpdateRankings(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")

	count, err := h.service.UpdateCompetitorRankings(r.Context(), appID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"updated": count})
}

func (h *CompetitorHandler) GetComparison(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")
	keywordID := chi.URLParam(r, "keywordID")

	comparison, err := h.service.GetComparison(r.Context(), appID, keywordID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, comparison)
}
