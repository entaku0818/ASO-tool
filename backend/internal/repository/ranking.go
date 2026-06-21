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

func (r *RankingRepository) ListByKeywordDays(ctx context.Context, keywordID string, days int) ([]*model.RankingHistory, error) {
	since := time.Now().AddDate(0, 0, -days)
	query := `
		SELECT id, keyword_id, rank, recorded_at
		FROM ranking_history
		WHERE keyword_id = $1
		  AND recorded_at >= $2
		ORDER BY recorded_at ASC
	`
	rows, err := r.pool.Query(ctx, query, keywordID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rankings []*model.RankingHistory
	for rows.Next() {
		ranking := &model.RankingHistory{}
		if err := rows.Scan(&ranking.ID, &ranking.KeywordID, &ranking.Rank, &ranking.RecordedAt); err != nil {
			return nil, err
		}
		rankings = append(rankings, ranking)
	}
	return rankings, nil
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

// RisingKeyword holds a keyword whose rank improved compared to N days ago.
type RisingKeyword struct {
	KeywordID    string `json:"keyword_id"`
	Keyword      string `json:"keyword"`
	Country      string `json:"country"`
	CurrentRank  int    `json:"current_rank"`
	PreviousRank int    `json:"previous_rank"`
	Improvement  int    `json:"improvement"`
}

// GetRisingKeywords returns keywords whose rank improved compared to `days` days ago, sorted by improvement desc.
func (r *RankingRepository) GetRisingKeywords(ctx context.Context, appID string, days int) ([]RisingKeyword, error) {
	since := time.Now().AddDate(0, 0, -days)
	rows, err := r.pool.Query(ctx, `
		WITH
		current_rankings AS (
			SELECT DISTINCT ON (k.id)
				k.id AS keyword_id,
				k.keyword,
				k.country,
				rh.rank
			FROM keywords k
			JOIN ranking_history rh ON k.id = rh.keyword_id
			WHERE k.app_id = $1
			ORDER BY k.id, rh.recorded_at DESC
		),
		previous_rankings AS (
			SELECT DISTINCT ON (k.id)
				k.id AS keyword_id,
				rh.rank
			FROM keywords k
			JOIN ranking_history rh ON k.id = rh.keyword_id
			WHERE k.app_id = $1
			  AND rh.recorded_at <= $2
			ORDER BY k.id, rh.recorded_at DESC
		)
		SELECT
			cr.keyword_id,
			cr.keyword,
			cr.country,
			cr.rank AS current_rank,
			pr.rank AS previous_rank,
			(pr.rank - cr.rank) AS improvement
		FROM current_rankings cr
		JOIN previous_rankings pr ON cr.keyword_id = pr.keyword_id
		WHERE pr.rank - cr.rank > 0
		ORDER BY improvement DESC
	`, appID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []RisingKeyword
	for rows.Next() {
		var k RisingKeyword
		if err := rows.Scan(&k.KeywordID, &k.Keyword, &k.Country, &k.CurrentRank, &k.PreviousRank, &k.Improvement); err != nil {
			return nil, err
		}
		keywords = append(keywords, k)
	}
	return keywords, nil
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

// KeywordRankSummary holds the latest rank and change for a keyword.
type KeywordRankSummary struct {
	KeywordID    string `json:"keyword_id"`
	Keyword      string `json:"keyword"`
	Country      string `json:"country"`
	CurrentRank  *int   `json:"current_rank"`
	PreviousRank *int   `json:"previous_rank"`
	Change       *int   `json:"change"`
}

// GetAllKeywordRanks returns current and previous ranks for all keywords of an app.
func (r *RankingRepository) GetAllKeywordRanks(ctx context.Context, appID string) ([]KeywordRankSummary, error) {
	rows, err := r.pool.Query(ctx, `
		WITH
		cur AS (
			SELECT DISTINCT ON (k.id)
				k.id AS keyword_id, k.keyword, k.country, rh.rank AS current_rank
			FROM keywords k
			LEFT JOIN ranking_history rh ON k.id = rh.keyword_id
			WHERE k.app_id = $1
			ORDER BY k.id, rh.recorded_at DESC NULLS LAST
		),
		prev AS (
			SELECT DISTINCT ON (k.id)
				k.id AS keyword_id, rh.rank AS previous_rank
			FROM keywords k
			JOIN ranking_history rh ON k.id = rh.keyword_id
			WHERE k.app_id = $1
			ORDER BY k.id, rh.recorded_at DESC
			OFFSET 0
		),
		prev2 AS (
			SELECT kp.keyword_id, kp.previous_rank
			FROM (
				SELECT k.id AS keyword_id, rh.rank AS previous_rank,
					ROW_NUMBER() OVER (PARTITION BY k.id ORDER BY rh.recorded_at DESC) AS rn
				FROM keywords k
				JOIN ranking_history rh ON k.id = rh.keyword_id
				WHERE k.app_id = $1
			) kp
			WHERE kp.rn = 2
		)
		SELECT c.keyword_id, c.keyword, c.country, c.current_rank, p.previous_rank
		FROM cur c
		LEFT JOIN prev2 p ON c.keyword_id = p.keyword_id
		ORDER BY c.keyword
	`, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []KeywordRankSummary
	for rows.Next() {
		var s KeywordRankSummary
		if err := rows.Scan(&s.KeywordID, &s.Keyword, &s.Country, &s.CurrentRank, &s.PreviousRank); err != nil {
			return nil, err
		}
		if s.CurrentRank != nil && s.PreviousRank != nil {
			change := *s.PreviousRank - *s.CurrentRank // positive = improved
			s.Change = &change
		}
		result = append(result, s)
	}
	return result, nil
}
