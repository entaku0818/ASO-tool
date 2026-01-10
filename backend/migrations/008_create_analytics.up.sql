-- App Store Connect analytics data (impressions and downloads)
CREATE TABLE IF NOT EXISTS asc_analytics (
    id VARCHAR(36) PRIMARY KEY,
    app_id VARCHAR(36) NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    impressions INTEGER NOT NULL DEFAULT 0,
    downloads INTEGER NOT NULL DEFAULT 0,
    page_views INTEGER DEFAULT 0,
    conversion_rate DECIMAL(5,4),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(app_id, date)
);

CREATE INDEX IF NOT EXISTS idx_asc_analytics_app_id ON asc_analytics(app_id);
CREATE INDEX IF NOT EXISTS idx_asc_analytics_date ON asc_analytics(date);
CREATE INDEX IF NOT EXISTS idx_asc_analytics_app_date ON asc_analytics(app_id, date DESC);

-- App version history for correlation analysis
CREATE TABLE IF NOT EXISTS app_versions (
    id VARCHAR(36) PRIMARY KEY,
    app_id VARCHAR(36) NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    version VARCHAR(50) NOT NULL,
    release_date DATE NOT NULL,
    release_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(app_id, version)
);

CREATE INDEX IF NOT EXISTS idx_app_versions_app_id ON app_versions(app_id);
CREATE INDEX IF NOT EXISTS idx_app_versions_release_date ON app_versions(release_date);
