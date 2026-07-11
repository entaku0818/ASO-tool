package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/entaku0818/aso-tool/backend/internal/middleware"
	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type KeywordHandler struct {
	service           *service.KeywordService
	rankingService    *service.RankingService
	competitorService *service.CompetitorService
	userRepo          *repository.UserRepository
}

func NewKeywordHandler(service *service.KeywordService, rankingService *service.RankingService, competitorService *service.CompetitorService, userRepo *repository.UserRepository) *KeywordHandler {
	return &KeywordHandler{
		service:           service,
		rankingService:    rankingService,
		competitorService: competitorService,
		userRepo:          userRepo,
	}
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

// ExportCSV handles GET /api/apps/{appID}/keywords/export.csv (Pro plan only).
func (h *KeywordHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
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
		respondError(w, http.StatusForbidden, "CSV export requires Pro plan")
		return
	}

	appID := chi.URLParam(r, "appID")

	keywords, err := h.service.ListByApp(r.Context(), appID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	gaps, err := h.competitorService.GetKeywordGap(r.Context(), appID)
	if err != nil {
		gaps = nil // proceed without gap data
	}
	gapByKeyword := make(map[string]*model.KeywordGap, len(gaps))
	for _, g := range gaps {
		gapByKeyword[g.KeywordID] = g
	}

	var sb strings.Builder
	sb.WriteString("\xEF\xBB\xBF") // UTF-8 BOM so Excel detects the encoding correctly
	sb.WriteString("キーワード,国,自社順位,人気スコア,競合名,競合順位,最終更新\n")
	for _, k := range keywords {
		rank := "圏外"
		if latest, err := h.rankingService.GetLatestByKeyword(r.Context(), k.ID); err == nil && latest.Rank != nil {
			rank = strconv.Itoa(*latest.Rank)
		}

		popularity := ""
		if k.PopularityScore != nil {
			popularity = strconv.Itoa(*k.PopularityScore)
		}

		competitorName, competitorRank := "", ""
		if gap, ok := gapByKeyword[k.ID]; ok {
			competitorName = gap.CompetitorName
			competitorRank = strconv.Itoa(gap.CompetitorRank)
		}

		fetchedAt := ""
		if k.PopularityFetchedAt != nil {
			fetchedAt = k.PopularityFetchedAt.Format("2006/1/2")
		}

		row := []string{
			csvQuote(k.Keyword),
			k.Country,
			rank,
			popularity,
			csvQuote(competitorName),
			competitorRank,
			fetchedAt,
		}
		sb.WriteString(strings.Join(row, ","))
		sb.WriteString("\n")
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=keywords_%s.csv", appID))
	_, _ = w.Write([]byte(sb.String()))
}

func csvQuote(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}
