package appstoreconnect

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// SearchTermEntry represents a single keyword row from an analytics report
type SearchTermEntry struct {
	Date        string `json:"date"`
	Keyword     string `json:"keyword"`
	Impressions int    `json:"impressions"`
	PageViews   int    `json:"page_views"`
	Installs    int    `json:"installs"`
}

type reportRequestResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

type reportsListResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Name     string `json:"name"`
			Category string `json:"category"`
		} `json:"attributes"`
	} `json:"data"`
}

type instancesListResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Granularity    string `json:"granularity"`
			ProcessingDate string `json:"processingDate"`
		} `json:"attributes"`
	} `json:"data"`
}

type segmentsListResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			URL            string `json:"url"`
			ExpirationDate string `json:"expirationDate"`
			Size           int    `json:"size"`
		} `json:"attributes"`
	} `json:"data"`
}

// CreateAnalyticsReportRequest creates a ONE_TIME_SNAPSHOT analytics report request.
// Returns the request ID to be polled later.
func (c *Client) CreateAnalyticsReportRequest(ctx context.Context, appStoreID string) (string, error) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "analyticsReportRequests",
			"attributes": map[string]interface{}{
				"accessType": "ONE_TIME_SNAPSHOT",
			},
			"relationships": map[string]interface{}{
				"app": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   appStoreID,
						"type": "apps",
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/analyticsReportRequests", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create report request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("create report request returned %d: %s", resp.StatusCode, string(b))
	}

	var result reportRequestResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode report request response: %w", err)
	}
	return result.Data.ID, nil
}

// GetSearchTermsFromRequest downloads search term data for a given report request ID.
// Returns nil, nil if the report is not yet ready.
func (c *Client) GetSearchTermsFromRequest(ctx context.Context, requestID string) ([]SearchTermEntry, error) {
	// List reports for this request
	resp, err := c.doRequest(ctx, http.MethodGet, "/analyticsReportRequests/"+requestID+"/reports", nil)
	if err != nil {
		return nil, fmt.Errorf("list reports: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // not ready yet
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list reports returned %d: %s", resp.StatusCode, string(b))
	}

	var reports reportsListResponse
	if err := json.NewDecoder(resp.Body).Decode(&reports); err != nil {
		return nil, fmt.Errorf("decode reports: %w", err)
	}
	if len(reports.Data) == 0 {
		return nil, nil // not ready yet
	}

	// Find APP_STORE_ENGAGEMENT report
	var reportID string
	for _, r := range reports.Data {
		if strings.Contains(r.Attributes.Category, "APP_STORE") || strings.Contains(r.Attributes.Name, "App Store") {
			reportID = r.ID
			break
		}
	}
	if reportID == "" {
		reportID = reports.Data[0].ID
	}

	// List instances of this report
	entries, err := c.downloadSearchTerms(ctx, reportID)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (c *Client) downloadSearchTerms(ctx context.Context, reportID string) ([]SearchTermEntry, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/analyticsReports/"+reportID+"/instances", nil)
	if err != nil {
		return nil, fmt.Errorf("list instances: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list instances returned %d: %s", resp.StatusCode, string(b))
	}

	var instances instancesListResponse
	if err := json.NewDecoder(resp.Body).Decode(&instances); err != nil {
		return nil, fmt.Errorf("decode instances: %w", err)
	}
	if len(instances.Data) == 0 {
		return nil, nil
	}

	var all []SearchTermEntry
	for _, inst := range instances.Data {
		entries, err := c.downloadInstanceSegments(ctx, inst.ID, inst.Attributes.ProcessingDate)
		if err != nil {
			continue // skip failed instances
		}
		all = append(all, entries...)
	}
	return all, nil
}

func (c *Client) downloadInstanceSegments(ctx context.Context, instanceID, date string) ([]SearchTermEntry, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/analyticsReportInstances/"+instanceID+"/segments", nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list segments returned %d", resp.StatusCode)
	}

	var segs segmentsListResponse
	if err := json.NewDecoder(resp.Body).Decode(&segs); err != nil {
		return nil, err
	}

	var all []SearchTermEntry
	for _, seg := range segs.Data {
		if seg.Attributes.URL == "" {
			continue
		}
		entries, err := c.downloadAndParseTSV(ctx, seg.Attributes.URL, date)
		if err != nil {
			continue
		}
		all = append(all, entries...)
	}
	return all, nil
}

// downloadAndParseTSV fetches a gzipped TSV segment and extracts search term rows.
func (c *Client) downloadAndParseTSV(ctx context.Context, segURL, date string) ([]SearchTermEntry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, segURL, nil)
	if err != nil {
		return nil, err
	}

	// Segment URLs are pre-signed and don't need the JWT
	dlResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = dlResp.Body.Close() }()

	if dlResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download segment returned %d", dlResp.StatusCode)
	}

	// Decompress gzip
	gz, err := gzip.NewReader(dlResp.Body)
	if err != nil {
		// Try reading as plain text if not gzipped
		gz = nil
	}

	var reader io.Reader
	if gz != nil {
		defer func() { _ = gz.Close() }()
		reader = gz
	} else {
		// Re-fetch without gzip reader (body already consumed on error)
		return nil, fmt.Errorf("expected gzipped content")
	}

	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return parseTSVSearchTerms(string(raw), date), nil
}

// parseTSVSearchTerms parses the TSV content and extracts search term data.
// The TSV format varies by report type; we look for a "Search Term" column
// and Source Type = "App Store Search".
func parseTSVSearchTerms(tsv, fallbackDate string) []SearchTermEntry {
	lines := strings.Split(tsv, "\n")
	if len(lines) < 2 {
		return nil
	}

	// Parse header
	header := strings.Split(lines[0], "\t")
	colIdx := map[string]int{}
	for i, h := range header {
		colIdx[strings.TrimSpace(h)] = i
	}

	dateIdx := colIdx["Date"]
	termIdx := colIdx["Search Term"]
	impIdx := colIdx["Impressions"]
	pvIdx := colIdx["Page Views"]
	instIdx := colIdx["App Units"]
	srcIdx := colIdx["Source Type"]

	// Search Term column must exist
	if _, ok := colIdx["Search Term"]; !ok {
		return nil
	}

	seen := map[string]int{} // keyword -> index in result
	var result []SearchTermEntry

	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		cols := strings.Split(line, "\t")
		if len(cols) <= termIdx {
			continue
		}

		// Filter: only rows with a non-empty search term
		keyword := strings.TrimSpace(cols[termIdx])
		if keyword == "" {
			continue
		}

		// Filter by Source Type = "App Store Search" if column exists
		if srcIdx > 0 && srcIdx < len(cols) {
			src := strings.TrimSpace(cols[srcIdx])
			if src != "" && !strings.Contains(strings.ToLower(src), "search") {
				continue
			}
		}

		date := fallbackDate
		if dateIdx < len(cols) {
			d := strings.TrimSpace(cols[dateIdx])
			if d != "" {
				// Normalize date format
				if t, err := time.Parse("2006-01-02", d); err == nil {
					date = t.Format("2006-01-02")
				} else if t, err := time.Parse("01/02/2006", d); err == nil {
					date = t.Format("2006-01-02")
				}
			}
		}

		imp := parseInt(safeCol(cols, impIdx))
		pv := parseInt(safeCol(cols, pvIdx))
		inst := parseInt(safeCol(cols, instIdx))

		key := date + "|" + keyword
		if idx, exists := seen[key]; exists {
			result[idx].Impressions += imp
			result[idx].PageViews += pv
			result[idx].Installs += inst
		} else {
			seen[key] = len(result)
			result = append(result, SearchTermEntry{
				Date:        date,
				Keyword:     keyword,
				Impressions: imp,
				PageViews:   pv,
				Installs:    inst,
			})
		}
	}

	return result
}

func safeCol(cols []string, idx int) string {
	if idx >= 0 && idx < len(cols) {
		return strings.TrimSpace(cols[idx])
	}
	return ""
}

func parseInt(s string) int {
	if s == "" {
		return 0
	}
	var n int
	_, _ = fmt.Sscanf(strings.ReplaceAll(s, ",", ""), "%d", &n)
	return n
}
