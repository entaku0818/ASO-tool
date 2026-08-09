package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/entaku0818/aso-compass/backend/internal/service"
)

type translationServiceInterface interface {
	Translate(ctx context.Context, req service.TranslateRequest) (*service.TranslateResponse, error)
}

type mockTranslationService struct {
	response *service.TranslateResponse
	err      error
}

func (m *mockTranslationService) Translate(_ context.Context, _ service.TranslateRequest) (*service.TranslateResponse, error) {
	return m.response, m.err
}

type translationHandlerWithMock struct {
	svc translationServiceInterface
}

func (h *translationHandlerWithMock) Translate(w http.ResponseWriter, r *http.Request) {
	var req service.TranslateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := h.svc.Translate(r.Context(), req)
	if err != nil {
		switch err {
		case model.ErrTranslationTextRequired, model.ErrTranslationTargetRequired:
			respondError(w, http.StatusBadRequest, err.Error())
		case model.ErrTranslationFailed:
			respondError(w, http.StatusUnprocessableEntity, "translation failed")
		default:
			respondError(w, http.StatusInternalServerError, "translation failed")
		}
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func TestTranslationHandler_Success(t *testing.T) {
	mock := &mockTranslationService{
		response: &service.TranslateResponse{TranslatedText: "photo editing", SourceLang: "ja", TargetLang: "en"},
	}
	h := &translationHandlerWithMock{svc: mock}

	body, _ := json.Marshal(service.TranslateRequest{Text: "写真編集", TargetLang: "en"})
	req := httptest.NewRequest(http.MethodPost, "/api/keywords/translate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Translate(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp service.TranslateResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.TranslatedText != "photo editing" {
		t.Errorf("got %q, want 'photo editing'", resp.TranslatedText)
	}
}

func TestTranslationHandler_TextRequired(t *testing.T) {
	mock := &mockTranslationService{err: model.ErrTranslationTextRequired}
	h := &translationHandlerWithMock{svc: mock}

	body, _ := json.Marshal(service.TranslateRequest{Text: "", TargetLang: "en"})
	req := httptest.NewRequest(http.MethodPost, "/api/keywords/translate", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Translate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestTranslationHandler_TranslationFailed(t *testing.T) {
	mock := &mockTranslationService{err: model.ErrTranslationFailed}
	h := &translationHandlerWithMock{svc: mock}

	body, _ := json.Marshal(service.TranslateRequest{Text: "写真編集", TargetLang: "en"})
	req := httptest.NewRequest(http.MethodPost, "/api/keywords/translate", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.Translate(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestTranslationHandler_InvalidBody(t *testing.T) {
	mock := &mockTranslationService{}
	h := &translationHandlerWithMock{svc: mock}

	req := httptest.NewRequest(http.MethodPost, "/api/keywords/translate", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()
	h.Translate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
