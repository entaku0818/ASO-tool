package repository

import (
	"context"
	"errors"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppRepository struct {
	pool *pgxpool.Pool
}

func NewAppRepository(pool *pgxpool.Pool) *AppRepository {
	return &AppRepository{pool: pool}
}

func (r *AppRepository) Create(ctx context.Context, userID string, req *model.CreateAppRequest) (*model.App, error) {
	app := &model.App{
		ID:       uuid.New().String(),
		Name:     req.Name,
		BundleID: req.BundleID,
		Platform: req.Platform,
		StoreURL: req.StoreURL,
		UserID:   userID,
	}

	query := `
		INSERT INTO apps (id, name, bundle_id, platform, store_url, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		app.ID, app.Name, app.BundleID, app.Platform, app.StoreURL, app.UserID,
	).Scan(&app.CreatedAt, &app.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return app, nil
}

func (r *AppRepository) Get(ctx context.Context, userID, id string) (*model.App, error) {
	query := `
		SELECT id, name, bundle_id, platform, store_url, user_id, created_at, updated_at
		FROM apps
		WHERE id = $1 AND user_id = $2
	`

	app := &model.App{}
	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&app.ID, &app.Name, &app.BundleID, &app.Platform, &app.StoreURL, &app.UserID,
		&app.CreatedAt, &app.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return app, nil
}

// GetByID returns an app without user filtering (for internal use like scraper)
func (r *AppRepository) GetByID(ctx context.Context, id string) (*model.App, error) {
	query := `
		SELECT id, name, bundle_id, platform, store_url, user_id, created_at, updated_at
		FROM apps
		WHERE id = $1
	`

	app := &model.App{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&app.ID, &app.Name, &app.BundleID, &app.Platform, &app.StoreURL, &app.UserID,
		&app.CreatedAt, &app.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (r *AppRepository) List(ctx context.Context, userID string) ([]*model.App, error) {
	query := `
		SELECT id, name, bundle_id, platform, store_url, user_id, created_at, updated_at
		FROM apps
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*model.App
	for rows.Next() {
		app := &model.App{}
		err := rows.Scan(
			&app.ID, &app.Name, &app.BundleID, &app.Platform, &app.StoreURL, &app.UserID,
			&app.CreatedAt, &app.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}

// ListAll returns all apps without user filtering (for internal use like scraper)
func (r *AppRepository) ListAll(ctx context.Context) ([]*model.App, error) {
	query := `
		SELECT id, name, bundle_id, platform, store_url, user_id, created_at, updated_at
		FROM apps
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*model.App
	for rows.Next() {
		app := &model.App{}
		err := rows.Scan(
			&app.ID, &app.Name, &app.BundleID, &app.Platform, &app.StoreURL, &app.UserID,
			&app.CreatedAt, &app.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}

func (r *AppRepository) Update(ctx context.Context, userID, id string, req *model.UpdateAppRequest) (*model.App, error) {
	query := `
		UPDATE apps
		SET name = COALESCE(NULLIF($3, ''), name),
		    store_url = COALESCE(NULLIF($4, ''), store_url),
		    updated_at = NOW()
		WHERE id = $1 AND user_id = $2
		RETURNING id, name, bundle_id, platform, store_url, user_id, created_at, updated_at
	`

	app := &model.App{}
	err := r.pool.QueryRow(ctx, query, id, userID, req.Name, req.StoreURL).Scan(
		&app.ID, &app.Name, &app.BundleID, &app.Platform, &app.StoreURL, &app.UserID,
		&app.CreatedAt, &app.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (r *AppRepository) Delete(ctx context.Context, userID, id string) error {
	query := `DELETE FROM apps WHERE id = $1 AND user_id = $2`

	result, err := r.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}
