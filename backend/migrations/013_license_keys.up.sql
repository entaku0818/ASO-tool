CREATE TABLE license_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(24) UNIQUE NOT NULL,
    email VARCHAR(255) NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT false,
    activated_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    stripe_session_id VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_license_keys_key ON license_keys(key);
CREATE INDEX idx_license_keys_email ON license_keys(email);
CREATE INDEX idx_license_keys_user_id ON license_keys(user_id);
