-- Screenshots テーブル: アプリのスクリーンショット管理
CREATE TABLE IF NOT EXISTS screenshots (
    id VARCHAR(36) PRIMARY KEY,
    app_id VARCHAR(36) NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    url VARCHAR(1024) NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    device_type VARCHAR(50) NOT NULL DEFAULT 'iphone',
    locale VARCHAR(10) NOT NULL DEFAULT 'ja-JP',
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_screenshots_app_id ON screenshots(app_id);
CREATE INDEX idx_screenshots_display_order ON screenshots(app_id, display_order);
