package repository

import (
	"context"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScreenshotRepository struct {
	pool *pgxpool.Pool
}

func NewScreenshotRepository(pool *pgxpool.Pool) *ScreenshotRepository {
	return &ScreenshotRepository{pool: pool}
}

func (r *ScreenshotRepository) Create(ctx context.Context, req *model.CreateScreenshotRequest) (*model.Screenshot, error) {
	id := uuid.New().String()
	now := time.Now()

	deviceType := req.DeviceType
	if deviceType == "" {
		deviceType = model.DeviceTypeIPhone
	}

	locale := req.Locale
	if locale == "" {
		locale = "ja-JP"
	}

	query := `
		INSERT INTO screenshots (id, app_id, url, display_order, device_type, locale, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, app_id, url, display_order, device_type, locale, description, created_at, updated_at
	`

	var s model.Screenshot
	err := r.pool.QueryRow(ctx, query,
		id, req.AppID, req.URL, req.DisplayOrder, deviceType, locale, req.Description, now, now,
	).Scan(&s.ID, &s.AppID, &s.URL, &s.DisplayOrder, &s.DeviceType, &s.Locale, &s.Description, &s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *ScreenshotRepository) Get(ctx context.Context, id string) (*model.Screenshot, error) {
	query := `
		SELECT id, app_id, url, display_order, device_type, locale, description, created_at, updated_at
		FROM screenshots
		WHERE id = $1
	`

	var s model.Screenshot
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.AppID, &s.URL, &s.DisplayOrder, &s.DeviceType, &s.Locale, &s.Description, &s.CreatedAt, &s.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *ScreenshotRepository) ListByApp(ctx context.Context, appID string) ([]*model.Screenshot, error) {
	query := `
		SELECT id, app_id, url, display_order, device_type, locale, description, created_at, updated_at
		FROM screenshots
		WHERE app_id = $1
		ORDER BY display_order ASC, created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var screenshots []*model.Screenshot
	for rows.Next() {
		var s model.Screenshot
		err := rows.Scan(&s.ID, &s.AppID, &s.URL, &s.DisplayOrder, &s.DeviceType, &s.Locale, &s.Description, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		screenshots = append(screenshots, &s)
	}

	return screenshots, nil
}

func (r *ScreenshotRepository) ListByAppAndDevice(ctx context.Context, appID string, deviceType model.DeviceType) ([]*model.Screenshot, error) {
	query := `
		SELECT id, app_id, url, display_order, device_type, locale, description, created_at, updated_at
		FROM screenshots
		WHERE app_id = $1 AND device_type = $2
		ORDER BY display_order ASC, created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, appID, deviceType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var screenshots []*model.Screenshot
	for rows.Next() {
		var s model.Screenshot
		err := rows.Scan(&s.ID, &s.AppID, &s.URL, &s.DisplayOrder, &s.DeviceType, &s.Locale, &s.Description, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		screenshots = append(screenshots, &s)
	}

	return screenshots, nil
}

func (r *ScreenshotRepository) Update(ctx context.Context, id string, req *model.UpdateScreenshotRequest) (*model.Screenshot, error) {
	now := time.Now()

	query := `
		UPDATE screenshots
		SET url = COALESCE($2, url),
			display_order = COALESCE($3, display_order),
			device_type = COALESCE($4, device_type),
			locale = COALESCE($5, locale),
			description = COALESCE($6, description),
			updated_at = $7
		WHERE id = $1
		RETURNING id, app_id, url, display_order, device_type, locale, description, created_at, updated_at
	`

	var s model.Screenshot
	err := r.pool.QueryRow(ctx, query,
		id, req.URL, req.DisplayOrder, req.DeviceType, req.Locale, req.Description, now,
	).Scan(&s.ID, &s.AppID, &s.URL, &s.DisplayOrder, &s.DeviceType, &s.Locale, &s.Description, &s.CreatedAt, &s.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *ScreenshotRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM screenshots WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *ScreenshotRepository) Reorder(ctx context.Context, appID string, screenshotIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, id := range screenshotIDs {
		query := `UPDATE screenshots SET display_order = $1, updated_at = $2 WHERE id = $3 AND app_id = $4`
		result, err := tx.Exec(ctx, query, i, time.Now(), id, appID)
		if err != nil {
			return err
		}
		if result.RowsAffected() == 0 {
			return model.ErrNotFound
		}
	}

	return tx.Commit(ctx)
}

func (r *ScreenshotRepository) GetMaxDisplayOrder(ctx context.Context, appID string) (int, error) {
	query := `SELECT COALESCE(MAX(display_order), -1) FROM screenshots WHERE app_id = $1`
	var maxOrder int
	err := r.pool.QueryRow(ctx, query, appID).Scan(&maxOrder)
	if err != nil {
		return 0, err
	}
	return maxOrder, nil
}
