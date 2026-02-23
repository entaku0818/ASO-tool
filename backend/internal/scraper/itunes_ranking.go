package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AppRankingEntry represents a single app in a ranking
type AppRankingEntry struct {
	Rank        int    `json:"rank"`
	Name        string `json:"name"`
	Developer   string `json:"developer"`
	IconURL     string `json:"icon_url"`
	Category    string `json:"category"`
	StoreURL    string `json:"store_url"`
	AppID       string `json:"app_id"`
	Price       string `json:"price"`
	ReleaseDate string `json:"release_date"`
}

// iTunesRSSResponse — entry is RawMessage because iTunes returns object when 1 result, array otherwise
type iTunesRSSResponse struct {
	Feed struct {
		Entry json.RawMessage `json:"entry"`
	} `json:"feed"`
}

type iTunesRSSEntry struct {
	ImName struct {
		Label string `json:"label"`
	} `json:"im:name"`
	ImImage []struct {
		Label      string `json:"label"`
		Attributes struct {
			Height string `json:"height"`
		} `json:"attributes"`
	} `json:"im:image"`
	ImPrice struct {
		Label      string `json:"label"`
		Attributes struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"attributes"`
	} `json:"im:price"`
	// link is also object or array depending on count
	Link json.RawMessage `json:"link"`
	ID   struct {
		Label      string `json:"label"`
		Attributes struct {
			ImID       string `json:"im:id"`
			ImBundleID string `json:"im:bundleId"`
		} `json:"attributes"`
	} `json:"id"`
	ImArtist struct {
		Label string `json:"label"`
	} `json:"im:artist"`
	Category struct {
		Attributes struct {
			Term string `json:"term"`
		} `json:"attributes"`
	} `json:"category"`
	ImReleaseDate struct {
		Label string `json:"label"`
	} `json:"im:releaseDate"`
}

type itunesLink struct {
	Attributes struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"attributes"`
}

// extractStoreURL handles link as either object or array
func extractStoreURL(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return ""
	}

	if trimmed[0] == '[' {
		var links []itunesLink
		if err := json.Unmarshal(trimmed, &links); err == nil {
			for _, l := range links {
				if l.Attributes.Rel == "alternate" {
					return l.Attributes.Href
				}
			}
		}
	} else {
		var link itunesLink
		if err := json.Unmarshal(trimmed, &link); err == nil {
			if link.Attributes.Rel == "alternate" {
				return link.Attributes.Href
			}
		}
	}
	return ""
}

// parseEntries handles entry as either single object or array
func parseEntries(raw json.RawMessage) ([]iTunesRSSEntry, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	trimmed := bytes.TrimSpace(raw)
	if trimmed[0] == '[' {
		var entries []iTunesRSSEntry
		if err := json.Unmarshal(trimmed, &entries); err != nil {
			return nil, err
		}
		return entries, nil
	}
	var entry iTunesRSSEntry
	if err := json.Unmarshal(trimmed, &entry); err != nil {
		return nil, err
	}
	return []iTunesRSSEntry{entry}, nil
}

// ITunesRankingScraper fetches app rankings from the iTunes RSS Feed API
type ITunesRankingScraper struct {
	client *http.Client
}

// NewITunesRankingScraper creates a new iTunes ranking scraper
func NewITunesRankingScraper() *ITunesRankingScraper {
	return &ITunesRankingScraper{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchRanking fetches the app ranking from iTunes RSS Feed
// rankingType: topfreeapplications, toppaidapplications, topgrossingapplications, newfreeapplications, newpaidapplications
// country: jp, us, gb, etc.
// genreID: "" for all genres, "6014" for games, etc.
// limit: max 200
func (s *ITunesRankingScraper) FetchRanking(ctx context.Context, country, rankingType, genreID string, limit int) ([]AppRankingEntry, error) {
	if country == "" {
		country = "jp"
	}
	if rankingType == "" {
		rankingType = "topfreeapplications"
	}
	if limit <= 0 || limit > 200 {
		limit = 100
	}

	var rssURL string
	if genreID != "" {
		rssURL = fmt.Sprintf(
			"https://itunes.apple.com/%s/rss/%s/limit=%d/genre=%s/json",
			country, rankingType, limit, genreID,
		)
	} else {
		rssURL = fmt.Sprintf(
			"https://itunes.apple.com/%s/rss/%s/limit=%d/json",
			country, rankingType, limit,
		)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", rssURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ranking: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rssResp iTunesRSSResponse
	if err := json.NewDecoder(resp.Body).Decode(&rssResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	rawEntries, err := parseEntries(rssResp.Feed.Entry)
	if err != nil {
		return nil, fmt.Errorf("failed to parse entries: %w", err)
	}

	entries := make([]AppRankingEntry, 0, len(rawEntries))
	for i, entry := range rawEntries {
		iconURL := ""
		for _, img := range entry.ImImage {
			iconURL = img.Label
		}

		price := entry.ImPrice.Label
		if price == "" {
			price = entry.ImPrice.Attributes.Amount
		}

		entries = append(entries, AppRankingEntry{
			Rank:        i + 1,
			Name:        entry.ImName.Label,
			Developer:   entry.ImArtist.Label,
			IconURL:     iconURL,
			Category:    entry.Category.Attributes.Term,
			StoreURL:    extractStoreURL(entry.Link),
			AppID:       entry.ID.Attributes.ImID,
			Price:       price,
			ReleaseDate: entry.ImReleaseDate.Label,
		})
	}

	return entries, nil
}
