package repository

import (
	"context"
	"errors"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepository struct {
	pool *pgxpool.Pool
}

func NewReviewRepository(pool *pgxpool.Pool) *ReviewRepository {
	return &ReviewRepository{pool: pool}
}

func (r *ReviewRepository) Create(ctx context.Context, req *model.CreateReviewRequest) (*model.Review, error) {
	review := &model.Review{
		ID:         uuid.New().String(),
		AppID:      req.AppID,
		ReviewID:   req.ReviewID,
		Author:     req.Author,
		Rating:     req.Rating,
		Title:      req.Title,
		Content:    req.Content,
		Version:    req.Version,
		ReviewedAt: req.ReviewedAt,
	}

	query := `
		INSERT INTO reviews (id, app_id, review_id, author, rating, title, content, version, reviewed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		ON CONFLICT (app_id, review_id) DO UPDATE SET
			rating = EXCLUDED.rating,
			title = EXCLUDED.title,
			content = EXCLUDED.content
		RETURNING created_at
	`

	err := r.pool.QueryRow(ctx, query,
		review.ID, review.AppID, review.ReviewID, review.Author, review.Rating,
		review.Title, review.Content, review.Version, review.ReviewedAt,
	).Scan(&review.CreatedAt)

	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *ReviewRepository) Get(ctx context.Context, id string) (*model.Review, error) {
	query := `
		SELECT id, app_id, review_id, author, rating, title, content, version, reviewed_at, created_at
		FROM reviews
		WHERE id = $1
	`

	review := &model.Review{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&review.ID, &review.AppID, &review.ReviewID, &review.Author, &review.Rating,
		&review.Title, &review.Content, &review.Version, &review.ReviewedAt, &review.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *ReviewRepository) ListByApp(ctx context.Context, appID string, limit, offset int) ([]*model.Review, error) {
	query := `
		SELECT id, app_id, review_id, author, rating, title, content, version, reviewed_at, created_at
		FROM reviews
		WHERE app_id = $1
		ORDER BY reviewed_at DESC NULLS LAST, created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, appID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*model.Review
	for rows.Next() {
		review := &model.Review{}
		err := rows.Scan(
			&review.ID, &review.AppID, &review.ReviewID, &review.Author, &review.Rating,
			&review.Title, &review.Content, &review.Version, &review.ReviewedAt, &review.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r *ReviewRepository) GetStats(ctx context.Context, appID string) (*model.ReviewStats, error) {
	query := `
		SELECT
			COUNT(*) as total_count,
			COALESCE(AVG(rating), 0) as average_rating,
			COUNT(*) FILTER (WHERE rating = 1) as rating_1_count,
			COUNT(*) FILTER (WHERE rating = 2) as rating_2_count,
			COUNT(*) FILTER (WHERE rating = 3) as rating_3_count,
			COUNT(*) FILTER (WHERE rating = 4) as rating_4_count,
			COUNT(*) FILTER (WHERE rating = 5) as rating_5_count
		FROM reviews
		WHERE app_id = $1
	`

	stats := &model.ReviewStats{AppID: appID}
	err := r.pool.QueryRow(ctx, query, appID).Scan(
		&stats.TotalCount, &stats.AverageRating,
		&stats.Rating1Count, &stats.Rating2Count, &stats.Rating3Count,
		&stats.Rating4Count, &stats.Rating5Count,
	)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *ReviewRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM reviews WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}
