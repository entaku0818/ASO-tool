package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/entaku0818/aso-compass/backend/internal/model"
)

func newTestSvc(handler http.HandlerFunc) (*TranslationService, func()) {
	srv := httptest.NewServer(handler)
	svc := &TranslationService{baseURL: srv.URL}
	return svc, srv.Close
}

func TestTranslate_EmptyText(t *testing.T) {
	svc := &TranslationService{baseURL: defaultMyMemoryURL}
	_, err := svc.Translate(context.Background(), TranslateRequest{Text: "", TargetLang: "en"})
	if err != model.ErrTranslationTextRequired {
		t.Errorf("got %v, want ErrTranslationTextRequired", err)
	}
}

func TestTranslate_EmptyTargetLang(t *testing.T) {
	svc := &TranslationService{baseURL: defaultMyMemoryURL}
	_, err := svc.Translate(context.Background(), TranslateRequest{Text: "写真編集", TargetLang: ""})
	if err != model.ErrTranslationTargetRequired {
		t.Errorf("got %v, want ErrTranslationTargetRequired", err)
	}
}

func TestTranslate_Success(t *testing.T) {
	svc, close := newTestSvc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"responseData":   map[string]interface{}{"translatedText": "photo editing"},
			"responseStatus": 200,
		})
	})
	defer close()

	resp, err := svc.Translate(context.Background(), TranslateRequest{Text: "写真編集", TargetLang: "en"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.TranslatedText != "photo editing" {
		t.Errorf("got %q, want 'photo editing'", resp.TranslatedText)
	}
	if resp.TargetLang != "en" {
		t.Errorf("got target_lang %q, want 'en'", resp.TargetLang)
	}
}

func TestTranslate_APIError(t *testing.T) {
	svc, close := newTestSvc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"responseData":   map[string]interface{}{"translatedText": ""},
			"responseStatus": 403,
		})
	})
	defer close()

	_, err := svc.Translate(context.Background(), TranslateRequest{Text: "写真編集", TargetLang: "en"})
	if err != model.ErrTranslationFailed {
		t.Errorf("got %v, want ErrTranslationFailed", err)
	}
}

func TestTranslate_EmptyResult(t *testing.T) {
	svc, close := newTestSvc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"responseData":   map[string]interface{}{"translatedText": ""},
			"responseStatus": 200,
		})
	})
	defer close()

	_, err := svc.Translate(context.Background(), TranslateRequest{Text: "写真編集", TargetLang: "en"})
	if err != model.ErrTranslationFailed {
		t.Errorf("got %v, want ErrTranslationFailed", err)
	}
}
