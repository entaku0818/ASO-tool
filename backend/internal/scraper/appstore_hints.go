package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

type KeywordSuggestion struct {
	Term string `json:"term"`
}

var appStoreHintsClient = &http.Client{Timeout: 5 * time.Second}

// storeFronts maps country codes to Apple Store Front IDs
var storeFronts = map[string]string{
	"jp": "143462-1,29",
	"us": "143441-1,29",
	"gb": "143444-1,29",
	"de": "143443-1,29",
	"kr": "143466-1,29",
	"cn": "143465-1,29",
	"fr": "143442-1,29",
}

// termPattern matches <key>term</key> followed by <string>VALUE</string> in Apple plist
var termPattern = regexp.MustCompile(`<key>term</key>\s*<string>([^<]+)</string>`)

// FetchKeywordSuggestions fetches App Store search suggestions for a given term.
func FetchKeywordSuggestions(ctx context.Context, term, country string) ([]KeywordSuggestion, error) {
	storeFront, ok := storeFronts[country]
	if !ok {
		storeFront = storeFronts["us"]
	}

	apiURL := fmt.Sprintf(
		"https://search.itunes.apple.com/WebObjects/MZSearchHints.woa/wa/hints?clientApplication=Software&term=%s&count=25",
		url.QueryEscape(term),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "AppStore/3.0 iOS/17.0 model/iPhone14,3 hwp/t8110 build/21A329 (6; dt:285)")
	req.Header.Set("X-Apple-Store-Front", storeFront)
	req.Header.Set("Accept-Language", country+"-"+country)

	resp, err := appStoreHintsClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hints request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hints API returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	matches := termPattern.FindAllSubmatch(body, -1)
	suggestions := make([]KeywordSuggestion, 0, len(matches))
	for _, m := range matches {
		suggestions = append(suggestions, KeywordSuggestion{Term: string(m[1])})
	}

	return suggestions, nil
}
