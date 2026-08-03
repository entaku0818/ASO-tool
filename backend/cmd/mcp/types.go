package main

import "net/url"

// --- tool arguments (jsonschema tags drive the generated input schema) ---

type AppIDArgs struct {
	AppID string `json:"app_id" jsonschema:"list_apps が返すアプリのUUID"`
}

type HistoryArgs struct {
	AppID string `json:"app_id" jsonschema:"list_apps が返すアプリのUUID"`
	Days  int    `json:"days,omitempty" jsonschema:"さかのぼる日数"`
}

type SearchArgs struct {
	Keyword string `json:"keyword" jsonschema:"検索するキーワード"`
	Country string `json:"country,omitempty" jsonschema:"国コード(既定 jp)"`
	Limit   int    `json:"limit,omitempty" jsonschema:"取得件数(既定 50)"`
}

// --- API response shapes (subset of the fields the backend returns) ---

type App struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	BundleID string `json:"bundle_id"`
	Platform string `json:"platform"`
	StoreURL string `json:"store_url,omitempty"`
}

type KeywordRank struct {
	KeywordID    string `json:"keyword_id"`
	Keyword      string `json:"keyword"`
	Country      string `json:"country"`
	CurrentRank  *int   `json:"current_rank"`
	PreviousRank *int   `json:"previous_rank"`
	Change       *int   `json:"change"`
}

type RankingRow struct {
	KeywordID  string `json:"keyword_id"`
	Keyword    string `json:"keyword"`
	Country    string `json:"country"`
	Rank       *int   `json:"rank"`
	RecordedAt string `json:"recorded_at"`
}

type RisingKeyword struct {
	KeywordID    string `json:"keyword_id"`
	Keyword      string `json:"keyword"`
	Country      string `json:"country"`
	CurrentRank  int    `json:"current_rank"`
	PreviousRank int    `json:"previous_rank"`
	Improvement  int    `json:"improvement"`
}

// SearchResult mirrors the scraper endpoint, which returns PascalCase keys.
type SearchResult struct {
	Rank    int `json:"Rank"`
	AppInfo struct {
		BundleID  string  `json:"BundleID"`
		Name      string  `json:"Name"`
		Developer string  `json:"Developer"`
		Rating    float64 `json:"Rating"`
		StoreURL  string  `json:"StoreURL"`
	} `json:"AppInfo"`
}

func urlQueryEscape(s string) string { return url.QueryEscape(s) }
