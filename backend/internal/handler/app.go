package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/entaku0818/aso-tool/backend/internal/middleware"
	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type AppHandler struct {
	service *service.AppService
}

func NewAppHandler(service *service.AppService) *AppHandler {
	return &AppHandler{service: service}
}

func (h *AppHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req model.CreateAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	app, err := h.service.Create(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, model.ErrAppLimitExceeded) {
			respondError(w, http.StatusPaymentRequired, err.Error())
			return
		}
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, app)
}

func (h *AppHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	app, err := h.service.Get(r.Context(), userID, id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, app)
}

func (h *AppHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	apps, err := h.service.List(r.Context(), userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if apps == nil {
		apps = []*model.App{}
	}

	respondJSON(w, http.StatusOK, apps)
}

func (h *AppHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	var req model.UpdateAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	app, err := h.service.Update(r.Context(), userID, id, &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, app)
}

func (h *AppHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	if err := h.service.Delete(r.Context(), userID, id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		respondError(w, http.StatusNotFound, "not found")
	case errors.Is(err, model.ErrNameRequired),
		errors.Is(err, model.ErrBundleIDRequired),
		errors.Is(err, model.ErrInvalidPlatform),
		errors.Is(err, model.ErrAppIDRequired),
		errors.Is(err, model.ErrKeywordRequired),
		errors.Is(err, model.ErrKeywordIDRequired),
		errors.Is(err, model.ErrReviewIDRequired),
		errors.Is(err, model.ErrInvalidRating):
		respondError(w, http.StatusBadRequest, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
