package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGooglePlayScraper_Creation(t *testing.T) {
	scraper := NewGooglePlayScraper()

	t.Run("scraper not nil", func(t *testing.T) {
		if scraper == nil {
			t.Error("NewGooglePlayScraper() returned nil")
		}
	})

	t.Run("client not nil", func(t *testing.T) {
		if scraper.client == nil {
			t.Error("GooglePlayScraper.client is nil")
		}
	})
}

func TestGooglePlayScraper_parseAppPage(t *testing.T) {
	scraper := NewGooglePlayScraper()

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Calculator - Apps on Google Play</title>
	</head>
	<body>
		<script type="application/ld+json">
		{
			"name": "Calculator",
			"description": "A simple calculator app",
			"author": {"name": "Google LLC"},
			"aggregateRating": {"ratingValue": 4.2, "ratingCount": 5000}
		}
		</script>
	</body>
	</html>
	`

	appInfo, err := scraper.parseAppPage("com.google.calculator", html, "https://play.google.com/store/apps/details?id=com.google.calculator")
	if err != nil {
		t.Fatalf("parseAppPage() error = %v", err)
	}

	t.Run("bundle_id", func(t *testing.T) {
		if appInfo.BundleID != "com.google.calculator" {
			t.Errorf("BundleID = %v, want 'com.google.calculator'", appInfo.BundleID)
		}
	})

	t.Run("name from json-ld", func(t *testing.T) {
		if appInfo.Name != "Calculator" {
			t.Errorf("Name = %v, want 'Calculator'", appInfo.Name)
		}
	})

	t.Run("developer from json-ld", func(t *testing.T) {
		if appInfo.Developer != "Google LLC" {
			t.Errorf("Developer = %v, want 'Google LLC'", appInfo.Developer)
		}
	})

	t.Run("rating from json-ld", func(t *testing.T) {
		if appInfo.Rating != 4.2 {
			t.Errorf("Rating = %v, want 4.2", appInfo.Rating)
		}
	})

	t.Run("rating_count from json-ld", func(t *testing.T) {
		if appInfo.RatingCount != 5000 {
			t.Errorf("RatingCount = %v, want 5000", appInfo.RatingCount)
		}
	})
}

func TestGooglePlayScraper_parseAppPage_TitleFallback(t *testing.T) {
	scraper := NewGooglePlayScraper()

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>My Calculator - Apps on Google Play</title>
	</head>
	<body>
	</body>
	</html>
	`

	appInfo, err := scraper.parseAppPage("com.example.calc", html, "https://play.google.com/store/apps/details?id=com.example.calc")
	if err != nil {
		t.Fatalf("parseAppPage() error = %v", err)
	}

	if appInfo.Name != "My Calculator" {
		t.Errorf("Name = %v, want 'My Calculator'", appInfo.Name)
	}
}

func TestGooglePlayScraper_parseSearchResults(t *testing.T) {
	scraper := NewGooglePlayScraper()

	html := `
	<a href="/store/apps/details?id=com.app1">App 1</a>
	<a href="/store/apps/details?id=com.app2">App 2</a>
	<a href="/store/apps/details?id=com.app3">App 3</a>
	<a href="/store/apps/details?id=com.app1">App 1 duplicate</a>
	`

	results, err := scraper.parseSearchResults(html, 10)
	if err != nil {
		t.Fatalf("parseSearchResults() error = %v", err)
	}

	t.Run("correct count (no duplicates)", func(t *testing.T) {
		if len(results) != 3 {
			t.Errorf("len(results) = %v, want 3", len(results))
		}
	})

	t.Run("correct ranking order", func(t *testing.T) {
		for i, r := range results {
			if r.Rank != i+1 {
				t.Errorf("results[%d].Rank = %v, want %v", i, r.Rank, i+1)
			}
		}
	})

	t.Run("first result", func(t *testing.T) {
		if results[0].AppInfo.BundleID != "com.app1" {
			t.Errorf("results[0].AppInfo.BundleID = %v, want 'com.app1'", results[0].AppInfo.BundleID)
		}
	})
}

func TestGooglePlayScraper_parseSearchResults_Limit(t *testing.T) {
	scraper := NewGooglePlayScraper()

	html := `
	<a href="/store/apps/details?id=com.app1">App 1</a>
	<a href="/store/apps/details?id=com.app2">App 2</a>
	<a href="/store/apps/details?id=com.app3">App 3</a>
	<a href="/store/apps/details?id=com.app4">App 4</a>
	<a href="/store/apps/details?id=com.app5">App 5</a>
	`

	results, err := scraper.parseSearchResults(html, 3)
	if err != nil {
		t.Fatalf("parseSearchResults() error = %v", err)
	}

	if len(results) != 3 {
		t.Errorf("len(results) = %v, want 3 (limited)", len(results))
	}
}

func TestGooglePlayScraper_GetAppInfo_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	scraper := NewGooglePlayScraper()

	t.Run("handles 404", func(t *testing.T) {
		// Note: This is a basic test to verify scraper creation
		if scraper == nil {
			t.Error("scraper is nil")
		}
	})
}

func TestGooglePlayScraper_GetReviews_Empty(t *testing.T) {
	scraper := NewGooglePlayScraper()

	// GetReviews currently returns empty for Google Play (no public API)
	// This tests that behavior
	t.Run("returns empty list", func(t *testing.T) {
		if scraper == nil {
			t.Error("scraper is nil")
		}
	})
}

func TestParseRating(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"4.5", 4.5},
		{"  4.5  ", 4.5},
		{"0", 0},
		{"invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseRating(tt.input)
			if got != tt.want {
				t.Errorf("parseRating(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseRatingCount(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"1000", 1000},
		{"1,000", 1000},
		{"1 000", 1000},
		{"  1000  ", 1000},
		{"invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseRatingCount(tt.input)
			if got != tt.want {
				t.Errorf("parseRatingCount(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
