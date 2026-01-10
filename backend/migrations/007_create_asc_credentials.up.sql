-- App Store Connect API credentials (encrypted storage)
CREATE TABLE IF NOT EXISTS asc_credentials (
    id VARCHAR(36) PRIMARY KEY,
    app_id VARCHAR(36) NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    issuer_id VARCHAR(255) NOT NULL,
    key_id VARCHAR(20) NOT NULL,
    private_key_encrypted BYTEA NOT NULL,
    encryption_nonce BYTEA NOT NULL,
    is_valid BOOLEAN DEFAULT true,
    last_validated_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(app_id)
);

CREATE INDEX IF NOT EXISTS idx_asc_credentials_app_id ON asc_credentials(app_id);
