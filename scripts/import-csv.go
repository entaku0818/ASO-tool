package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const apiURL = "https://aso-api-671942133800.asia-northeast1.run.app"

type App struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	BundleID string `json:"bundle_id"`
	Platform string `json:"platform"`
}

type Keyword struct {
	ID      string `json:"id"`
	Keyword string `json:"keyword"`
	Country string `json:"country"`
}

type CSVRow struct {
	AppName    string
	AppID      string
	Platform   string
	Keyword    string
	Country    string
	Ranking    int
	Popularity int
	Difficulty int
	LastUpdate string
}

var token string

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run import-csv.go <csv-file> <auth-token>")
		fmt.Println("  Get auth token by logging in at the frontend")
		os.Exit(1)
	}

	csvFile := os.Args[1]
	token = os.Args[2]

	rows, err := readCSV(csvFile)
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d rows\n", len(rows))

	// Group by app
	appRows := make(map[string][]CSVRow)
	for _, row := range rows {
		key := row.AppID
		appRows[key] = append(appRows[key], row)
	}

	for appID, rows := range appRows {
		fmt.Printf("\nProcessing app: %s (%s)\n", rows[0].AppName, appID)

		// Find or create app
		app, err := findOrCreateApp(rows[0])
		if err != nil {
			fmt.Printf("  Error creating app: %v\n", err)
			continue
		}
		fmt.Printf("  App ID: %s\n", app.ID)

		// Create keywords and rankings
		for _, row := range rows {
			if row.Ranking >= 1000 {
				fmt.Printf("  Skipping '%s' (rank >= 1000)\n", row.Keyword)
				continue
			}

			keyword, err := findOrCreateKeyword(app.ID, row)
			if err != nil {
				fmt.Printf("  Error creating keyword '%s': %v\n", row.Keyword, err)
				continue
			}

			err = createRanking(app.ID, keyword.ID, row)
			if err != nil {
				fmt.Printf("  Error creating ranking for '%s': %v\n", row.Keyword, err)
				continue
			}

			fmt.Printf("  Imported: %s -> rank %d\n", row.Keyword, row.Ranking)
		}
	}

	fmt.Println("\nImport completed!")
}

func readCSV(filename string) ([]CSVRow, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var rows []CSVRow
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		ranking, _ := strconv.Atoi(record[8])
		popularity, _ := strconv.Atoi(record[10])
		difficulty, _ := strconv.Atoi(record[11])

		platform := record[2]
		if platform == "iphone" {
			platform = "ios"
		}

		rows = append(rows, CSVRow{
			AppName:    record[0],
			AppID:      record[1],
			Platform:   platform,
			Keyword:    record[3],
			Country:    record[4],
			Ranking:    ranking,
			Popularity: popularity,
			Difficulty: difficulty,
			LastUpdate: record[7],
		})
	}

	return rows, nil
}

func findOrCreateApp(row CSVRow) (*App, error) {
	// Try to find existing app
	resp, err := doRequest("GET", fmt.Sprintf("/api/apps"), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apps []App
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		return nil, err
	}

	// Look for existing app by bundle_id
	bundleID := fmt.Sprintf("id%s", row.AppID)
	for _, app := range apps {
		if app.BundleID == bundleID || strings.Contains(app.BundleID, row.AppID) {
			return &app, nil
		}
	}

	// Create new app
	body := map[string]interface{}{
		"name":      row.AppName,
		"bundle_id": bundleID,
		"platform":  row.Platform,
	}
	jsonBody, _ := json.Marshal(body)

	resp, err = doRequest("POST", "/api/apps", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var app App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, err
	}

	return &app, nil
}

func findOrCreateKeyword(appID string, row CSVRow) (*Keyword, error) {
	// Get existing keywords
	resp, err := doRequest("GET", fmt.Sprintf("/api/apps/%s/keywords", appID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var keywords []Keyword
	if err := json.NewDecoder(resp.Body).Decode(&keywords); err != nil {
		return nil, err
	}

	// Look for existing keyword
	for _, kw := range keywords {
		if kw.Keyword == row.Keyword && kw.Country == row.Country {
			return &kw, nil
		}
	}

	// Create new keyword
	body := map[string]interface{}{
		"keyword": row.Keyword,
		"country": row.Country,
	}
	jsonBody, _ := json.Marshal(body)

	resp, err = doRequest("POST", fmt.Sprintf("/api/apps/%s/keywords", appID), bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var keyword Keyword
	if err := json.NewDecoder(resp.Body).Decode(&keyword); err != nil {
		return nil, err
	}

	return &keyword, nil
}

func createRanking(appID, keywordID string, row CSVRow) error {
	body := map[string]interface{}{
		"app_id":     appID,
		"keyword_id": keywordID,
		"rank":       row.Ranking,
		"checked_at": time.Now().Format(time.RFC3339),
	}
	jsonBody, _ := json.Marshal(body)

	resp, err := doRequest("POST", "/api/rankings", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func doRequest(method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, apiURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 30 * time.Second}
	return client.Do(req)
}
