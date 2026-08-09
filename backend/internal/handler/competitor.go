package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/entaku0818/aso-compass/backend/internal/middleware"
	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/entaku0818/aso-compass/backend/internal/repository"
	"github.com/entaku0818/aso-compass/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type CompetitorHandler struct {
	service  *service.CompetitorService
	userRepo *repository.UserRepository
}

func NewCompetitorHandler(service *service.CompetitorService, userRepo *repository.UserRepository) *CompetitorHandler {
	return &CompetitorHandler{service: service, userRepo: userRepo}
}

// requirePro returns false (and has already written an error response) if the
// requesting user is not on the Pro plan.
func (h *CompetitorHandler) requirePro(w http.ResponseWriter, r *http.Request) bool {
	userID := middleware.GetUserID(r.Context())
	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			respondError(w, http.StatusUnauthorized, "user not found")
		} else {
			respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return false
	}
	if !user.IsPro() {
		respondError(w, http.StatusForbidden, "keyword gap analysis requires Pro plan")
		return false
	}
	return true
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

// GetKeywordGap returns keywords where a competitor outranks the app.
// GET /api/apps/{appID}/competitors/keyword-gap
func (h *CompetitorHandler) GetKeywordGap(w http.ResponseWriter, r *http.Request) {
	if !h.requirePro(w, r) {
		return
	}

	appID := chi.URLParam(r, "appID")

	gaps, err := h.service.GetKeywordGap(r.Context(), appID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if gaps == nil {
		gaps = []*model.KeywordGap{}
	}

	respondJSON(w, http.StatusOK, gaps)
}
