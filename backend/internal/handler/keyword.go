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

type KeywordHandler struct {
	service *service.KeywordService
}

func NewKeywordHandler(service *service.KeywordService) *KeywordHandler {
	return &KeywordHandler{service: service}
}

func (h *KeywordHandler) Create(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")
	userID := middleware.GetUserID(r.Context())

	var req model.CreateKeywordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.AppID = appID

	keyword, err := h.service.Create(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, model.ErrPlanLimitExceeded) {
			respondError(w, http.StatusPaymentRequired, err.Error())
			return
		}
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, keyword)
}

func (h *KeywordHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "keywordID")

	keyword, err := h.service.Get(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, keyword)
}

func (h *KeywordHandler) ListByApp(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")

	keywords, err := h.service.ListByApp(r.Context(), appID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if keywords == nil {
		keywords = []*model.Keyword{}
	}

	respondJSON(w, http.StatusOK, keywords)
}

func (h *KeywordHandler) Import(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")

	var req model.ImportKeywordsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.Keywords) == 0 {
		respondError(w, http.StatusBadRequest, "keywords is empty")
		return
	}

	resp, err := h.service.BulkImport(r.Context(), appID, &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h *KeywordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "keywordID")

	if err := h.service.Delete(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
