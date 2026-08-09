package repository

import (
	"context"
	"errors"
	"time"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash, name string, isAdmin bool) (*model.User, error) {
	user := &model.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
		IsAdmin:      isAdmin,
		Plan:         "free",
	}

	query := `
		INSERT INTO users (id, email, password_hash, name, is_admin, plan, created_at)
		VALUES ($1, $2, $3, $4, $5, 'free', NOW())
		RETURNING created_at
	`

	err := r.pool.QueryRow(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Name, user.IsAdmin,
	).Scan(&user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, name, is_admin, plan,
		       COALESCE(stripe_customer_id, ''), COALESCE(stripe_subscription_id, ''),
		       plan_expires_at, created_at
		FROM users WHERE email = $1
	`
	return r.scanUser(r.pool.QueryRow(ctx, query, email))
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, name, is_admin, plan,
		       COALESCE(stripe_customer_id, ''), COALESCE(stripe_subscription_id, ''),
		       plan_expires_at, created_at
		FROM users WHERE id = $1
	`
	return r.scanUser(r.pool.QueryRow(ctx, query, id))
}

func (r *UserRepository) GetByStripeCustomerID(ctx context.Context, customerID string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, name, is_admin, plan,
		       COALESCE(stripe_customer_id, ''), COALESCE(stripe_subscription_id, ''),
		       plan_expires_at, created_at
		FROM users WHERE stripe_customer_id = $1
	`
	return r.scanUser(r.pool.QueryRow(ctx, query, customerID))
}

func (r *UserRepository) List(ctx context.Context) ([]*model.User, error) {
	query := `
		SELECT id, email, password_hash, name, is_admin, plan,
		       COALESCE(stripe_customer_id, ''), COALESCE(stripe_subscription_id, ''),
		       plan_expires_at, created_at
		FROM users ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.IsAdmin,
			&user.Plan, &user.StripeCustomerID, &user.StripeSubscriptionID,
			&user.PlanExpiresAt, &user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdatePlan sets a user's plan and associated Stripe identifiers.
func (r *UserRepository) UpdatePlan(ctx context.Context, userID, plan, stripeCustomerID, stripeSubscriptionID string, expiresAt *time.Time) error {
	query := `
		UPDATE users
		SET plan = $2,
		    stripe_customer_id = NULLIF($3, ''),
		    stripe_subscription_id = NULLIF($4, ''),
		    plan_expires_at = $5
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, userID, plan, stripeCustomerID, stripeSubscriptionID, expiresAt)
	return err
}

// SetStripeCustomerID stores the Stripe customer ID for a user.
func (r *UserRepository) SetStripeCustomerID(ctx context.Context, userID, customerID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET stripe_customer_id = $2 WHERE id = $1`,
		userID, customerID,
	)
	return err
}

func (r *UserRepository) scanUser(row pgx.Row) (*model.User, error) {
	user := &model.User{}
	err := row.Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.IsAdmin,
		&user.Plan, &user.StripeCustomerID, &user.StripeSubscriptionID,
		&user.PlanExpiresAt, &user.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
