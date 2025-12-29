package model

import (
	"errors"
	"testing"
)

func TestCreateAppRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateAppRequest
		wantErr error
	}{
		{
			name: "valid iOS request",
			req: CreateAppRequest{
				Name:     "Test App",
				BundleID: "com.example.test",
				Platform: PlatformIOS,
			},
			wantErr: nil,
		},
		{
			name: "valid Android request",
			req: CreateAppRequest{
				Name:     "Test App",
				BundleID: "com.example.test",
				Platform: PlatformAndroid,
			},
			wantErr: nil,
		},
		{
			name: "valid with store URL",
			req: CreateAppRequest{
				Name:     "Test App",
				BundleID: "com.example.test",
				Platform: PlatformIOS,
				StoreURL: "https://apps.apple.com/app/id123456",
			},
			wantErr: nil,
		},
		{
			name: "missing name",
			req: CreateAppRequest{
				BundleID: "com.example.test",
				Platform: PlatformIOS,
			},
			wantErr: ErrNameRequired,
		},
		{
			name: "empty name",
			req: CreateAppRequest{
				Name:     "",
				BundleID: "com.example.test",
				Platform: PlatformIOS,
			},
			wantErr: ErrNameRequired,
		},
		{
			name: "missing bundle_id",
			req: CreateAppRequest{
				Name:     "Test App",
				Platform: PlatformIOS,
			},
			wantErr: ErrBundleIDRequired,
		},
		{
			name: "empty bundle_id",
			req: CreateAppRequest{
				Name:     "Test App",
				BundleID: "",
				Platform: PlatformIOS,
			},
			wantErr: ErrBundleIDRequired,
		},
		{
			name: "invalid platform",
			req: CreateAppRequest{
				Name:     "Test App",
				BundleID: "com.example.test",
				Platform: "windows",
			},
			wantErr: ErrInvalidPlatform,
		},
		{
			name: "empty platform",
			req: CreateAppRequest{
				Name:     "Test App",
				BundleID: "com.example.test",
				Platform: "",
			},
			wantErr: ErrInvalidPlatform,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlatformConstants(t *testing.T) {
	if PlatformIOS != "ios" {
		t.Errorf("PlatformIOS = %v, want 'ios'", PlatformIOS)
	}
	if PlatformAndroid != "android" {
		t.Errorf("PlatformAndroid = %v, want 'android'", PlatformAndroid)
	}
}
