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

// DailyRankingComparison holds a keyword's ranking for today vs yesterday.
type DailyRankingComparison struct {
	Keyword       string `json:"keyword"`
	Country       string `json:"country"`
	TodayRank     *int   `json:"today_rank"`
	YesterdayRank *int   `json:"yesterday_rank"`
	Change        *int   `json:"change"` // positive = improved (lower rank number is better)
}

// CompareWithYesterday returns keyword rankings for today and yesterday for the given app.
func (r *RankingRepository) CompareWithYesterday(ctx context.Context, appID string, today time.Time) ([]DailyRankingComparison, error) {
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	yesterdayStart := todayStart.AddDate(0, 0, -1)

	rows, err := r.pool.Query(ctx, `
		WITH
		today_rankings AS (
			SELECT DISTINCT ON (k.id)
				k.id AS keyword_id,
				k.keyword,
				k.country,
				rh.rank
			FROM keywords k
			JOIN ranking_history rh ON k.id = rh.keyword_id
			WHERE k.app_id = $1
			  AND rh.recorded_at >= $2
			ORDER BY k.id, rh.recorded_at DESC
		),
		yesterday_rankings AS (
			SELECT DISTINCT ON (k.id)
				k.id AS keyword_id,
				rh.rank
			FROM keywords k
			JOIN ranking_history rh ON k.id = rh.keyword_id
			WHERE k.app_id = $1
			  AND rh.recorded_at >= $3
			  AND rh.recorded_at < $2
			ORDER BY k.id, rh.recorded_at DESC
		),
		all_keywords AS (
			SELECT DISTINCT k.id AS keyword_id, k.keyword, k.country
			FROM keywords k
			WHERE k.app_id = $1
			  AND EXISTS (
				SELECT 1 FROM ranking_history rh
				WHERE rh.keyword_id = k.id
				  AND rh.recorded_at >= $3
			)
		)
		SELECT
			ak.keyword,
			ak.country,
			t.rank AS today_rank,
			y.rank AS yesterday_rank
		FROM all_keywords ak
		LEFT JOIN today_rankings t ON ak.keyword_id = t.keyword_id
		LEFT JOIN yesterday_rankings y ON ak.keyword_id = y.keyword_id
		ORDER BY ak.keyword
	`, appID, todayStart, yesterdayStart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comparisons []DailyRankingComparison
	for rows.Next() {
		var c DailyRankingComparison
		if err := rows.Scan(&c.Keyword, &c.Country, &c.TodayRank, &c.YesterdayRank); err != nil {
			return nil, err
		}
		if c.TodayRank != nil && c.YesterdayRank != nil {
			change := *c.YesterdayRank - *c.TodayRank
			c.Change = &change
		}
		comparisons = append(comparisons, c)
	}
	return comparisons, nil
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
