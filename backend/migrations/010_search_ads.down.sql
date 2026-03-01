ALTER TABLE keywords
    DROP COLUMN IF EXISTS popularity_score,
    DROP COLUMN IF EXISTS popularity_fetched_at;

DROP TABLE IF EXISTS search_ads_credentials;
