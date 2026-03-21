DROP INDEX IF EXISTS idx_users_stripe_customer_id;

ALTER TABLE users
    DROP COLUMN IF EXISTS plan_expires_at,
    DROP COLUMN IF EXISTS stripe_subscription_id,
    DROP COLUMN IF EXISTS stripe_customer_id,
    DROP COLUMN IF EXISTS plan;
