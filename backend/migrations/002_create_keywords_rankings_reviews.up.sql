-- Keywords table: tracks keywords for each app
CREATE TABLE IF NOT EXISTS keywords (
    id VARCHAR(36) PRIMARY KEY,
    app_id VARCHAR(36) NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    keyword VARCHAR(255) NOT NULL,
    country VARCHAR(2) NOT NULL DEFAULT 'JP',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(app_id, keyword, country)
);

CREATE INDEX idx_keywords_app_id ON keywords(app_id);

-- Ranking history table: stores daily ranking for each keyword
CREATE TABLE IF NOT EXISTS ranking_history (
    id VARCHAR(36) PRIMARY KEY,
    keyword_id VARCHAR(36) NOT NULL REFERENCES keywords(id) ON DELETE CASCADE,
    rank INTEGER,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_ranking_history_keyword_id ON ranking_history(keyword_id);
CREATE INDEX idx_ranking_history_recorded_at ON ranking_history(recorded_at);

-- Reviews table: stores app reviews
CREATE TABLE IF NOT EXISTS reviews (
    id VARCHAR(36) PRIMARY KEY,
    app_id VARCHAR(36) NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    review_id VARCHAR(255) NOT NULL,
    author VARCHAR(255),
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title TEXT,
    content TEXT,
    version VARCHAR(50),
    reviewed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(app_id, review_id)
);

CREATE INDEX idx_reviews_app_id ON reviews(app_id);
CREATE INDEX idx_reviews_rating ON reviews(rating);
CREATE INDEX idx_reviews_reviewed_at ON reviews(reviewed_at);
