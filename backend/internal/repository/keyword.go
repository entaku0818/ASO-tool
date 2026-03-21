package repository

import (
	"context"
	"errors"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type KeywordRepository struct {
	pool *pgxpool.Pool
}

func NewKeywordRepository(pool *pgxpool.Pool) *KeywordRepository {
	return &KeywordRepository{pool: pool}
}

func (r *KeywordRepository) Create(ctx context.Context, req *model.CreateKeywordRequest) (*model.Keyword, error) {
	keyword := &model.Keyword{
		ID:      uuid.New().String(),
		AppID:   req.AppID,
		Keyword: req.Keyword,
		Country: req.Country,
	}

	query := `
		INSERT INTO keywords (id, app_id, keyword, country, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING created_at
	`

	err := r.pool.QueryRow(ctx, query,
		keyword.ID, keyword.AppID, keyword.Keyword, keyword.Country,
	).Scan(&keyword.CreatedAt)

	if err != nil {
		return nil, err
	}

	return keyword, nil
}

func (r *KeywordRepository) Get(ctx context.Context, id string) (*model.Keyword, error) {
	query := `
		SELECT id, app_id, keyword, country, popularity_score, popularity_fetched_at, created_at
		FROM keywords
		WHERE id = $1
	`

	keyword := &model.Keyword{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&keyword.ID, &keyword.AppID, &keyword.Keyword, &keyword.Country,
		&keyword.PopularityScore, &keyword.PopularityFetchedAt, &keyword.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return keyword, nil
}

func (r *KeywordRepository) ListByApp(ctx context.Context, appID string) ([]*model.Keyword, error) {
	query := `
		SELECT id, app_id, keyword, country, popularity_score, popularity_fetched_at, created_at
		FROM keywords
		WHERE app_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []*model.Keyword
	for rows.Next() {
		keyword := &model.Keyword{}
		err := rows.Scan(
			&keyword.ID, &keyword.AppID, &keyword.Keyword, &keyword.Country,
			&keyword.PopularityScore, &keyword.PopularityFetchedAt, &keyword.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		keywords = append(keywords, keyword)
	}

	return keywords, nil
}

func (r *KeywordRepository) UpdatePopularity(ctx context.Context, id string, score int) error {
	query := `
		UPDATE keywords
		SET popularity_score = $2, popularity_fetched_at = NOW()
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, id, score)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}

func (r *KeywordRepository) CountByApp(ctx context.Context, appID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM keywords WHERE app_id = $1`, appID,
	).Scan(&count)
	return count, err
}

func (r *KeywordRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM keywords WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}
