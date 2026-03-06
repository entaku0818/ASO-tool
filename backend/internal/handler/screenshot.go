package handler

import (
	"encoding/base64"
	"encoding/json"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/entaku0818/aso-tool/backend/internal/imgproc"
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

// Generate renders a framed screenshot image for each provided caption language
// and returns the results as base64-encoded PNG data URLs.
func (h *ScreenshotHandler) Generate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "failed to parse form: "+err.Error())
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		respondError(w, http.StatusBadRequest, "image field is required")
		return
	}
	defer func() { _ = file.Close() }()

	img, _, err := image.Decode(file)
	if err != nil {
		respondError(w, http.StatusBadRequest, "failed to decode image: "+err.Error())
		return
	}

	device := r.FormValue("device")
	if device == "" {
		device = "iphone67"
	}
	bgColor := r.FormValue("bg_color")
	if bgColor == "" {
		bgColor = "#4F46E5"
	}
	textColor := r.FormValue("text_color")
	if textColor == "" {
		textColor = "#FFFFFF"
	}

	var captions map[string]string
	if raw := r.FormValue("captions"); raw != "" {
		if err := json.Unmarshal([]byte(raw), &captions); err != nil {
			respondError(w, http.StatusBadRequest, "invalid captions JSON: "+err.Error())
			return
		}
	}

	result, err := imgproc.Generate(&imgproc.GenerateRequest{
		Image:     img,
		Device:    device,
		BGColor:   bgColor,
		TextColor: textColor,
		Captions:  captions,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "image generation failed: "+err.Error())
		return
	}

	images := make(map[string]string, len(result.Images))
	for lang, pngBytes := range result.Images {
		images[lang] = "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"images": images})
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
