-- Search Ads API credentials (encrypted storage)
CREATE TABLE IF NOT EXISTS search_ads_credentials (
    id VARCHAR(36) PRIMARY KEY,
    app_id VARCHAR(36) NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    client_id VARCHAR(255) NOT NULL,
    team_id VARCHAR(255) NOT NULL,
    key_id VARCHAR(20) NOT NULL,
    private_key_encrypted BYTEA NOT NULL,
    encryption_nonce BYTEA NOT NULL,
    org_id VARCHAR(255),
    adam_id BIGINT,
    is_valid BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(app_id)
);

CREATE INDEX IF NOT EXISTS idx_search_ads_credentials_app_id ON search_ads_credentials(app_id);

-- Add popularity score columns to keywords table
ALTER TABLE keywords
    ADD COLUMN IF NOT EXISTS popularity_score INTEGER,
    ADD COLUMN IF NOT EXISTS popularity_fetched_at TIMESTAMP WITH TIME ZONE;
