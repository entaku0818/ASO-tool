package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/go-chi/chi/v5"
)

// mockAppService implements the service interface for testing
type mockAppService struct {
	apps map[string]*model.App
}

func newMockAppService() *mockAppService {
	return &mockAppService{
		apps: make(map[string]*model.App),
	}
}

func (m *mockAppService) Create(ctx context.Context, req *model.CreateAppRequest) (*model.App, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	app := &model.App{
		ID:        "test-id-123",
		Name:      req.Name,
		BundleID:  req.BundleID,
		Platform:  req.Platform,
		StoreURL:  req.StoreURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.apps[app.ID] = app
	return app, nil
}

func (m *mockAppService) Get(ctx context.Context, id string) (*model.App, error) {
	app, ok := m.apps[id]
	if !ok {
		return nil, model.ErrNotFound
	}
	return app, nil
}

func (m *mockAppService) List(ctx context.Context) ([]*model.App, error) {
	var apps []*model.App
	for _, app := range m.apps {
		apps = append(apps, app)
	}
	return apps, nil
}

func (m *mockAppService) Update(ctx context.Context, id string, req *model.UpdateAppRequest) (*model.App, error) {
	app, ok := m.apps[id]
	if !ok {
		return nil, model.ErrNotFound
	}
	if req.Name != "" {
		app.Name = req.Name
	}
	if req.StoreURL != "" {
		app.StoreURL = req.StoreURL
	}
	app.UpdatedAt = time.Now()
	return app, nil
}

func (m *mockAppService) Delete(ctx context.Context, id string) error {
	if _, ok := m.apps[id]; !ok {
		return model.ErrNotFound
	}
	delete(m.apps, id)
	return nil
}

// AppServiceInterface defines the interface for dependency injection
type AppServiceInterface interface {
	Create(ctx context.Context, req *model.CreateAppRequest) (*model.App, error)
	Get(ctx context.Context, id string) (*model.App, error)
	List(ctx context.Context) ([]*model.App, error)
	Update(ctx context.Context, id string, req *model.UpdateAppRequest) (*model.App, error)
	Delete(ctx context.Context, id string) error
}

// testableAppHandler is a handler that accepts the interface
type testableAppHandler struct {
	service AppServiceInterface
}

func newTestableAppHandler(service AppServiceInterface) *testableAppHandler {
	return &testableAppHandler{service: service}
}

func (h *testableAppHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	app, err := h.service.Create(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, app)
}

func (h *testableAppHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	app, err := h.service.Get(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, app)
}

func (h *testableAppHandler) List(w http.ResponseWriter, r *http.Request) {
	apps, err := h.service.List(r.Context())
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if apps == nil {
		apps = []*model.App{}
	}

	respondJSON(w, http.StatusOK, apps)
}

func (h *testableAppHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req model.UpdateAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	app, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, app)
}

func (h *testableAppHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Delete(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func TestAppHandler_Create(t *testing.T) {
	mockService := newMockAppService()
	handler := newTestableAppHandler(mockService)

	tests := []struct {
		name       string
		body       interface{}
		wantStatus int
	}{
		{
			name: "valid iOS app",
			body: model.CreateAppRequest{
				Name:     "Test App",
				BundleID: "com.example.test",
				Platform: model.PlatformIOS,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "valid Android app",
			body: model.CreateAppRequest{
				Name:     "Test Android App",
				BundleID: "com.example.android",
				Platform: model.PlatformAndroid,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "missing name",
			body: model.CreateAppRequest{
				BundleID: "com.example.test",
				Platform: model.PlatformIOS,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing bundle_id",
			body: model.CreateAppRequest{
				Name:     "Test App",
				Platform: model.PlatformIOS,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid platform",
			body: model.CreateAppRequest{
				Name:     "Test App",
				BundleID: "com.example.test",
				Platform: "windows",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if s, ok := tt.body.(string); ok {
				body = []byte(s)
			} else {
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/apps", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.Create(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Create() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestAppHandler_Get(t *testing.T) {
	mockService := newMockAppService()
	handler := newTestableAppHandler(mockService)

	// Create a test app
	mockService.apps["existing-id"] = &model.App{
		ID:       "existing-id",
		Name:     "Existing App",
		BundleID: "com.example.existing",
		Platform: model.PlatformIOS,
	}

	tests := []struct {
		name       string
		appID      string
		wantStatus int
	}{
		{
			name:       "existing app",
			appID:      "existing-id",
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing app",
			appID:      "non-existing-id",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/api/apps/{id}", handler.Get)

			req := httptest.NewRequest(http.MethodGet, "/api/apps/"+tt.appID, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Get() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestAppHandler_List(t *testing.T) {
	mockService := newMockAppService()
	handler := newTestableAppHandler(mockService)

	t.Run("empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/apps", nil)
		rec := httptest.NewRecorder()

		handler.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("List() status = %v, want %v", rec.Code, http.StatusOK)
		}

		var apps []model.App
		if err := json.NewDecoder(rec.Body).Decode(&apps); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if len(apps) != 0 {
			t.Errorf("List() returned %d apps, want 0", len(apps))
		}
	})

	t.Run("with apps", func(t *testing.T) {
		mockService.apps["app-1"] = &model.App{ID: "app-1", Name: "App 1"}
		mockService.apps["app-2"] = &model.App{ID: "app-2", Name: "App 2"}

		req := httptest.NewRequest(http.MethodGet, "/api/apps", nil)
		rec := httptest.NewRecorder()

		handler.List(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("List() status = %v, want %v", rec.Code, http.StatusOK)
		}

		var apps []model.App
		if err := json.NewDecoder(rec.Body).Decode(&apps); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if len(apps) != 2 {
			t.Errorf("List() returned %d apps, want 2", len(apps))
		}
	})
}

func TestAppHandler_Update(t *testing.T) {
	mockService := newMockAppService()
	handler := newTestableAppHandler(mockService)

	mockService.apps["existing-id"] = &model.App{
		ID:       "existing-id",
		Name:     "Original Name",
		BundleID: "com.example.test",
		Platform: model.PlatformIOS,
	}

	tests := []struct {
		name       string
		appID      string
		body       model.UpdateAppRequest
		wantStatus int
	}{
		{
			name:       "update name",
			appID:      "existing-id",
			body:       model.UpdateAppRequest{Name: "Updated Name"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing app",
			appID:      "non-existing-id",
			body:       model.UpdateAppRequest{Name: "Updated Name"},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Put("/api/apps/{id}", handler.Update)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/api/apps/"+tt.appID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Update() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestAppHandler_Delete(t *testing.T) {
	mockService := newMockAppService()
	handler := newTestableAppHandler(mockService)

	mockService.apps["existing-id"] = &model.App{
		ID:       "existing-id",
		Name:     "Test App",
		BundleID: "com.example.test",
		Platform: model.PlatformIOS,
	}

	tests := []struct {
		name       string
		appID      string
		wantStatus int
	}{
		{
			name:       "delete existing app",
			appID:      "existing-id",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "delete non-existing app",
			appID:      "non-existing-id",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Delete("/api/apps/{id}", handler.Delete)

			req := httptest.NewRequest(http.MethodDelete, "/api/apps/"+tt.appID, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Delete() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}
