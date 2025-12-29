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

func (r *AppRepository) Create(ctx context.Context, req *model.CreateAppRequest) (*model.App, error) {
	app := &model.App{
		ID:       uuid.New().String(),
		Name:     req.Name,
		BundleID: req.BundleID,
		Platform: req.Platform,
		StoreURL: req.StoreURL,
	}

	query := `
		INSERT INTO apps (id, name, bundle_id, platform, store_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		app.ID, app.Name, app.BundleID, app.Platform, app.StoreURL,
	).Scan(&app.CreatedAt, &app.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return app, nil
}

func (r *AppRepository) Get(ctx context.Context, id string) (*model.App, error) {
	query := `
		SELECT id, name, bundle_id, platform, store_url, created_at, updated_at
		FROM apps
		WHERE id = $1
	`

	app := &model.App{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&app.ID, &app.Name, &app.BundleID, &app.Platform, &app.StoreURL,
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

func (r *AppRepository) List(ctx context.Context) ([]*model.App, error) {
	query := `
		SELECT id, name, bundle_id, platform, store_url, created_at, updated_at
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
			&app.ID, &app.Name, &app.BundleID, &app.Platform, &app.StoreURL,
			&app.CreatedAt, &app.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}

func (r *AppRepository) Update(ctx context.Context, id string, req *model.UpdateAppRequest) (*model.App, error) {
	query := `
		UPDATE apps
		SET name = COALESCE(NULLIF($2, ''), name),
		    store_url = COALESCE(NULLIF($3, ''), store_url),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, bundle_id, platform, store_url, created_at, updated_at
	`

	app := &model.App{}
	err := r.pool.QueryRow(ctx, query, id, req.Name, req.StoreURL).Scan(
		&app.ID, &app.Name, &app.BundleID, &app.Platform, &app.StoreURL,
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

func (r *AppRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM apps WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}
