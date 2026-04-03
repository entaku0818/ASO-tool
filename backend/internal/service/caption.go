package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const claudeAPIURL = "https://api.anthropic.com/v1/messages"
const claudeModel   = "claude-sonnet-4-6"

// CaptionService generates App Store screenshot captions using Claude API.
type CaptionService struct {
	apiKey     string
	httpClient *http.Client
}

func NewCaptionService() *CaptionService {
	return &CaptionService{
		apiKey: os.Getenv("ANTHROPIC_API_KEY"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *CaptionService) IsConfigured() bool {
	return s.apiKey != ""
}

// Generate calls Claude to produce exactly 3 caption candidates.
func (s *CaptionService) Generate(ctx context.Context, keywords []string, language string) ([]string, error) {
	if !s.IsConfigured() {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is not set")
	}

	prompt := buildCaptionPrompt(keywords, language)

	reqBody, err := json.Marshal(map[string]any{
		"model":      claudeModel,
		"max_tokens": 512,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, claudeAPIURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("claude api request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("claude api error %d: %s", resp.StatusCode, body)
	}

	// Parse Claude response: content[0].text contains the JSON array of captions
	var apiResp struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if len(apiResp.Content) == 0 {
		return nil, fmt.Errorf("empty response from claude")
	}

	captions, err := parseCaptionJSON(apiResp.Content[0].Text)
	if err != nil {
		return nil, fmt.Errorf("parse captions: %w", err)
	}
	return captions, nil
}

// buildCaptionPrompt constructs a tightly scoped prompt for caption generation.
func buildCaptionPrompt(keywords []string, language string) string {
	langName := map[string]string{
		"ja":      "Japanese",
		"en":      "English",
		"zh-Hans": "Simplified Chinese",
		"ko":      "Korean",
		"fr":      "French",
		"de":      "German",
	}
	lang := langName[language]
	if lang == "" {
		lang = language
	}

	kw := strings.Join(keywords, ", ")
	return fmt.Sprintf(`You are an App Store copywriter. Generate exactly 3 short, punchy screenshot captions in %s for an app related to these keywords: %s.

Rules:
- Each caption must be 1 sentence, under 30 characters if %s, or under 40 characters otherwise
- Captions should be compelling, action-oriented, and highlight user benefit
- Return ONLY a JSON array of 3 strings, no explanation

Example output format:
["Caption one", "Caption two", "Caption three"]`, lang, kw, lang)
}

// parseCaptionJSON extracts the string array from Claude's response text.
func parseCaptionJSON(text string) ([]string, error) {
	// Find the JSON array in the response (Claude may add surrounding text)
	start := strings.Index(text, "[")
	end := strings.LastIndex(text, "]")
	if start == -1 || end == -1 || end <= start {
		return nil, fmt.Errorf("no JSON array found in response")
	}

	var captions []string
	if err := json.Unmarshal([]byte(text[start:end+1]), &captions); err != nil {
		return nil, fmt.Errorf("unmarshal captions array: %w", err)
	}
	if len(captions) == 0 {
		return nil, fmt.Errorf("empty captions array")
	}
	return captions, nil
}
