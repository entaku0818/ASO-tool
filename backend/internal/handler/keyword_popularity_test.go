package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/entaku0818/aso-compass/backend/internal/scraper"
	"github.com/go-chi/chi/v5"
)

// keywordPopularityServiceInterface defines the methods the handler requires.
type keywordPopularityServiceInterface interface {
	SetCredentials(ctx context.Context, userID, appID string, req *model.CreateSearchAdsCredentialsRequest) error
	GetCredentials(ctx context.Context, userID, appID string) (*model.SearchAdsCredentials, error)
	DeleteCredentials(ctx context.Context, userID, appID string) error
	RefreshScores(ctx context.Context, userID, appID string) (int, error)
	GetSuggestions(ctx context.Context, userID, appID string, limit int) ([]scraper.KeywordPopularity, error)
	GetCompetitorSuggestions(ctx context.Context, userID, appID string, competitorAdamID int64, limit int) ([]scraper.KeywordPopularity, error)
}

// mockKeywordPopularityService implements keywordPopularityServiceInterface.
type mockKeywordPopularityService struct {
	creds       *model.SearchAdsCredentials
	suggestions []scraper.KeywordPopularity
	err         error
	updated     int
}

func (m *mockKeywordPopularityService) SetCredentials(_ context.Context, _, _ string, _ *model.CreateSearchAdsCredentialsRequest) error {
	return m.err
}

func (m *mockKeywordPopularityService) GetCredentials(_ context.Context, _, _ string) (*model.SearchAdsCredentials, error) {
	return m.creds, m.err
}

func (m *mockKeywordPopularityService) DeleteCredentials(_ context.Context, _, _ string) error {
	return m.err
}

func (m *mockKeywordPopularityService) RefreshScores(_ context.Context, _, _ string) (int, error) {
	return m.updated, m.err
}

func (m *mockKeywordPopularityService) GetSuggestions(_ context.Context, _, _ string, _ int) ([]scraper.KeywordPopularity, error) {
	return m.suggestions, m.err
}

func (m *mockKeywordPopularityService) GetCompetitorSuggestions(_ context.Context, _, _ string, _ int64, _ int) ([]scraper.KeywordPopularity, error) {
	return m.suggestions, m.err
}

// testableKeywordPopularityHandler wraps the interface for testing.
type testableKeywordPopularityHandler struct {
	service keywordPopularityServiceInterface
}

func newTestableKeywordPopularityHandler(svc keywordPopularityServiceInterface) *testableKeywordPopularityHandler {
	return &testableKeywordPopularityHandler{service: svc}
}

func (h *testableKeywordPopularityHandler) SetCredentials(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")
	var req model.CreateSearchAdsCredentialsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.SetCredentials(r.Context(), "", appID, &req); err != nil {
		handleSearchAdsError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "credentials saved"})
}

func (h *testableKeywordPopularityHandler) GetCredentials(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")
	creds, err := h.service.GetCredentials(r.Context(), "", appID)
	if err != nil {
		handleSearchAdsError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, creds)
}

func (h *testableKeywordPopularityHandler) DeleteCredentials(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")
	if err := h.service.DeleteCredentials(r.Context(), "", appID); err != nil {
		handleSearchAdsError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *testableKeywordPopularityHandler) RefreshScores(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")
	updated, err := h.service.RefreshScores(r.Context(), "", appID)
	if err != nil {
		handleSearchAdsError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, map[string]int{"updated": updated})
}

func (h *testableKeywordPopularityHandler) GetSuggestions(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")
	suggestions, err := h.service.GetSuggestions(r.Context(), "", appID, 25)
	if err != nil {
		handleSearchAdsError(w, err)
		return
	}
	if suggestions == nil {
		suggestions = []scraper.KeywordPopularity{}
	}
	respondJSON(w, http.StatusOK, suggestions)
}

func (h *testableKeywordPopularityHandler) GetCompetitorSuggestions(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "appID")
	adamIDStr := r.URL.Query().Get("adam_id")
	if adamIDStr == "" {
		respondError(w, http.StatusBadRequest, "adam_id is required")
		return
	}
	suggestions, err := h.service.GetCompetitorSuggestions(r.Context(), "", appID, 0, 25)
	if err != nil {
		handleSearchAdsError(w, err)
		return
	}
	if suggestions == nil {
		suggestions = []scraper.KeywordPopularity{}
	}
	respondJSON(w, http.StatusOK, suggestions)
}

// Helpers for routing in tests
func newRouterWith(method, path string, handlerFn http.HandlerFunc) *chi.Mux {
	r := chi.NewRouter()
	switch method {
	case http.MethodGet:
		r.Get(path, handlerFn)
	case http.MethodPost:
		r.Post(path, handlerFn)
	case http.MethodDelete:
		r.Delete(path, handlerFn)
	}
	return r
}

func TestKeywordPopularityHandler_SetCredentials(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		serviceErr error
		wantStatus int
	}{
		{
			name: "valid credentials",
			body: model.CreateSearchAdsCredentialsRequest{
				ClientID:   "client-id",
				TeamID:     "team-id",
				KeyID:      "key-id",
				PrivateKey: "base64key",
				AdamID:     123456789,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "missing client_id",
			body: model.CreateSearchAdsCredentialsRequest{
				TeamID:     "team-id",
				KeyID:      "key-id",
				PrivateKey: "base64key",
				AdamID:     123456789,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing adam_id",
			body: model.CreateSearchAdsCredentialsRequest{
				ClientID:   "client-id",
				TeamID:     "team-id",
				KeyID:      "key-id",
				PrivateKey: "base64key",
				AdamID:     0,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid JSON",
			body:       "not-json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service returns not found",
			body: model.CreateSearchAdsCredentialsRequest{
				ClientID:   "client-id",
				TeamID:     "team-id",
				KeyID:      "key-id",
				PrivateKey: "base64key",
				AdamID:     123456789,
			},
			serviceErr: model.ErrNotFound,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockKeywordPopularityService{err: tt.serviceErr}
			h := newTestableKeywordPopularityHandler(mock)
			r := newRouterWith(http.MethodPost, "/api/apps/{appID}/search-ads-credentials", h.SetCredentials)

			var body []byte
			if s, ok := tt.body.(string); ok {
				body = []byte(s)
			} else {
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/apps/app-1/search-ads-credentials", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("SetCredentials() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestKeywordPopularityHandler_GetCredentials(t *testing.T) {
	tests := []struct {
		name       string
		creds      *model.SearchAdsCredentials
		serviceErr error
		wantStatus int
	}{
		{
			name: "returns credentials",
			creds: &model.SearchAdsCredentials{
				ID:       "cred-id",
				AppID:    "app-1",
				ClientID: "client-id",
				AdamID:   123456789,
				IsValid:  true,
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "credentials not found",
			serviceErr: model.ErrSearchAdsCredentialsNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "app not found",
			serviceErr: model.ErrNotFound,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockKeywordPopularityService{creds: tt.creds, err: tt.serviceErr}
			h := newTestableKeywordPopularityHandler(mock)
			r := newRouterWith(http.MethodGet, "/api/apps/{appID}/search-ads-credentials", h.GetCredentials)

			req := httptest.NewRequest(http.MethodGet, "/api/apps/app-1/search-ads-credentials", nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("GetCredentials() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestKeywordPopularityHandler_DeleteCredentials(t *testing.T) {
	tests := []struct {
		name       string
		serviceErr error
		wantStatus int
	}{
		{
			name:       "success",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "not found",
			serviceErr: model.ErrSearchAdsCredentialsNotFound,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockKeywordPopularityService{err: tt.serviceErr}
			h := newTestableKeywordPopularityHandler(mock)
			r := newRouterWith(http.MethodDelete, "/api/apps/{appID}/search-ads-credentials", h.DeleteCredentials)

			req := httptest.NewRequest(http.MethodDelete, "/api/apps/app-1/search-ads-credentials", nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("DeleteCredentials() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestKeywordPopularityHandler_RefreshScores(t *testing.T) {
	tests := []struct {
		name        string
		updated     int
		serviceErr  error
		wantStatus  int
		wantUpdated int
	}{
		{
			name:        "refreshes 3 keywords",
			updated:     3,
			wantStatus:  http.StatusOK,
			wantUpdated: 3,
		},
		{
			name:       "invalid credentials",
			serviceErr: model.ErrInvalidSearchAdsCredentials,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "app not found",
			serviceErr: model.ErrNotFound,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockKeywordPopularityService{updated: tt.updated, err: tt.serviceErr}
			h := newTestableKeywordPopularityHandler(mock)
			r := newRouterWith(http.MethodPost, "/api/apps/{appID}/keywords/refresh-popularity", h.RefreshScores)

			req := httptest.NewRequest(http.MethodPost, "/api/apps/app-1/keywords/refresh-popularity", nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("RefreshScores() status = %v, want %v", rec.Code, tt.wantStatus)
			}

			if rec.Code == http.StatusOK {
				var resp map[string]int
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if resp["updated"] != tt.wantUpdated {
					t.Errorf("updated = %v, want %v", resp["updated"], tt.wantUpdated)
				}
			}
		})
	}
}

func TestKeywordPopularityHandler_GetSuggestions(t *testing.T) {
	tests := []struct {
		name        string
		suggestions []scraper.KeywordPopularity
		serviceErr  error
		wantStatus  int
		wantCount   int
	}{
		{
			name: "returns suggestions",
			suggestions: []scraper.KeywordPopularity{
				{Text: "calculator", PopularityScore: 5},
				{Text: "math", PopularityScore: 3},
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:        "returns empty list when nil",
			suggestions: nil,
			wantStatus:  http.StatusOK,
			wantCount:   0,
		},
		{
			name:       "credentials not found",
			serviceErr: model.ErrSearchAdsCredentialsNotFound,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockKeywordPopularityService{suggestions: tt.suggestions, err: tt.serviceErr}
			h := newTestableKeywordPopularityHandler(mock)
			r := newRouterWith(http.MethodGet, "/api/apps/{appID}/keywords/suggestions", h.GetSuggestions)

			req := httptest.NewRequest(http.MethodGet, "/api/apps/app-1/keywords/suggestions", nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("GetSuggestions() status = %v, want %v", rec.Code, tt.wantStatus)
			}

			if rec.Code == http.StatusOK {
				var resp []scraper.KeywordPopularity
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if len(resp) != tt.wantCount {
					t.Errorf("len(suggestions) = %v, want %v", len(resp), tt.wantCount)
				}
			}
		})
	}
}

func TestKeywordPopularityHandler_GetCompetitorSuggestions(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{
			name:       "valid adam_id",
			query:      "?adam_id=123456789",
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing adam_id",
			query:      "",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockKeywordPopularityService{
				suggestions: []scraper.KeywordPopularity{
					{Text: "competitor-keyword", PopularityScore: 4},
				},
			}
			h := newTestableKeywordPopularityHandler(mock)
			r := newRouterWith(http.MethodGet, "/api/apps/{appID}/keywords/competitor-suggestions", h.GetCompetitorSuggestions)

			req := httptest.NewRequest(http.MethodGet, "/api/apps/app-1/keywords/competitor-suggestions"+tt.query, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("GetCompetitorSuggestions() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandleSearchAdsError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "ErrNotFound",
			err:        model.ErrNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "ErrSearchAdsCredentialsNotFound",
			err:        model.ErrSearchAdsCredentialsNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "ErrInvalidPrivateKey",
			err:        model.ErrInvalidPrivateKey,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "ErrClientIDRequired",
			err:        model.ErrClientIDRequired,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "ErrAdamIDRequired",
			err:        model.ErrAdamIDRequired,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "ErrInvalidSearchAdsCredentials",
			err:        model.ErrInvalidSearchAdsCredentials,
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handleSearchAdsError(rec, tt.err)
			if rec.Code != tt.wantStatus {
				t.Errorf("handleSearchAdsError(%v) status = %v, want %v", tt.err, rec.Code, tt.wantStatus)
			}
		})
	}
}
