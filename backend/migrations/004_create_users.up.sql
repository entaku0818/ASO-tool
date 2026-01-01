-- Create users table
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- Add user_id to apps table
ALTER TABLE apps ADD COLUMN user_id VARCHAR(36) REFERENCES users(id);
CREATE INDEX idx_apps_user_id ON apps(user_id);

-- Add user_id to tracked_keywords table
ALTER TABLE tracked_keywords ADD COLUMN user_id VARCHAR(36) REFERENCES users(id);
CREATE INDEX idx_tracked_keywords_user_id ON tracked_keywords(user_id);

-- Update unique constraint on tracked_keywords to include user_id
ALTER TABLE tracked_keywords DROP CONSTRAINT tracked_keywords_keyword_country_platform_key;
ALTER TABLE tracked_keywords ADD CONSTRAINT tracked_keywords_keyword_country_platform_user_id_key UNIQUE(keyword, country, platform, user_id);
