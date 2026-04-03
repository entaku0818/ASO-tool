CREATE TABLE screenshot_templates (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    category    VARCHAR(50)  NOT NULL,
    description VARCHAR(500) NOT NULL DEFAULT '',
    device      VARCHAR(50)  NOT NULL DEFAULT 'iphone67',
    style       JSONB        NOT NULL DEFAULT '{}',
    is_pro      BOOLEAN      NOT NULL DEFAULT false,
    sort_order  INTEGER      NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_templates_category ON screenshot_templates(category);
CREATE INDEX idx_templates_sort ON screenshot_templates(sort_order, created_at);

-- Initial data: 9 templates matching the design doc lineup
INSERT INTO screenshot_templates (name, category, description, device, style, is_pro, sort_order) VALUES
  ('ダークファンタジー RPG',       'game',      '重厚な世界観・戦闘系アプリ向け',     'iphone67', '{"bg_gradient_from":"#0f0c29","bg_gradient_to":"#302b63","bg_gradient_dir":"tb","text_color":"#E2E8F0","image_align":"bottom"}', false, 10),
  ('ポップ＆カジュアル',           'game',      '明るい色調・パズル・カジュアル向け', 'iphone67', '{"bg_gradient_from":"#f093fb","bg_gradient_to":"#f5576c","bg_gradient_dir":"tb","text_color":"#FFFFFF","image_align":"bottom"}', false, 20),
  ('スポーツ＆アクション',         'game',      '躍動感・スポーツ・アクション向け',   'iphone67', '{"bg_gradient_from":"#4facfe","bg_gradient_to":"#00f2fe","bg_gradient_dir":"tb","text_color":"#FFFFFF","image_align":"bottom"}', true,  30),
  ('クリーン・プロダクティビティ', 'business',  'シンプル・ミニマル・業務効率向け',   'iphone67', '{"bg_color":"#F8FAFC","text_color":"#1E293B","image_align":"center"}',                                                          false, 10),
  ('エンタープライズ SaaS',        'business',  '信頼感・B2B・ビジネスツール向け',    'iphone67', '{"bg_gradient_from":"#1D4ED8","bg_gradient_to":"#1E3A8A","bg_gradient_dir":"tb","text_color":"#FFFFFF","image_align":"bottom"}', true,  20),
  ('スタートアップ モダン',        'business',  'グラデーション・Fintech・新興向け',  'iphone67', '{"bg_gradient_from":"#667eea","bg_gradient_to":"#764ba2","bg_gradient_dir":"tlbr","text_color":"#FFFFFF","image_align":"bottom"}', true, 30),
  ('キッズ＆ラーニング',           'education', 'カラフル・子供向け学習アプリ',       'iphone67', '{"bg_gradient_from":"#f7971e","bg_gradient_to":"#ffd200","bg_gradient_dir":"tb","text_color":"#1E293B","image_align":"bottom"}', false, 10),
  ('語学・スキルアップ',           'education', '落ち着いた・大人の学習向け',         'iphone67', '{"bg_gradient_from":"#134e5e","bg_gradient_to":"#71b280","bg_gradient_dir":"tb","text_color":"#FFFFFF","image_align":"bottom"}', false, 20),
  ('資格・受験対策',               'education', 'シリアス・試験・問題集向け',         'iphone67', '{"bg_gradient_from":"#1a1a2e","bg_gradient_to":"#16213e","bg_gradient_dir":"tb","text_color":"#E2E8F0","image_align":"bottom"}', true,  30);
