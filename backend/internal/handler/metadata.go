package handler

import (
	"encoding/json"
	"net/http"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type MetadataHandler struct {
	service *service.MetadataService
}

func NewMetadataHandler(service *service.MetadataService) *MetadataHandler {
	return &MetadataHandler{service: service}
}

func (h *MetadataHandler) List(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")
	items, err := h.service.List(r.Context(), appID)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	if items == nil {
		items = []*model.AppMetadataVersion{}
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *MetadataHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")
	var req model.UpsertMetadataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.AppID = appID

	item, err := h.service.Upsert(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, item)
}

func (h *MetadataHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "metadataID")
	if err := h.service.Delete(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
