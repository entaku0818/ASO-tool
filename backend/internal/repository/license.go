package repository

import (
	"context"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LicenseRepository struct {
	pool *pgxpool.Pool
}

func NewLicenseRepository(pool *pgxpool.Pool) *LicenseRepository {
	return &LicenseRepository{pool: pool}
}

func (r *LicenseRepository) Create(ctx context.Context, key, email string, stripeSessionID *string) (*model.LicenseKey, error) {
	query := `
		INSERT INTO license_keys (key, email, stripe_session_id)
		VALUES ($1, $2, $3)
		RETURNING id, key, email, user_id, is_active, activated_at, expires_at, stripe_session_id, created_at
	`
	return r.scan(r.pool.QueryRow(ctx, query, key, email, stripeSessionID))
}

func (r *LicenseRepository) GetByKey(ctx context.Context, key string) (*model.LicenseKey, error) {
	query := `
		SELECT id, key, email, user_id, is_active, activated_at, expires_at, stripe_session_id, created_at
		FROM license_keys WHERE key = $1
	`
	return r.scan(r.pool.QueryRow(ctx, query, key))
}

func (r *LicenseRepository) GetByUserID(ctx context.Context, userID string) (*model.LicenseKey, error) {
	query := `
		SELECT id, key, email, user_id, is_active, activated_at, expires_at, stripe_session_id, created_at
		FROM license_keys WHERE user_id = $1 AND is_active = true
		ORDER BY activated_at DESC LIMIT 1
	`
	return r.scan(r.pool.QueryRow(ctx, query, userID))
}

func (r *LicenseRepository) Activate(ctx context.Context, keyID, userID string) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx,
		`UPDATE license_keys SET is_active = true, user_id = $1, activated_at = $2 WHERE id = $3`,
		userID, now, keyID,
	)
	return err
}

func (r *LicenseRepository) List(ctx context.Context) ([]*model.LicenseKey, error) {
	query := `
		SELECT id, key, email, user_id, is_active, activated_at, expires_at, stripe_session_id, created_at
		FROM license_keys ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []*model.LicenseKey
	for rows.Next() {
		lk := &model.LicenseKey{}
		err := rows.Scan(
			&lk.ID, &lk.Key, &lk.Email, &lk.UserID,
			&lk.IsActive, &lk.ActivatedAt, &lk.ExpiresAt,
			&lk.StripeSessionID, &lk.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		keys = append(keys, lk)
	}
	return keys, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func (r *LicenseRepository) scan(row scannable) (*model.LicenseKey, error) {
	lk := &model.LicenseKey{}
	err := row.Scan(
		&lk.ID, &lk.Key, &lk.Email, &lk.UserID,
		&lk.IsActive, &lk.ActivatedAt, &lk.ExpiresAt,
		&lk.StripeSessionID, &lk.CreatedAt,
	)
	if err != nil {
		return nil, model.ErrNotFound
	}
	return lk, nil
}
