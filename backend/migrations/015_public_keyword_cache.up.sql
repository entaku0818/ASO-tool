CREATE TABLE IF NOT EXISTS public_keyword_cache (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    keyword     TEXT NOT NULL,
    country     TEXT NOT NULL,
    genre       TEXT NOT NULL DEFAULT '',
    popularity  INTEGER NOT NULL DEFAULT 0,
    fetched_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(keyword, country, genre)
);
CREATE INDEX IF NOT EXISTS idx_public_keyword_cache_country_genre ON public_keyword_cache(country, genre);
CREATE INDEX IF NOT EXISTS idx_public_keyword_cache_popularity ON public_keyword_cache(popularity DESC);
