package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/entaku0818/aso-tool/backend/internal/middleware"
	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type CaptionHandler struct {
	service  *service.CaptionService
	userRepo *repository.UserRepository
}

func NewCaptionHandler(svc *service.CaptionService, userRepo *repository.UserRepository) *CaptionHandler {
	return &CaptionHandler{service: svc, userRepo: userRepo}
}

// Generate handles POST /api/apps/{appID}/captions/generate
//
// Body: { "keywords": ["word1", "word2"], "language": "ja", "bulk": false }
//   - language: BCP-47 tag (e.g. "ja", "en") or "all" (Pro only)
//   - bulk:     true → generate for all 6 languages (Pro only)
//
// Free plan: single-language requests only. bulk=true or language="all" → 403.
func (h *CaptionHandler) Generate(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "appID")
	userID := middleware.GetUserID(r.Context())

	var req struct {
		Keywords []string `json:"keywords"`
		Language string   `json:"language"`
		Bulk     bool     `json:"bulk"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.Keywords) == 0 {
		respondError(w, http.StatusBadRequest, "keywords is required")
		return
	}
	req.Language = strings.TrimSpace(req.Language)
	if req.Language == "" {
		req.Language = "ja"
	}

	// Pro gate: bulk generation and "all" language mode require Pro plan
	isBulk := req.Bulk || req.Language == "all"
	if isBulk {
		user, err := h.userRepo.GetByID(r.Context(), userID)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				respondError(w, http.StatusUnauthorized, "user not found")
			} else {
				respondError(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}
		if !user.IsPro() {
			respondError(w, http.StatusForbidden, "bulk AI caption generation requires Pro plan")
			return
		}
	}

	if !h.service.IsConfigured() {
		respondError(w, http.StatusServiceUnavailable, "AI caption service is not configured")
		return
	}

	captions, err := h.service.Generate(r.Context(), req.Keywords, req.Language)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate captions")
		return
	}

	respondJSON(w, http.StatusOK, map[string][]string{"captions": captions})
}
