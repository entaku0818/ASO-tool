package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/go-chi/chi/v5"
)

type mockMetadataService struct {
	versions []*model.AppMetadataVersion
	err      error
}

func (m *mockMetadataService) List(_ context.Context, _ string) ([]*model.AppMetadataVersion, error) {
	return m.versions, m.err
}

func (m *mockMetadataService) Get(_ context.Context, id string) (*model.AppMetadataVersion, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, v := range m.versions {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, model.ErrNotFound
}

func (m *mockMetadataService) Upsert(_ context.Context, req *model.UpsertMetadataRequest) (*model.AppMetadataVersion, error) {
	if m.err != nil {
		return nil, m.err
	}
	title := req.Title
	return &model.AppMetadataVersion{
		ID:         "meta-1",
		AppID:      req.AppID,
		Locale:     req.Locale,
		VersionTag: req.VersionTag,
		Title:      title,
		UpdatedAt:  time.Now(),
		CreatedAt:  time.Now(),
	}, nil
}

func (m *mockMetadataService) Delete(_ context.Context, _ string) error {
	return m.err
}

func newMetadataRouter(svc *mockMetadataService) http.Handler {
	r := chi.NewRouter()
	r.Get("/apps/{appID}/metadata", func(w http.ResponseWriter, r *http.Request) {
		appID := chi.URLParam(r, "appID")
		items, err := svc.List(r.Context(), appID)
		if err != nil {
			handleServiceError(w, err)
			return
		}
		if items == nil {
			items = []*model.AppMetadataVersion{}
		}
		respondJSON(w, http.StatusOK, items)
	})
	r.Put("/apps/{appID}/metadata", func(w http.ResponseWriter, r *http.Request) {
		appID := chi.URLParam(r, "appID")
		var req model.UpsertMetadataRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		req.AppID = appID
		item, err := svc.Upsert(r.Context(), &req)
		if err != nil {
			handleServiceError(w, err)
			return
		}
		respondJSON(w, http.StatusOK, item)
	})
	r.Delete("/apps/{appID}/metadata/{metadataID}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "metadataID")
		if err := svc.Delete(r.Context(), id); err != nil {
			handleServiceError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	return r
}

func TestMetadataList(t *testing.T) {
	title := "My App"
	svc := &mockMetadataService{
		versions: []*model.AppMetadataVersion{
			{ID: "m1", AppID: "app-1", Locale: "ja", VersionTag: "draft", Title: &title},
		},
	}
	r := newMetadataRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/apps/app-1/metadata", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result []*model.AppMetadataVersion
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}
	if *result[0].Title != "My App" {
		t.Errorf("unexpected title: %v", result[0].Title)
	}
}

func TestMetadataListEmpty(t *testing.T) {
	svc := &mockMetadataService{}
	r := newMetadataRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/apps/app-1/metadata", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result []*model.AppMetadataVersion
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("expected empty slice, not nil")
	}
}

func TestMetadataUpsert(t *testing.T) {
	svc := &mockMetadataService{}
	r := newMetadataRouter(svc)

	title := "テストアプリ"
	body := map[string]interface{}{
		"locale":      "ja",
		"version_tag": "draft",
		"title":       title,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/apps/app-1/metadata", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var result model.AppMetadataVersion
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.Locale != "ja" {
		t.Errorf("unexpected locale: %s", result.Locale)
	}
	if result.Title == nil || *result.Title != title {
		t.Errorf("unexpected title: %v", result.Title)
	}
}

func TestMetadataUpsertInvalidBody(t *testing.T) {
	svc := &mockMetadataService{}
	r := newMetadataRouter(svc)
	req := httptest.NewRequest(http.MethodPut, "/apps/app-1/metadata", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestMetadataDelete(t *testing.T) {
	svc := &mockMetadataService{}
	r := newMetadataRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/apps/app-1/metadata/meta-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestMetadataDeleteNotFound(t *testing.T) {
	svc := &mockMetadataService{err: model.ErrNotFound}
	r := newMetadataRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/apps/app-1/metadata/missing", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
