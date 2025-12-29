package repository

import (
	"context"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RankingRepository struct {
	pool *pgxpool.Pool
}

func NewRankingRepository(pool *pgxpool.Pool) *RankingRepository {
	return &RankingRepository{pool: pool}
}

func (r *RankingRepository) Create(ctx context.Context, req *model.CreateRankingRequest) (*model.RankingHistory, error) {
	ranking := &model.RankingHistory{
		ID:        uuid.New().String(),
		KeywordID: req.KeywordID,
		Rank:      req.Rank,
	}

	query := `
		INSERT INTO ranking_history (id, keyword_id, rank, recorded_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING recorded_at
	`

	err := r.pool.QueryRow(ctx, query,
		ranking.ID, ranking.KeywordID, ranking.Rank,
	).Scan(&ranking.RecordedAt)

	if err != nil {
		return nil, err
	}

	return ranking, nil
}

func (r *RankingRepository) ListByKeyword(ctx context.Context, keywordID string, limit int) ([]*model.RankingHistory, error) {
	query := `
		SELECT id, keyword_id, rank, recorded_at
		FROM ranking_history
		WHERE keyword_id = $1
		ORDER BY recorded_at DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, keywordID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rankings []*model.RankingHistory
	for rows.Next() {
		ranking := &model.RankingHistory{}
		err := rows.Scan(
			&ranking.ID, &ranking.KeywordID, &ranking.Rank, &ranking.RecordedAt,
		)
		if err != nil {
			return nil, err
		}
		rankings = append(rankings, ranking)
	}

	return rankings, nil
}

func (r *RankingRepository) ListByAppWithKeyword(ctx context.Context, appID string, from, to time.Time) ([]*model.RankingWithKeyword, error) {
	query := `
		SELECT rh.id, rh.keyword_id, rh.rank, rh.recorded_at, k.keyword, k.country
		FROM ranking_history rh
		JOIN keywords k ON rh.keyword_id = k.id
		WHERE k.app_id = $1
		AND rh.recorded_at >= $2 AND rh.recorded_at <= $3
		ORDER BY rh.recorded_at DESC
	`

	rows, err := r.pool.Query(ctx, query, appID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rankings []*model.RankingWithKeyword
	for rows.Next() {
		ranking := &model.RankingWithKeyword{}
		err := rows.Scan(
			&ranking.ID, &ranking.KeywordID, &ranking.Rank, &ranking.RecordedAt,
			&ranking.Keyword, &ranking.Country,
		)
		if err != nil {
			return nil, err
		}
		rankings = append(rankings, ranking)
	}

	return rankings, nil
}

func (r *RankingRepository) GetLatestByKeyword(ctx context.Context, keywordID string) (*model.RankingHistory, error) {
	query := `
		SELECT id, keyword_id, rank, recorded_at
		FROM ranking_history
		WHERE keyword_id = $1
		ORDER BY recorded_at DESC
		LIMIT 1
	`

	ranking := &model.RankingHistory{}
	err := r.pool.QueryRow(ctx, query, keywordID).Scan(
		&ranking.ID, &ranking.KeywordID, &ranking.Rank, &ranking.RecordedAt,
	)

	if err != nil {
		return nil, err
	}

	return ranking, nil
}
