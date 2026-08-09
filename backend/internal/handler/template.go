package handler

import (
	"net/http"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/entaku0818/aso-compass/backend/internal/service"
)

type TemplateHandler struct {
	service *service.TemplateService
}

func NewTemplateHandler(svc *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{service: svc}
}

// List handles GET /api/templates?category=...
func (h *TemplateHandler) List(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	templates, err := h.service.List(r.Context(), category)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch templates")
		return
	}

	if templates == nil {
		templates = []*model.Template{}
	}

	respondJSON(w, http.StatusOK, templates)
}
