-- Competitors テーブル: アプリごとの競合アプリ
CREATE TABLE IF NOT EXISTS competitors (
    id VARCHAR(36) PRIMARY KEY,
    app_id VARCHAR(36) NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    competitor_bundle_id VARCHAR(255) NOT NULL,
    competitor_name VARCHAR(255) NOT NULL,
    platform VARCHAR(20) NOT NULL CHECK (platform IN ('ios', 'android')),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(app_id, competitor_bundle_id)
);

CREATE INDEX idx_competitors_app_id ON competitors(app_id);

-- Competitor Rankings テーブル: 競合アプリのキーワード順位履歴
CREATE TABLE IF NOT EXISTS competitor_rankings (
    id VARCHAR(36) PRIMARY KEY,
    competitor_id VARCHAR(36) NOT NULL REFERENCES competitors(id) ON DELETE CASCADE,
    keyword_id VARCHAR(36) NOT NULL REFERENCES keywords(id) ON DELETE CASCADE,
    rank INTEGER,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_competitor_rankings_competitor_id ON competitor_rankings(competitor_id);
CREATE INDEX idx_competitor_rankings_keyword_id ON competitor_rankings(keyword_id);
CREATE INDEX idx_competitor_rankings_recorded_at ON competitor_rankings(recorded_at);
