package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GooglePlayScraper implements Scraper for Google Play Store
// Note: Google Play doesn't have an official API, so this uses limited public endpoints
type GooglePlayScraper struct {
	client *http.Client
}

// NewGooglePlayScraper creates a new Google Play scraper
func NewGooglePlayScraper() *GooglePlayScraper {
	return &GooglePlayScraper{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetAppInfo fetches app information from Google Play
// Uses web scraping as there's no official API
func (s *GooglePlayScraper) GetAppInfo(ctx context.Context, bundleID string, country string) (*AppInfo, error) {
	if country == "" {
		country = "jp"
	}

	appURL := fmt.Sprintf(
		"https://play.google.com/store/apps/details?id=%s&hl=%s&gl=%s",
		url.QueryEscape(bundleID),
		url.QueryEscape(country),
		url.QueryEscape(country),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", appURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept-Language", country)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch app info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("app not found: %s", bundleID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return s.parseAppPage(bundleID, string(body), appURL)
}

// SearchKeyword searches for apps on Google Play
// Note: This is limited as Google Play doesn't provide a public search API
func (s *GooglePlayScraper) SearchKeyword(ctx context.Context, keyword string, country string, limit int) ([]SearchResult, error) {
	if country == "" {
		country = "jp"
	}

	searchURL := fmt.Sprintf(
		"https://play.google.com/store/search?q=%s&c=apps&hl=%s&gl=%s",
		url.QueryEscape(keyword),
		url.QueryEscape(country),
		url.QueryEscape(country),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search apps: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return s.parseSearchResults(string(body), limit)
}

// GetReviews fetches reviews from Google Play
func (s *GooglePlayScraper) GetReviews(ctx context.Context, appID string, country string, page int) ([]Review, error) {
	// Google Play doesn't have a public RSS feed like App Store
	// This would require more complex scraping
	// For now, return empty results
	return []Review{}, nil
}

// GetAppRanking finds the rank of an app for a specific keyword
func (s *GooglePlayScraper) GetAppRanking(ctx context.Context, bundleID string, keyword string, country string) (*int, error) {
	results, err := s.SearchKeyword(ctx, keyword, country, 100)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		if result.AppInfo.BundleID == bundleID {
			rank := result.Rank
			return &rank, nil
		}
	}

	return nil, nil
}

// parseAppPage extracts app info from Google Play HTML page
func (s *GooglePlayScraper) parseAppPage(bundleID, html, storeURL string) (*AppInfo, error) {
	app := &AppInfo{
		BundleID: bundleID,
		StoreURL: storeURL,
	}

	// Extract app name from title tag
	titleRe := regexp.MustCompile(`<title>([^<]+)\s*-\s*Apps on Google Play</title>`)
	if matches := titleRe.FindStringSubmatch(html); len(matches) > 1 {
		app.Name = strings.TrimSpace(matches[1])
	}

	// Try to extract from JSON-LD data if available
	jsonLdRe := regexp.MustCompile(`<script type="application/ld\+json"[^>]*>([^<]+)</script>`)
	if matches := jsonLdRe.FindStringSubmatch(html); len(matches) > 1 {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(matches[1]), &jsonData); err == nil {
			if name, ok := jsonData["name"].(string); ok && app.Name == "" {
				app.Name = name
			}
			if desc, ok := jsonData["description"].(string); ok {
				app.Description = desc
			}
			if author, ok := jsonData["author"].(map[string]interface{}); ok {
				if name, ok := author["name"].(string); ok {
					app.Developer = name
				}
			}
			if rating, ok := jsonData["aggregateRating"].(map[string]interface{}); ok {
				if ratingValue, ok := rating["ratingValue"].(float64); ok {
					app.Rating = ratingValue
				}
				if ratingCount, ok := rating["ratingCount"].(float64); ok {
					app.RatingCount = int(ratingCount)
				}
			}
		}
	}

	// Extract icon URL
	iconRe := regexp.MustCompile(`<img[^>]+itemprop="image"[^>]+src="([^"]+)"`)
	if matches := iconRe.FindStringSubmatch(html); len(matches) > 1 {
		app.IconURL = matches[1]
	}

	return app, nil
}

// parseSearchResults extracts search results from Google Play HTML
func (s *GooglePlayScraper) parseSearchResults(html string, limit int) ([]SearchResult, error) {
	results := []SearchResult{}

	// Look for app package IDs in search results
	// Pattern: /store/apps/details?id=com.example.app
	appIdRe := regexp.MustCompile(`/store/apps/details\?id=([a-zA-Z0-9_.]+)`)
	matches := appIdRe.FindAllStringSubmatch(html, -1)

	seen := make(map[string]bool)
	rank := 1

	for _, match := range matches {
		if len(match) > 1 {
			bundleID := match[1]
			if seen[bundleID] {
				continue
			}
			seen[bundleID] = true

			results = append(results, SearchResult{
				Rank: rank,
				AppInfo: AppInfo{
					BundleID: bundleID,
					StoreURL: fmt.Sprintf("https://play.google.com/store/apps/details?id=%s", bundleID),
				},
			})
			rank++

			if limit > 0 && rank > limit {
				break
			}
		}
	}

	// Try to extract more details for each app
	for i := range results {
		s.enrichAppFromHTML(html, &results[i].AppInfo)
	}

	return results, nil
}

// enrichAppFromHTML tries to extract additional app info from HTML
func (s *GooglePlayScraper) enrichAppFromHTML(html string, app *AppInfo) {
	// Look for app name near the bundle ID
	// This is a simplified extraction - a real implementation would need more sophisticated parsing
	pattern := regexp.MustCompile(
		fmt.Sprintf(`id=%s[^"]*"[^>]*>[^<]*</a>[^<]*<[^>]+>([^<]+)`, regexp.QuoteMeta(app.BundleID)),
	)
	if matches := pattern.FindStringSubmatch(html); len(matches) > 1 {
		app.Name = strings.TrimSpace(matches[1])
	}
}

// parseRating parses rating string to float
func parseRating(s string) float64 {
	s = strings.TrimSpace(s)
	rating, _ := strconv.ParseFloat(s, 64)
	return rating
}

// parseRatingCount parses rating count string to int
func parseRatingCount(s string) int {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, " ", "")
	count, _ := strconv.Atoi(s)
	return count
}
