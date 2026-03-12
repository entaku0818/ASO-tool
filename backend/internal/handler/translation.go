package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/service"
)

type TranslationHandler struct {
	service *service.TranslationService
}

func NewTranslationHandler(svc *service.TranslationService) *TranslationHandler {
	return &TranslationHandler{service: svc}
}

// Translate handles POST /api/keywords/translate
func (h *TranslationHandler) Translate(w http.ResponseWriter, r *http.Request) {
	var req service.TranslateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.service.Translate(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrTranslationTextRequired),
			errors.Is(err, model.ErrTranslationTargetRequired):
			respondError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, model.ErrTranslationFailed):
			respondError(w, http.StatusUnprocessableEntity, "translation failed")
		default:
			respondError(w, http.StatusInternalServerError, "translation failed")
		}
		return
	}

	respondJSON(w, http.StatusOK, result)
}
