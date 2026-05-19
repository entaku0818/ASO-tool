-- ASC Analytics search keyword report requests and results
CREATE TABLE IF NOT EXISTS asc_report_requests (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id      UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    request_id  TEXT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'pending', -- pending, ready, failed
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(app_id, request_id)
);
CREATE INDEX IF NOT EXISTS idx_asc_report_requests_app_id ON asc_report_requests(app_id);

CREATE TABLE IF NOT EXISTS asc_search_keywords (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id      UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    date        DATE NOT NULL,
    keyword     TEXT NOT NULL,
    impressions INTEGER NOT NULL DEFAULT 0,
    page_views  INTEGER NOT NULL DEFAULT 0,
    installs    INTEGER NOT NULL DEFAULT 0,
    fetched_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(app_id, date, keyword)
);
CREATE INDEX IF NOT EXISTS idx_asc_search_keywords_app_id ON asc_search_keywords(app_id);
CREATE INDEX IF NOT EXISTS idx_asc_search_keywords_app_date ON asc_search_keywords(app_id, date DESC);
CREATE INDEX IF NOT EXISTS idx_asc_search_keywords_installs ON asc_search_keywords(app_id, installs DESC);
