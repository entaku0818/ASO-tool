-- Tracked keywords table: keywords to track search results (not tied to a specific app)
CREATE TABLE IF NOT EXISTS tracked_keywords (
    id VARCHAR(36) PRIMARY KEY,
    keyword VARCHAR(255) NOT NULL,
    country VARCHAR(2) NOT NULL DEFAULT 'JP',
    platform VARCHAR(10) NOT NULL DEFAULT 'ios',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(keyword, country, platform)
);

CREATE INDEX idx_tracked_keywords_keyword ON tracked_keywords(keyword);

-- Search results table: stores all apps found in search results for tracked keywords
CREATE TABLE IF NOT EXISTS search_results (
    id VARCHAR(36) PRIMARY KEY,
    tracked_keyword_id VARCHAR(36) NOT NULL REFERENCES tracked_keywords(id) ON DELETE CASCADE,
    rank INTEGER NOT NULL,
    app_name VARCHAR(500) NOT NULL,
    bundle_id VARCHAR(255),
    developer VARCHAR(255),
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_search_results_tracked_keyword_id ON search_results(tracked_keyword_id);
CREATE INDEX idx_search_results_recorded_at ON search_results(recorded_at);
CREATE INDEX idx_search_results_rank ON search_results(rank);
