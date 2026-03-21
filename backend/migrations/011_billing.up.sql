ALTER TABLE users
    ADD COLUMN plan                   VARCHAR(20)  NOT NULL DEFAULT 'free',
    ADD COLUMN stripe_customer_id     VARCHAR(255),
    ADD COLUMN stripe_subscription_id VARCHAR(255),
    ADD COLUMN plan_expires_at        TIMESTAMPTZ;

CREATE INDEX idx_users_stripe_customer_id ON users(stripe_customer_id) WHERE stripe_customer_id IS NOT NULL;
