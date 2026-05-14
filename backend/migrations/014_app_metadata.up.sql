CREATE TABLE app_metadata_versions (
    id            VARCHAR(36) PRIMARY KEY,
    app_id        VARCHAR(36) NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    locale        VARCHAR(10) NOT NULL DEFAULT 'ja',
    version_tag   VARCHAR(50) NOT NULL DEFAULT 'draft',
    title         VARCHAR(30),
    subtitle      VARCHAR(30),
    description   VARCHAR(4000),
    keywords      VARCHAR(100),
    promotional_text VARCHAR(170),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_app_metadata_versions_app_id ON app_metadata_versions(app_id);
