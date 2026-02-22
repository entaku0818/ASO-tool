-- Store rankings table: iTunes App Store chart rankings (fetched daily via batch)
CREATE TABLE IF NOT EXISTS store_rankings (
    id VARCHAR(36) PRIMARY KEY,
    country VARCHAR(10) NOT NULL,
    ranking_type VARCHAR(50) NOT NULL,
    genre_id VARCHAR(10) NOT NULL DEFAULT '',
    rank INTEGER NOT NULL,
    app_id VARCHAR(50) NOT NULL,
    app_name VARCHAR(500) NOT NULL,
    developer VARCHAR(255) NOT NULL DEFAULT '',
    icon_url TEXT NOT NULL DEFAULT '',
    category VARCHAR(255) NOT NULL DEFAULT '',
    store_url TEXT NOT NULL DEFAULT '',
    price VARCHAR(50) NOT NULL DEFAULT '',
    release_date VARCHAR(100) NOT NULL DEFAULT '',
    fetched_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_store_rankings_lookup ON store_rankings(country, ranking_type, genre_id, fetched_at DESC);
CREATE INDEX IF NOT EXISTS idx_store_rankings_fetched_at ON store_rankings(fetched_at DESC);
