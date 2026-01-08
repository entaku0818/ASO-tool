package handler

import (
	"encoding/json"
	"net/http"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type ScreenshotHandler struct {
	service *service.ScreenshotService
}

func NewScreenshotHandler(service *service.ScreenshotService) *ScreenshotHandler {
	return &ScreenshotHandler{service: service}
}

func (h *ScreenshotHandler) Create(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")

	var req model.CreateScreenshotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.AppID = appID

	screenshot, err := h.service.Create(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, screenshot)
}

func (h *ScreenshotHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "screenshotID")

	screenshot, err := h.service.Get(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, screenshot)
}

func (h *ScreenshotHandler) List(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")
	deviceType := r.URL.Query().Get("device_type")

	var screenshots []*model.Screenshot
	var err error

	if deviceType != "" {
		screenshots, err = h.service.ListByAppAndDevice(r.Context(), appID, model.DeviceType(deviceType))
	} else {
		screenshots, err = h.service.ListByApp(r.Context(), appID)
	}

	if err != nil {
		handleServiceError(w, err)
		return
	}

	if screenshots == nil {
		screenshots = []*model.Screenshot{}
	}

	respondJSON(w, http.StatusOK, screenshots)
}

func (h *ScreenshotHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "screenshotID")

	var req model.UpdateScreenshotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	screenshot, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, screenshot)
}

func (h *ScreenshotHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "screenshotID")

	if err := h.service.Delete(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ScreenshotHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")

	var req model.ReorderScreenshotsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.Reorder(r.Context(), appID, &req); err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "reordered successfully"})
}
