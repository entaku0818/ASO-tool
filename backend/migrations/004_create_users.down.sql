-- Restore original unique constraint on tracked_keywords
ALTER TABLE tracked_keywords DROP CONSTRAINT IF EXISTS tracked_keywords_keyword_country_platform_user_id_key;
ALTER TABLE tracked_keywords ADD CONSTRAINT tracked_keywords_keyword_country_platform_key UNIQUE(keyword, country, platform);

-- Remove user_id from tracked_keywords
DROP INDEX IF EXISTS idx_tracked_keywords_user_id;
ALTER TABLE tracked_keywords DROP COLUMN IF EXISTS user_id;

-- Remove user_id from apps
DROP INDEX IF EXISTS idx_apps_user_id;
ALTER TABLE apps DROP COLUMN IF EXISTS user_id;

-- Drop users table
DROP TABLE IF EXISTS users;
