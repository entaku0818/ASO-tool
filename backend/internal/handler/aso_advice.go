package handler

import (
	"net/http"

	"github.com/entaku0818/aso-tool/backend/internal/middleware"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type ASOAdviceHandler struct {
	service *service.ASOAdviceService
}

func NewASOAdviceHandler(s *service.ASOAdviceService) *ASOAdviceHandler {
	return &ASOAdviceHandler{service: s}
}

// Generate handles GET /api/apps/{appID}/aso-advice
func (h *ASOAdviceHandler) Generate(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	appID := chi.URLParam(r, "appID")

	if !h.service.IsConfigured() {
		respondError(w, http.StatusServiceUnavailable, "AI advice not configured")
		return
	}

	advice, err := h.service.Generate(r.Context(), userID, appID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, advice)
}
