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

type AppVersionRepository struct {
	pool *pgxpool.Pool
}

func NewAppVersionRepository(pool *pgxpool.Pool) *AppVersionRepository {
	return &AppVersionRepository{pool: pool}
}

func (r *AppVersionRepository) Create(ctx context.Context, version *model.AppVersion) error {
	if version.ID == "" {
		version.ID = uuid.New().String()
	}

	query := `
		INSERT INTO app_versions (id, app_id, version, release_date, release_notes, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (app_id, version) DO UPDATE SET
			release_date = EXCLUDED.release_date,
			release_notes = EXCLUDED.release_notes
		RETURNING created_at
	`

	err := r.pool.QueryRow(ctx, query,
		version.ID, version.AppID, version.Version, version.ReleaseDate, version.ReleaseNotes,
	).Scan(&version.CreatedAt)

	return err
}

func (r *AppVersionRepository) GetByAppID(ctx context.Context, appID string) ([]*model.AppVersion, error) {
	query := `
		SELECT id, app_id, version, release_date, release_notes, created_at
		FROM app_versions
		WHERE app_id = $1
		ORDER BY release_date DESC
	`

	rows, err := r.pool.Query(ctx, query, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*model.AppVersion
	for rows.Next() {
		v := &model.AppVersion{}
		err := rows.Scan(
			&v.ID, &v.AppID, &v.Version, &v.ReleaseDate, &v.ReleaseNotes, &v.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}

	return versions, nil
}

func (r *AppVersionRepository) GetLatest(ctx context.Context, appID string) (*model.AppVersion, error) {
	query := `
		SELECT id, app_id, version, release_date, release_notes, created_at
		FROM app_versions
		WHERE app_id = $1
		ORDER BY release_date DESC
		LIMIT 1
	`

	v := &model.AppVersion{}
	err := r.pool.QueryRow(ctx, query, appID).Scan(
		&v.ID, &v.AppID, &v.Version, &v.ReleaseDate, &v.ReleaseNotes, &v.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (r *AppVersionRepository) GetByDateRange(ctx context.Context, appID string, start, end time.Time) ([]*model.AppVersion, error) {
	query := `
		SELECT id, app_id, version, release_date, release_notes, created_at
		FROM app_versions
		WHERE app_id = $1 AND release_date >= $2 AND release_date <= $3
		ORDER BY release_date DESC
	`

	rows, err := r.pool.Query(ctx, query, appID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []*model.AppVersion
	for rows.Next() {
		v := &model.AppVersion{}
		err := rows.Scan(
			&v.ID, &v.AppID, &v.Version, &v.ReleaseDate, &v.ReleaseNotes, &v.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}

	return versions, nil
}

func (r *AppVersionRepository) GetVersionOnDate(ctx context.Context, appID string, date time.Time) (*model.AppVersion, error) {
	query := `
		SELECT id, app_id, version, release_date, release_notes, created_at
		FROM app_versions
		WHERE app_id = $1 AND release_date = $2
		LIMIT 1
	`

	v := &model.AppVersion{}
	err := r.pool.QueryRow(ctx, query, appID, date).Scan(
		&v.ID, &v.AppID, &v.Version, &v.ReleaseDate, &v.ReleaseNotes, &v.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return v, nil
}
