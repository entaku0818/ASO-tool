package handler

import (
	"encoding/json"
	"net/http"

	"github.com/entaku0818/aso-tool/backend/internal/middleware"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type TrackedKeywordHandler struct {
	service *service.TrackedKeywordService
}

func NewTrackedKeywordHandler(service *service.TrackedKeywordService) *TrackedKeywordHandler {
	return &TrackedKeywordHandler{service: service}
}

func (h *TrackedKeywordHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	keywords, err := h.service.List(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keywords)
}

func (h *TrackedKeywordHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	keyword, err := h.service.Get(r.Context(), userID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keyword)
}

func (h *TrackedKeywordHandler) GetSearchResults(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	results, err := h.service.GetSearchResults(r.Context(), userID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return empty array if no results
	if results == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
