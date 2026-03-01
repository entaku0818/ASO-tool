package scraper

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	searchAdsTokenURL = "https://appleid.apple.com/auth/oauth2/token"
	searchAdsBaseURL  = "https://api.searchads.apple.com/api/v4"
	searchAdsAudience = "https://appleid.apple.com"
)

// SearchAdsClient is the Apple Search Ads API client
type SearchAdsClient struct {
	httpClient  *http.Client
	clientID    string
	teamID      string
	keyID       string
	privateKey  *ecdsa.PrivateKey
	accessToken string
	tokenExpiry time.Time
	orgID       string
}

// KeywordPopularity holds a keyword and its popularity score (0–5)
type KeywordPopularity struct {
	Text            string `json:"text"`
	PopularityScore int    `json:"popularityScore"`
}

// NewSearchAdsClient creates a new Apple Search Ads API client
func NewSearchAdsClient(clientID, teamID, keyID string, privateKeyPEM []byte, orgID string) (*SearchAdsClient, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not ECDSA type")
	}

	return &SearchAdsClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		clientID:   clientID,
		teamID:     teamID,
		keyID:      keyID,
		privateKey: ecdsaKey,
		orgID:      orgID,
	}, nil
}

// generateClientSecret creates a JWT client secret for OAuth 2.0
func (c *SearchAdsClient) generateClientSecret() (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": c.clientID,
		"iss": c.teamID,
		"aud": searchAdsAudience,
		"iat": now.Unix(),
		"exp": now.Add(180 * 24 * time.Hour).Unix(), // 180 days max
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = c.keyID

	return token.SignedString(c.privateKey)
}

// refreshToken fetches a new OAuth 2.0 access token
func (c *SearchAdsClient) refreshToken(ctx context.Context) error {
	clientSecret, err := c.generateClientSecret()
	if err != nil {
		return fmt.Errorf("failed to generate client secret: %w", err)
	}

	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", c.clientID)
	formData.Set("client_secret", clientSecret)
	formData.Set("scope", "searchadsorg")

	req, err := http.NewRequestWithContext(ctx, "POST", searchAdsTokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token request failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return nil
}

// ensureToken ensures we have a valid access token
func (c *SearchAdsClient) ensureToken(ctx context.Context) error {
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry.Add(-60*time.Second)) {
		return nil
	}
	return c.refreshToken(ctx)
}

// doRequest performs an authenticated API request
func (c *SearchAdsClient) doRequest(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, searchAdsBaseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	if c.orgID != "" {
		req.Header.Set("X-AP-Context", fmt.Sprintf("orgId=%s", c.orgID))
	}

	return c.httpClient.Do(req)
}

type keywordSearchRequest struct {
	AdamID   int64    `json:"adamId"`
	Keywords []string `json:"keywords,omitempty"`
	Limit    int      `json:"limit,omitempty"`
}

type keywordSearchResponse struct {
	Data []KeywordPopularity `json:"data"`
}

// GetKeywordPopularity fetches keyword popularity scores from Search Ads API.
// adamID is the app's iTunes numeric ID. If keywords is empty, returns suggestions.
// limit controls the max number of results (default 25).
func (c *SearchAdsClient) GetKeywordPopularity(ctx context.Context, adamID int64, keywords []string, limit int) ([]KeywordPopularity, error) {
	if limit <= 0 {
		limit = 25
	}

	reqBody := keywordSearchRequest{
		AdamID:   adamID,
		Keywords: keywords,
		Limit:    limit,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, "POST", "/keywords/search", bodyBytes)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search ads API error: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result keywordSearchResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Data, nil
}
