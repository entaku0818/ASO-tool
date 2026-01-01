package repository

import (
	"context"
	"errors"

	"github.com/entaku0818/aso-tool/backend/internal/model"
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
	}

	query := `
		INSERT INTO users (id, email, password_hash, name, is_admin, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
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
		SELECT id, email, password_hash, name, is_admin, created_at
		FROM users
		WHERE email = $1
	`

	user := &model.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.IsAdmin, &user.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, name, is_admin, created_at
		FROM users
		WHERE id = $1
	`

	user := &model.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.IsAdmin, &user.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]*model.User, error) {
	query := `
		SELECT id, email, password_hash, name, is_admin, created_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.IsAdmin, &user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
