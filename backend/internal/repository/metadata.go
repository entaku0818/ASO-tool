package repository

import (
	"context"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MetadataRepository struct {
	pool *pgxpool.Pool
}

func NewMetadataRepository(pool *pgxpool.Pool) *MetadataRepository {
	return &MetadataRepository{pool: pool}
}

func (r *MetadataRepository) scanRow(row pgx.Row) (*model.AppMetadataVersion, error) {
	var m model.AppMetadataVersion
	err := row.Scan(
		&m.ID, &m.AppID, &m.Locale, &m.VersionTag,
		&m.Title, &m.Subtitle, &m.Description, &m.Keywords, &m.PromotionalText,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

const selectCols = `id, app_id, locale, version_tag, title, subtitle, description, keywords, promotional_text, created_at, updated_at`

func (r *MetadataRepository) List(ctx context.Context, appID string) ([]*model.AppMetadataVersion, error) {
	query := `SELECT ` + selectCols + ` FROM app_metadata_versions WHERE app_id = $1 ORDER BY updated_at DESC`
	rows, err := r.pool.Query(ctx, query, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.AppMetadataVersion
	for rows.Next() {
		var m model.AppMetadataVersion
		err := rows.Scan(
			&m.ID, &m.AppID, &m.Locale, &m.VersionTag,
			&m.Title, &m.Subtitle, &m.Description, &m.Keywords, &m.PromotionalText,
			&m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &m)
	}
	return results, nil
}

func (r *MetadataRepository) Get(ctx context.Context, id string) (*model.AppMetadataVersion, error) {
	query := `SELECT ` + selectCols + ` FROM app_metadata_versions WHERE id = $1`
	m, err := r.scanRow(r.pool.QueryRow(ctx, query, id))
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return m, err
}

func (r *MetadataRepository) Upsert(ctx context.Context, req *model.UpsertMetadataRequest) (*model.AppMetadataVersion, error) {
	now := time.Now()
	// try update first
	updateQ := `
		UPDATE app_metadata_versions
		SET title=$1, subtitle=$2, description=$3, keywords=$4, promotional_text=$5, updated_at=$6
		WHERE app_id=$7 AND locale=$8 AND version_tag=$9
		RETURNING ` + selectCols
	m, err := r.scanRow(r.pool.QueryRow(ctx, updateQ,
		req.Title, req.Subtitle, req.Description, req.Keywords, req.PromotionalText, now,
		req.AppID, req.Locale, req.VersionTag,
	))
	if err == nil {
		return m, nil
	}
	if err != pgx.ErrNoRows {
		return nil, err
	}

	// insert
	insertQ := `
		INSERT INTO app_metadata_versions (id, app_id, locale, version_tag, title, subtitle, description, keywords, promotional_text, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$10)
		RETURNING ` + selectCols
	return r.scanRow(r.pool.QueryRow(ctx, insertQ,
		uuid.New().String(), req.AppID, req.Locale, req.VersionTag,
		req.Title, req.Subtitle, req.Description, req.Keywords, req.PromotionalText, now,
	))
}

func (r *MetadataRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM app_metadata_versions WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
