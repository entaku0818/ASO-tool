package scraper

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// generateTestECKey creates a test ECDSA P-256 private key in PEM format.
func generateTestECKey(t *testing.T) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate EC key: %v", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("failed to marshal private key: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
}

func TestNewSearchAdsClient_Valid(t *testing.T) {
	keyPEM := generateTestECKey(t)
	client, err := NewSearchAdsClient("client-id", "team-id", "key-id", keyPEM, "")
	if err != nil {
		t.Fatalf("NewSearchAdsClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("NewSearchAdsClient() returned nil")
	}
}

func TestNewSearchAdsClient_InvalidPEM(t *testing.T) {
	_, err := NewSearchAdsClient("client-id", "team-id", "key-id", []byte("not-valid-pem"), "")
	if err == nil {
		t.Error("NewSearchAdsClient() expected error for invalid PEM, got nil")
	}
}

func TestNewSearchAdsClient_EmptyPEM(t *testing.T) {
	_, err := NewSearchAdsClient("client-id", "team-id", "key-id", []byte(""), "")
	if err == nil {
		t.Error("NewSearchAdsClient() expected error for empty PEM, got nil")
	}
}

func TestNewSearchAdsClient_RSAKey(t *testing.T) {
	// RSA key should fail because we expect ECDSA
	rsaPEM := []byte(`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC7o4qne60TB8+F
-----END PRIVATE KEY-----`)
	_, err := NewSearchAdsClient("client-id", "team-id", "key-id", rsaPEM, "")
	if err == nil {
		t.Error("NewSearchAdsClient() expected error for invalid key format, got nil")
	}
}

func TestNewSearchAdsClient_WithOrgID(t *testing.T) {
	keyPEM := generateTestECKey(t)
	client, err := NewSearchAdsClient("client-id", "team-id", "key-id", keyPEM, "org-123")
	if err != nil {
		t.Fatalf("NewSearchAdsClient() error = %v", err)
	}
	if client.orgID != "org-123" {
		t.Errorf("orgID = %v, want 'org-123'", client.orgID)
	}
}

// redirectTransport rewrites request URLs so we can intercept calls to external APIs.
type redirectTransport struct {
	tokenServer    *httptest.Server
	searchAdsServer *httptest.Server
}

func (rt *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var targetBase string
	if strings.Contains(req.URL.Host, "appleid.apple.com") {
		targetBase = rt.tokenServer.URL
	} else if strings.Contains(req.URL.Host, "searchads.apple.com") {
		targetBase = rt.searchAdsServer.URL
	}

	if targetBase != "" {
		newURL := *req.URL
		newURL.Scheme = "http"
		newURL.Host = strings.TrimPrefix(targetBase, "http://")
		req = req.Clone(req.Context())
		req.URL = &newURL
	}

	return http.DefaultTransport.RoundTrip(req)
}

func TestSearchAdsClient_GetKeywordPopularity(t *testing.T) {
	// Mock token server
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "test-access-token",
			"expires_in":   3600,
		})
	}))
	defer tokenServer.Close()

	// Mock Search Ads API server
	searchAdsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-access-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(keywordSearchResponse{
			Data: []KeywordPopularity{
				{Text: "calculator", PopularityScore: 5},
				{Text: "math", PopularityScore: 3},
			},
		})
	}))
	defer searchAdsServer.Close()

	keyPEM := generateTestECKey(t)
	client, err := NewSearchAdsClient("client-id", "team-id", "key-id", keyPEM, "")
	if err != nil {
		t.Fatalf("NewSearchAdsClient() error = %v", err)
	}

	// Inject the redirect transport
	client.httpClient = &http.Client{
		Transport: &redirectTransport{
			tokenServer:    tokenServer,
			searchAdsServer: searchAdsServer,
		},
	}

	results, err := client.GetKeywordPopularity(context.Background(), 123456789, []string{"calculator", "math"}, "jp", 10)
	if err != nil {
		t.Fatalf("GetKeywordPopularity() error = %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("GetKeywordPopularity() returned %d results, want 2", len(results))
	}
	if results[0].Text != "calculator" {
		t.Errorf("results[0].Text = %v, want 'calculator'", results[0].Text)
	}
	if results[0].PopularityScore != 5 {
		t.Errorf("results[0].PopularityScore = %v, want 5", results[0].PopularityScore)
	}
	if results[1].PopularityScore != 3 {
		t.Errorf("results[1].PopularityScore = %v, want 3", results[1].PopularityScore)
	}
}

func TestSearchAdsClient_GetKeywordPopularity_DefaultLimit(t *testing.T) {
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "test-token",
			"expires_in":   3600,
		})
	}))
	defer tokenServer.Close()

	var capturedBody keywordSearchRequest
	searchAdsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&capturedBody)
		_ = json.NewEncoder(w).Encode(keywordSearchResponse{Data: []KeywordPopularity{}})
	}))
	defer searchAdsServer.Close()

	keyPEM := generateTestECKey(t)
	client, _ := NewSearchAdsClient("c", "t", "k", keyPEM, "")
	client.httpClient = &http.Client{
		Transport: &redirectTransport{
			tokenServer:    tokenServer,
			searchAdsServer: searchAdsServer,
		},
	}

	_, _ = client.GetKeywordPopularity(context.Background(), 1, nil, "", 0) // limit=0 should default to 25
	if capturedBody.Limit != 25 {
		t.Errorf("default limit = %v, want 25", capturedBody.Limit)
	}
}

func TestSearchAdsClient_GetKeywordPopularity_APIError(t *testing.T) {
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "test-token",
			"expires_in":   3600,
		})
	}))
	defer tokenServer.Close()

	searchAdsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer searchAdsServer.Close()

	keyPEM := generateTestECKey(t)
	client, _ := NewSearchAdsClient("c", "t", "k", keyPEM, "")
	client.httpClient = &http.Client{
		Transport: &redirectTransport{
			tokenServer:    tokenServer,
			searchAdsServer: searchAdsServer,
		},
	}

	_, err := client.GetKeywordPopularity(context.Background(), 1, []string{"test"}, "", 10)
	if err == nil {
		t.Error("GetKeywordPopularity() expected error for API 500, got nil")
	}
}

func TestKeywordPopularity_Struct(t *testing.T) {
	kp := KeywordPopularity{
		Text:            "calculator",
		PopularityScore: 4,
	}
	if kp.Text != "calculator" {
		t.Errorf("Text = %v, want 'calculator'", kp.Text)
	}
	if kp.PopularityScore != 4 {
		t.Errorf("PopularityScore = %v, want 4", kp.PopularityScore)
	}
}
