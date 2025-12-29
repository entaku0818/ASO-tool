CREATE TABLE IF NOT EXISTS apps (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    bundle_id VARCHAR(255) NOT NULL,
    platform VARCHAR(20) NOT NULL CHECK (platform IN ('ios', 'android')),
    store_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_apps_bundle_id ON apps(bundle_id);
CREATE INDEX idx_apps_platform ON apps(platform);
