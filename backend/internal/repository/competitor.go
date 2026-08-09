package repository

import (
	"context"
	"time"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CompetitorRepository struct {
	pool *pgxpool.Pool
}

func NewCompetitorRepository(pool *pgxpool.Pool) *CompetitorRepository {
	return &CompetitorRepository{pool: pool}
}

func (r *CompetitorRepository) Create(ctx context.Context, req *model.CreateCompetitorRequest) (*model.Competitor, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO competitors (id, app_id, competitor_bundle_id, competitor_name, platform, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, app_id, competitor_bundle_id, competitor_name, platform, notes, created_at
	`

	var c model.Competitor
	err := r.pool.QueryRow(ctx, query,
		id, req.AppID, req.CompetitorBundleID, req.CompetitorName, req.Platform, req.Notes, now,
	).Scan(&c.ID, &c.AppID, &c.CompetitorBundleID, &c.CompetitorName, &c.Platform, &c.Notes, &c.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CompetitorRepository) Get(ctx context.Context, id string) (*model.Competitor, error) {
	query := `
		SELECT id, app_id, competitor_bundle_id, competitor_name, platform, notes, created_at
		FROM competitors
		WHERE id = $1
	`

	var c model.Competitor
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.AppID, &c.CompetitorBundleID, &c.CompetitorName, &c.Platform, &c.Notes, &c.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CompetitorRepository) ListByApp(ctx context.Context, appID string) ([]*model.Competitor, error) {
	query := `
		SELECT id, app_id, competitor_bundle_id, competitor_name, platform, notes, created_at
		FROM competitors
		WHERE app_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var competitors []*model.Competitor
	for rows.Next() {
		var c model.Competitor
		err := rows.Scan(&c.ID, &c.AppID, &c.CompetitorBundleID, &c.CompetitorName, &c.Platform, &c.Notes, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		competitors = append(competitors, &c)
	}

	return competitors, nil
}

func (r *CompetitorRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM competitors WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

// Ranking operations

func (r *CompetitorRepository) CreateRanking(ctx context.Context, req *model.CreateCompetitorRankingRequest) (*model.CompetitorRanking, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO competitor_rankings (id, competitor_id, keyword_id, rank, recorded_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, competitor_id, keyword_id, rank, recorded_at
	`

	var cr model.CompetitorRanking
	err := r.pool.QueryRow(ctx, query,
		id, req.CompetitorID, req.KeywordID, req.Rank, now,
	).Scan(&cr.ID, &cr.CompetitorID, &cr.KeywordID, &cr.Rank, &cr.RecordedAt)

	if err != nil {
		return nil, err
	}

	return &cr, nil
}

func (r *CompetitorRepository) GetLatestRankings(ctx context.Context, competitorID string, keywordID string) (*model.CompetitorRanking, error) {
	query := `
		SELECT id, competitor_id, keyword_id, rank, recorded_at
		FROM competitor_rankings
		WHERE competitor_id = $1 AND keyword_id = $2
		ORDER BY recorded_at DESC
		LIMIT 1
	`

	var cr model.CompetitorRanking
	err := r.pool.QueryRow(ctx, query, competitorID, keywordID).Scan(
		&cr.ID, &cr.CompetitorID, &cr.KeywordID, &cr.Rank, &cr.RecordedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &cr, nil
}

func (r *CompetitorRepository) ListRankingsByKeyword(ctx context.Context, keywordID string, limit int) ([]*model.CompetitorRanking, error) {
	query := `
		SELECT id, competitor_id, keyword_id, rank, recorded_at
		FROM competitor_rankings
		WHERE keyword_id = $1
		ORDER BY recorded_at DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, keywordID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rankings []*model.CompetitorRanking
	for rows.Next() {
		var cr model.CompetitorRanking
		err := rows.Scan(&cr.ID, &cr.CompetitorID, &cr.KeywordID, &cr.Rank, &cr.RecordedAt)
		if err != nil {
			return nil, err
		}
		rankings = append(rankings, &cr)
	}

	return rankings, nil
}

// GetComparisonData returns comparison data for an app's keyword with all its competitors
func (r *CompetitorRepository) GetComparisonData(ctx context.Context, appID string, keywordID string) (*model.CompetitorComparison, error) {
	// Get keyword info
	keywordQuery := `SELECT keyword FROM keywords WHERE id = $1`
	var keyword string
	err := r.pool.QueryRow(ctx, keywordQuery, keywordID).Scan(&keyword)
	if err != nil {
		return nil, err
	}

	// Get app's latest ranking
	appRankQuery := `
		SELECT rank, recorded_at FROM ranking_history
		WHERE keyword_id = $1
		ORDER BY recorded_at DESC
		LIMIT 1
	`
	var appRank *int
	var recordedAt time.Time
	err = r.pool.QueryRow(ctx, appRankQuery, keywordID).Scan(&appRank, &recordedAt)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	// Get competitors and their latest rankings
	competitorsQuery := `
		SELECT c.id, c.competitor_name, c.competitor_bundle_id, cr.rank
		FROM competitors c
		LEFT JOIN LATERAL (
			SELECT rank FROM competitor_rankings
			WHERE competitor_id = c.id AND keyword_id = $2
			ORDER BY recorded_at DESC
			LIMIT 1
		) cr ON true
		WHERE c.app_id = $1
		ORDER BY cr.rank NULLS LAST
	`

	rows, err := r.pool.Query(ctx, competitorsQuery, appID, keywordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var competitors []model.CompetitorRankData
	for rows.Next() {
		var cd model.CompetitorRankData
		err := rows.Scan(&cd.CompetitorID, &cd.CompetitorName, &cd.BundleID, &cd.Rank)
		if err != nil {
			return nil, err
		}
		competitors = append(competitors, cd)
	}

	return &model.CompetitorComparison{
		Keyword:     keyword,
		KeywordID:   keywordID,
		AppRank:     appRank,
		Competitors: competitors,
		RecordedAt:  recordedAt,
	}, nil
}

// GetKeywordGap returns keywords where any competitor ranks <= competitorThreshold
// but our app ranks > ourThreshold or has no ranking.
func (r *CompetitorRepository) GetKeywordGap(ctx context.Context, appID string, competitorThreshold, ourThreshold int) ([]*model.KeywordGap, error) {
	query := `
		SELECT DISTINCT ON (k.id)
			k.id          AS keyword_id,
			k.keyword,
			k.country,
			c.competitor_name,
			cr_latest.rank AS competitor_rank,
			rh_latest.rank AS our_rank
		FROM keywords k
		JOIN competitors c ON c.app_id = k.app_id
		JOIN LATERAL (
			SELECT rank
			FROM competitor_rankings
			WHERE competitor_id = c.id AND keyword_id = k.id
			ORDER BY recorded_at DESC
			LIMIT 1
		) cr_latest ON true
		LEFT JOIN LATERAL (
			SELECT rank
			FROM ranking_history
			WHERE keyword_id = k.id
			ORDER BY recorded_at DESC
			LIMIT 1
		) rh_latest ON true
		WHERE k.app_id = $1
		  AND cr_latest.rank IS NOT NULL
		  AND cr_latest.rank <= $2
		  AND (rh_latest.rank IS NULL OR rh_latest.rank > $3)
		ORDER BY k.id, cr_latest.rank ASC
	`

	rows, err := r.pool.Query(ctx, query, appID, competitorThreshold, ourThreshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gaps []*model.KeywordGap
	for rows.Next() {
		g := &model.KeywordGap{}
		if err := rows.Scan(&g.KeywordID, &g.Keyword, &g.Country, &g.CompetitorName, &g.CompetitorRank, &g.OurRank); err != nil {
			return nil, err
		}
		gaps = append(gaps, g)
	}
	return gaps, nil
}
