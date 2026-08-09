package service

import (
	"context"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/entaku0818/aso-compass/backend/internal/repository"
)

type ScreenshotService struct {
	repo *repository.ScreenshotRepository
}

func NewScreenshotService(repo *repository.ScreenshotRepository) *ScreenshotService {
	return &ScreenshotService{repo: repo}
}

func (s *ScreenshotService) Create(ctx context.Context, req *model.CreateScreenshotRequest) (*model.Screenshot, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// If display_order is 0, set it to the next available order
	if req.DisplayOrder == 0 {
		maxOrder, err := s.repo.GetMaxDisplayOrder(ctx, req.AppID)
		if err != nil {
			return nil, err
		}
		req.DisplayOrder = maxOrder + 1
	}

	return s.repo.Create(ctx, req)
}

func (s *ScreenshotService) Get(ctx context.Context, id string) (*model.Screenshot, error) {
	return s.repo.Get(ctx, id)
}

func (s *ScreenshotService) ListByApp(ctx context.Context, appID string) ([]*model.Screenshot, error) {
	return s.repo.ListByApp(ctx, appID)
}

func (s *ScreenshotService) ListByAppAndDevice(ctx context.Context, appID string, deviceType model.DeviceType) ([]*model.Screenshot, error) {
	if !deviceType.IsValid() {
		return nil, model.ErrInvalidDeviceType
	}
	return s.repo.ListByAppAndDevice(ctx, appID, deviceType)
}

func (s *ScreenshotService) Update(ctx context.Context, id string, req *model.UpdateScreenshotRequest) (*model.Screenshot, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, req)
}

func (s *ScreenshotService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *ScreenshotService) Reorder(ctx context.Context, appID string, req *model.ReorderScreenshotsRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	return s.repo.Reorder(ctx, appID, req.ScreenshotIDs)
}
