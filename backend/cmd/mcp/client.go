package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// DefaultBaseURL is the production API. Override with ASO_API_BASE.
const DefaultBaseURL = "https://aso-api-671942133800.asia-northeast1.run.app"

type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// apiError carries the status code so callers can distinguish auth failures
// from plan restrictions (403) and missing resources (404).
type apiError struct {
	Status int
	Body   string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.Status, e.Body)
}

// Client talks to the ASO-tool REST API on behalf of one license.
//
// The backend issues JWTs that expire after 24h, so rather than caching a token
// across runs the client activates on first use and re-activates once on a 401.
// Activation is idempotent for the same email (service/license.go), which makes
// this safe to do on every process start.
type Client struct {
	BaseURL string
	HTTP    httpDoer

	licenseKey string
	email      string

	mu    sync.Mutex
	token string
}

func NewClient(baseURL, licenseKey, email string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		HTTP:       &http.Client{Timeout: 60 * time.Second},
		licenseKey: licenseKey,
		email:      email,
	}
}

type activateResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Plan  string `json:"plan"`
	} `json:"user"`
}

// activate exchanges the license key for a JWT. The endpoint takes no auth
// header, so calling it from the 401 path cannot recurse.
func (c *Client) activate(ctx context.Context) (*activateResponse, error) {
	body, err := json.Marshal(map[string]string{
		"key":   c.licenseKey,
		"email": c.email,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/api/licenses/activate", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	data, err := c.send(req)
	if err != nil {
		return nil, fmt.Errorf("activate: %w", err)
	}
	var out activateResponse
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("activate: decode: %w", err)
	}
	if out.Token == "" {
		return nil, fmt.Errorf("activate: response contained no token")
	}

	c.mu.Lock()
	c.token = out.Token
	c.mu.Unlock()
	return &out, nil
}

func (c *Client) currentToken() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.token
}

// Get fetches path and decodes the JSON body into out. If the token is missing
// or rejected, it activates (once) and retries.
func (c *Client) Get(ctx context.Context, path string, out any) error {
	if c.currentToken() == "" {
		if _, err := c.activate(ctx); err != nil {
			return err
		}
	}

	data, err := c.getOnce(ctx, path)
	if apiErr, ok := err.(*apiError); ok && apiErr.Status == http.StatusUnauthorized {
		if _, aerr := c.activate(ctx); aerr != nil {
			return aerr
		}
		data, err = c.getOnce(ctx, path)
	}
	if err != nil {
		return err
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}

func (c *Client) getOnce(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.currentToken())
	return c.send(req)
}

func (c *Client) send(req *http.Request) ([]byte, error) {
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &apiError{Status: resp.StatusCode, Body: strings.TrimSpace(string(data))}
	}
	return data, nil
}
