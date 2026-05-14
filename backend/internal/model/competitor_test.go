package model

import (
	"errors"
	"testing"
)

func TestCreateCompetitorRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateCompetitorRequest
		wantErr error
	}{
		{
			name: "valid iOS competitor",
			req: CreateCompetitorRequest{
				AppID:              "app-1",
				CompetitorBundleID: "com.rival.app",
				CompetitorName:     "Rival App",
				Platform:           PlatformIOS,
			},
			wantErr: nil,
		},
		{
			name: "valid Android competitor",
			req: CreateCompetitorRequest{
				AppID:              "app-1",
				CompetitorBundleID: "com.rival.app",
				CompetitorName:     "Rival App",
				Platform:           PlatformAndroid,
			},
			wantErr: nil,
		},
		{
			name:    "missing app_id",
			req:     CreateCompetitorRequest{CompetitorBundleID: "com.rival.app", CompetitorName: "Rival", Platform: PlatformIOS},
			wantErr: ErrAppIDRequired,
		},
		{
			name:    "missing bundle_id",
			req:     CreateCompetitorRequest{AppID: "app-1", CompetitorName: "Rival", Platform: PlatformIOS},
			wantErr: ErrBundleIDRequired,
		},
		{
			name:    "missing name",
			req:     CreateCompetitorRequest{AppID: "app-1", CompetitorBundleID: "com.rival.app", Platform: PlatformIOS},
			wantErr: ErrNameRequired,
		},
		{
			name:    "invalid platform",
			req:     CreateCompetitorRequest{AppID: "app-1", CompetitorBundleID: "com.rival.app", CompetitorName: "Rival", Platform: "windows"},
			wantErr: ErrInvalidPlatform,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeywordGap_Fields(t *testing.T) {
	rank := 45
	g := KeywordGap{
		KeywordID:      "kw-1",
		Keyword:        "running app",
		Country:        "JP",
		CompetitorName: "Rival",
		CompetitorRank: 8,
		OurRank:        &rank,
	}

	if g.CompetitorRank != 8 {
		t.Errorf("expected competitor_rank=8, got %d", g.CompetitorRank)
	}
	if *g.OurRank != 45 {
		t.Errorf("expected our_rank=45, got %d", *g.OurRank)
	}
}

func TestKeywordGap_NilOurRank(t *testing.T) {
	g := KeywordGap{
		KeywordID:      "kw-2",
		Keyword:        "fitness",
		Country:        "US",
		CompetitorName: "Rival",
		CompetitorRank: 3,
		OurRank:        nil,
	}

	if g.OurRank != nil {
		t.Errorf("expected nil our_rank")
	}
}
