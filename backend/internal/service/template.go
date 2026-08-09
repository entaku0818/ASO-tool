package service

import (
	"context"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/entaku0818/aso-compass/backend/internal/repository"
)

type TemplateService struct {
	repo *repository.TemplateRepository
}

func NewTemplateService(repo *repository.TemplateRepository) *TemplateService {
	return &TemplateService{repo: repo}
}

func (s *TemplateService) List(ctx context.Context, category string) ([]*model.Template, error) {
	return s.repo.List(ctx, category)
}
