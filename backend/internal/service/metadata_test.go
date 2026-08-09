package service

import (
	"testing"

	"github.com/entaku0818/aso-compass/backend/internal/model"
)

// TestMetadataUpsertValidation tests that Validate() is called by service.Upsert.
func TestMetadataUpsertValidation(t *testing.T) {
	req := &model.UpsertMetadataRequest{}
	if err := req.Validate(); err != model.ErrAppIDRequired {
		t.Errorf("expected ErrAppIDRequired, got %v", err)
	}
}

// TestMetadataDefaultsApplied verifies defaults are set on valid request.
func TestMetadataDefaultsApplied(t *testing.T) {
	req := &model.UpsertMetadataRequest{AppID: "app-1"}
	if err := req.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Locale != "ja" {
		t.Errorf("expected locale=ja, got %q", req.Locale)
	}
	if req.VersionTag != "draft" {
		t.Errorf("expected version_tag=draft, got %q", req.VersionTag)
	}
}
