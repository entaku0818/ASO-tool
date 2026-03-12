package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/entaku0818/aso-tool/backend/internal/model"
)

const defaultMyMemoryURL = "https://api.mymemory.translated.net/get"

type TranslationService struct {
	baseURL string
}

func NewTranslationService() *TranslationService {
	return &TranslationService{baseURL: defaultMyMemoryURL}
}

type TranslateRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang,omitempty"`
	TargetLang string `json:"target_lang"`
}

type TranslateResponse struct {
	TranslatedText string `json:"translated_text"`
	SourceLang     string `json:"source_lang"`
	TargetLang     string `json:"target_lang"`
}

type myMemoryResponse struct {
	ResponseData struct {
		TranslatedText string `json:"translatedText"`
	} `json:"responseData"`
	ResponseStatus  int    `json:"responseStatus"`
	ResponseDetails string `json:"responseDetails"`
}

func (s *TranslationService) Translate(ctx context.Context, req TranslateRequest) (*TranslateResponse, error) {
	if req.Text == "" {
		return nil, model.ErrTranslationTextRequired
	}
	if req.TargetLang == "" {
		return nil, model.ErrTranslationTargetRequired
	}

	sourceLang := req.SourceLang
	if sourceLang == "" {
		sourceLang = "autodetect"
	}
	langPair := fmt.Sprintf("%s|%s", sourceLang, req.TargetLang)

	apiURL := s.baseURL + "?q=" + url.QueryEscape(req.Text) + "&langpair=" + url.QueryEscape(langPair)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("mymemory request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mymemory API error: status %d", resp.StatusCode)
	}

	var result myMemoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if result.ResponseStatus != 200 {
		return nil, model.ErrTranslationFailed
	}

	text := result.ResponseData.TranslatedText
	if text == "" {
		return nil, model.ErrTranslationFailed
	}

	return &TranslateResponse{
		TranslatedText: text,
		SourceLang:     sourceLang,
		TargetLang:     req.TargetLang,
	}, nil
}
