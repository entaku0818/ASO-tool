package repository

import (
	"context"
	"encoding/json"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TemplateRepository struct {
	pool *pgxpool.Pool
}

func NewTemplateRepository(pool *pgxpool.Pool) *TemplateRepository {
	return &TemplateRepository{pool: pool}
}

func (r *TemplateRepository) List(ctx context.Context, category string) ([]*model.Template, error) {
	query := `
		SELECT id, name, category, description, device, style, is_pro, sort_order, created_at
		FROM screenshot_templates
	`
	args := []any{}

	if category != "" {
		query += ` WHERE category = $1`
		args = append(args, category)
	}

	query += ` ORDER BY sort_order ASC, created_at ASC`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*model.Template
	for rows.Next() {
		t := &model.Template{}
		var styleJSON []byte
		if err := rows.Scan(
			&t.ID, &t.Name, &t.Category, &t.Description,
			&t.Device, &styleJSON, &t.IsPro, &t.SortOrder, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(styleJSON, &t.Style); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}

	return templates, nil
}
