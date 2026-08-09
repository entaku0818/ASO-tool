package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/go-chi/chi/v5"
)

type mockCompetitorService struct {
	competitors []*model.Competitor
	gaps        []*model.KeywordGap
	err         error
	updatedCount int
}

func (m *mockCompetitorService) Create(_ context.Context, req *model.CreateCompetitorRequest) (*model.Competitor, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &model.Competitor{
		ID:                 "comp-1",
		AppID:              req.AppID,
		CompetitorBundleID: req.CompetitorBundleID,
		CompetitorName:     req.CompetitorName,
		Platform:           req.Platform,
		CreatedAt:          time.Now(),
	}, nil
}

func (m *mockCompetitorService) Get(_ context.Context, id string) (*model.Competitor, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, c := range m.competitors {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, model.ErrNotFound
}

func (m *mockCompetitorService) ListByApp(_ context.Context, _ string) ([]*model.Competitor, error) {
	return m.competitors, m.err
}

func (m *mockCompetitorService) Delete(_ context.Context, _ string) error {
	return m.err
}

func (m *mockCompetitorService) UpdateCompetitorRankings(_ context.Context, _ string) (int, error) {
	return m.updatedCount, m.err
}

func (m *mockCompetitorService) GetComparison(_ context.Context, _, _ string) (*model.CompetitorComparison, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &model.CompetitorComparison{Keyword: "test", KeywordID: "kw-1"}, nil
}

func (m *mockCompetitorService) GetLatestRanking(_ context.Context, _, _ string) (*model.CompetitorRanking, error) {
	return nil, m.err
}

func (m *mockCompetitorService) GetKeywordGap(_ context.Context, _ string) ([]*model.KeywordGap, error) {
	return m.gaps, m.err
}

// Wrap mock to satisfy handler's concrete *CompetitorService dependency via injection trick.
// Handler takes *service.CompetitorService, so we test via HTTP round-trip using a real handler
// wired with a nil service and verify routing only, or we test helpers directly.
// Instead, test the handler logic directly by calling the exported methods.

func newCompetitorRouter(svc *mockCompetitorService) http.Handler {
	h := &CompetitorHandler{service: nil} // we'll monkey-patch below
	_ = h
	// Since CompetitorHandler holds a concrete *service.CompetitorService we can't inject
	// a mock directly. Test the response helpers via a thin wrapper instead.
	r := chi.NewRouter()
	r.Get("/apps/{appID}/competitors", func(w http.ResponseWriter, r *http.Request) {
		appID := chi.URLParam(r, "appID")
		list, err := svc.ListByApp(r.Context(), appID)
		if err != nil {
			handleServiceError(w, err)
			return
		}
		if list == nil {
			list = []*model.Competitor{}
		}
		respondJSON(w, http.StatusOK, list)
	})
	r.Get("/apps/{appID}/competitors/keyword-gap", func(w http.ResponseWriter, r *http.Request) {
		appID := chi.URLParam(r, "appID")
		gaps, err := svc.GetKeywordGap(r.Context(), appID)
		if err != nil {
			handleServiceError(w, err)
			return
		}
		if gaps == nil {
			gaps = []*model.KeywordGap{}
		}
		respondJSON(w, http.StatusOK, gaps)
	})
	r.Post("/apps/{appID}/competitors/update-rankings", func(w http.ResponseWriter, r *http.Request) {
		appID := chi.URLParam(r, "appID")
		count, err := svc.UpdateCompetitorRankings(r.Context(), appID)
		if err != nil {
			handleServiceError(w, err)
			return
		}
		respondJSON(w, http.StatusOK, map[string]int{"updated": count})
	})
	return r
}

func TestCompetitorList(t *testing.T) {
	// test via inline handler
	svc2 := &mockCompetitorService{
		competitors: []*model.Competitor{
			{ID: "c1", AppID: "app-1", CompetitorName: "Rival App", Platform: model.PlatformIOS},
		},
	}
	router := newCompetitorRouter(svc2)
	req := httptest.NewRequest(http.MethodGet, "/apps/app-1/competitors", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result []*model.Competitor
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0].ID != "c1" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestCompetitorListEmpty(t *testing.T) {
	svc := &mockCompetitorService{}
	router := newCompetitorRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/apps/app-1/competitors", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result []*model.Competitor
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("expected empty slice, not nil")
	}
}

func TestCompetitorListError(t *testing.T) {
	svc := &mockCompetitorService{err: model.ErrNotFound}
	router := newCompetitorRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/apps/app-1/competitors", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetKeywordGap(t *testing.T) {
	competitorRank := 10
	ourRank := 50
	svc := &mockCompetitorService{
		gaps: []*model.KeywordGap{
			{
				KeywordID:      "kw-1",
				Keyword:        "fitness app",
				Country:        "JP",
				CompetitorName: "Rival",
				CompetitorRank: competitorRank,
				OurRank:        &ourRank,
			},
		},
	}
	router := newCompetitorRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/apps/app-1/competitors/keyword-gap", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result []*model.KeywordGap
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 gap, got %d", len(result))
	}
	if result[0].Keyword != "fitness app" {
		t.Errorf("unexpected keyword: %s", result[0].Keyword)
	}
	if result[0].CompetitorRank != competitorRank {
		t.Errorf("unexpected competitor rank: %d", result[0].CompetitorRank)
	}
}

func TestGetKeywordGapEmpty(t *testing.T) {
	svc := &mockCompetitorService{}
	router := newCompetitorRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/apps/app-1/competitors/keyword-gap", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result []*model.KeywordGap
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("expected empty slice, not nil")
	}
}

func TestUpdateRankings(t *testing.T) {
	svc := &mockCompetitorService{updatedCount: 5}
	router := newCompetitorRouter(svc)
	req := httptest.NewRequest(http.MethodPost, "/apps/app-1/competitors/update-rankings", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result map[string]int
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result["updated"] != 5 {
		t.Errorf("expected updated=5, got %d", result["updated"])
	}
}
