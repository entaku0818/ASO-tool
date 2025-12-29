package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type ReviewHandler struct {
	service *service.ReviewService
}

func NewReviewHandler(service *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")

	var req model.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.AppID = appID

	review, err := h.service.Create(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, review)
}

func (h *ReviewHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "reviewID")

	review, err := h.service.Get(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, review)
}

func (h *ReviewHandler) ListByApp(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")

	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	reviews, err := h.service.ListByApp(r.Context(), appID, limit, offset)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if reviews == nil {
		reviews = []*model.Review{}
	}

	respondJSON(w, http.StatusOK, reviews)
}

func (h *ReviewHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")

	stats, err := h.service.GetStats(r.Context(), appID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "reviewID")

	if err := h.service.Delete(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
