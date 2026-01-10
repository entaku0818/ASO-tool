package service

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/appstoreconnect"
	"github.com/entaku0818/aso-tool/backend/internal/crypto"
	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
)

type AnalyticsService struct {
	credentialsRepo *repository.ASCCredentialsRepository
	analyticsRepo   *repository.AnalyticsRepository
	appVersionRepo  *repository.AppVersionRepository
	rankingRepo     *repository.RankingRepository
	reviewRepo      *repository.ReviewRepository
	appRepo         *repository.AppRepository
}

func NewAnalyticsService(
	credentialsRepo *repository.ASCCredentialsRepository,
	analyticsRepo *repository.AnalyticsRepository,
	appVersionRepo *repository.AppVersionRepository,
	rankingRepo *repository.RankingRepository,
	reviewRepo *repository.ReviewRepository,
	appRepo *repository.AppRepository,
) *AnalyticsService {
	return &AnalyticsService{
		credentialsRepo: credentialsRepo,
		analyticsRepo:   analyticsRepo,
		appVersionRepo:  appVersionRepo,
		rankingRepo:     rankingRepo,
		reviewRepo:      reviewRepo,
		appRepo:         appRepo,
	}
}

// SetCredentials stores App Store Connect API credentials
func (s *AnalyticsService) SetCredentials(ctx context.Context, userID, appID string, req *model.CreateASCCredentialsRequest) error {
	// Verify app belongs to user and is iOS
	app, err := s.appRepo.Get(ctx, userID, appID)
	if err != nil {
		return err
	}
	if app.Platform != model.PlatformIOS {
		return model.ErrIOSOnly
	}

	// Decode base64 private key
	privateKeyBytes, err := base64.StdEncoding.DecodeString(req.PrivateKey)
	if err != nil {
		return model.ErrInvalidPrivateKey
	}

	// Encrypt the private key
	encryptedKey, nonce, err := crypto.Encrypt(privateKeyBytes)
	if err != nil {
		return err
	}

	// Check if credentials already exist
	existing, err := s.credentialsRepo.GetByAppID(ctx, appID)
	if err != nil && err != model.ErrCredentialsNotFound {
		return err
	}

	if existing != nil {
		// Update existing
		_, err = s.credentialsRepo.Update(ctx, appID, req.IssuerID, req.KeyID, encryptedKey, nonce)
	} else {
		// Create new
		_, err = s.credentialsRepo.Create(ctx, appID, req.IssuerID, req.KeyID, encryptedKey, nonce)
	}

	return err
}

// GetCredentials retrieves credentials metadata (without private key)
func (s *AnalyticsService) GetCredentials(ctx context.Context, userID, appID string) (*model.ASCCredentials, error) {
	// Verify app belongs to user
	_, err := s.appRepo.Get(ctx, userID, appID)
	if err != nil {
		return nil, err
	}

	return s.credentialsRepo.GetByAppID(ctx, appID)
}

// DeleteCredentials removes stored credentials
func (s *AnalyticsService) DeleteCredentials(ctx context.Context, userID, appID string) error {
	// Verify app belongs to user
	_, err := s.appRepo.Get(ctx, userID, appID)
	if err != nil {
		return err
	}

	return s.credentialsRepo.Delete(ctx, appID)
}

// ValidateCredentials tests if stored credentials are valid
func (s *AnalyticsService) ValidateCredentials(ctx context.Context, userID, appID string) error {
	// Verify app belongs to user
	_, err := s.appRepo.Get(ctx, userID, appID)
	if err != nil {
		return err
	}

	creds, err := s.credentialsRepo.GetByAppID(ctx, appID)
	if err != nil {
		return err
	}

	// Decrypt private key
	privateKey, err := crypto.Decrypt(creds.PrivateKeyEncrypted, creds.EncryptionNonce)
	if err != nil {
		return err
	}

	// Create client and validate
	client, err := appstoreconnect.NewClient(appstoreconnect.Config{
		IssuerID:      creds.IssuerID,
		KeyID:         creds.KeyID,
		PrivateKeyPEM: privateKey,
	})
	if err != nil {
		_ = s.credentialsRepo.SetValidationStatus(ctx, appID, false)
		return model.ErrInvalidASCCredentials
	}

	if err := client.ValidateCredentials(ctx); err != nil {
		_ = s.credentialsRepo.SetValidationStatus(ctx, appID, false)
		return model.ErrInvalidASCCredentials
	}

	_ = s.credentialsRepo.SetValidationStatus(ctx, appID, true)
	return nil
}

// FetchAnalytics fetches analytics from App Store Connect and stores it
func (s *AnalyticsService) FetchAnalytics(ctx context.Context, appID string, startDate, endDate time.Time) (int, error) {
	app, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return 0, err
	}

	creds, err := s.credentialsRepo.GetByAppID(ctx, appID)
	if err != nil {
		return 0, err
	}

	// Decrypt private key
	privateKey, err := crypto.Decrypt(creds.PrivateKeyEncrypted, creds.EncryptionNonce)
	if err != nil {
		return 0, err
	}

	// Create client
	client, err := appstoreconnect.NewClient(appstoreconnect.Config{
		IssuerID:      creds.IssuerID,
		KeyID:         creds.KeyID,
		PrivateKeyPEM: privateKey,
	})
	if err != nil {
		return 0, err
	}

	// Get App Store app ID from bundle ID
	ascApp, err := client.GetAppByBundleID(ctx, app.BundleID)
	if err != nil {
		return 0, err
	}

	// Fetch analytics
	analyticsData, err := client.GetAnalytics(ctx, ascApp.ID, startDate, endDate)
	if err != nil {
		return 0, err
	}

	// Store analytics
	count := 0
	for _, data := range analyticsData {
		date, _ := time.Parse("2006-01-02", data.Date)
		var conversionRate *float64
		if data.Impressions > 0 {
			rate := float64(data.Downloads) / float64(data.Impressions)
			conversionRate = &rate
		}

		analytics := &model.ASCAnalytics{
			AppID:          appID,
			Date:           date,
			Impressions:    data.Impressions,
			Downloads:      data.Downloads,
			PageViews:      data.PageViews,
			ConversionRate: conversionRate,
		}

		if err := s.analyticsRepo.Upsert(ctx, analytics); err != nil {
			continue
		}
		count++
	}

	// Also fetch and store version info
	versions, err := client.GetAppVersions(ctx, ascApp.ID)
	if err == nil {
		for _, v := range versions {
			releaseDate, _ := time.Parse(time.RFC3339, v.ReleaseDate)
			appVersion := &model.AppVersion{
				AppID:       appID,
				Version:     v.Version,
				ReleaseDate: releaseDate,
			}
			_ = s.appVersionRepo.Create(ctx, appVersion)
		}
	}

	return count, nil
}

// FetchAllAppsAnalytics fetches analytics for all apps with valid credentials
func (s *AnalyticsService) FetchAllAppsAnalytics(ctx context.Context) (map[string]int, error) {
	creds, err := s.credentialsRepo.ListAppsWithCredentials(ctx)
	if err != nil {
		return nil, err
	}

	results := make(map[string]int)
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	for _, cred := range creds {
		count, err := s.FetchAnalytics(ctx, cred.AppID, startDate, endDate)
		if err != nil {
			results[cred.AppID] = -1 // Indicate error
			continue
		}
		results[cred.AppID] = count
	}

	return results, nil
}

// GetAnalytics retrieves stored analytics data
func (s *AnalyticsService) GetAnalytics(ctx context.Context, userID, appID string, days int) ([]*model.ASCAnalytics, error) {
	// Verify app belongs to user
	_, err := s.appRepo.Get(ctx, userID, appID)
	if err != nil {
		return nil, err
	}

	return s.analyticsRepo.GetLatest(ctx, appID, days)
}

// GetAnalyticsWithCorrelation retrieves analytics with correlation indicators
func (s *AnalyticsService) GetAnalyticsWithCorrelation(ctx context.Context, userID, appID string, days int) ([]*model.AnalyticsWithCorrelation, error) {
	// Verify app belongs to user
	_, err := s.appRepo.Get(ctx, userID, appID)
	if err != nil {
		return nil, err
	}

	analytics, err := s.analyticsRepo.GetLatest(ctx, appID, days)
	if err != nil {
		return nil, err
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// Get versions released in this period
	versions, err := s.appVersionRepo.GetByDateRange(ctx, appID, startDate, endDate)
	if err != nil {
		versions = []*model.AppVersion{}
	}

	// Build version map by date
	versionByDate := make(map[string]*model.AppVersion)
	for _, v := range versions {
		dateKey := v.ReleaseDate.Format("2006-01-02")
		versionByDate[dateKey] = v
	}

	// Get rankings for this period
	rankings, err := s.rankingRepo.ListByAppWithKeyword(ctx, appID, startDate, endDate)
	if err != nil {
		rankings = []*model.RankingWithKeyword{}
	}

	// Build ranking changes map
	rankingChanges := s.calculateRankingChanges(rankings)

	// Build result with correlations
	result := make([]*model.AnalyticsWithCorrelation, len(analytics))
	for i, a := range analytics {
		corr := &model.AnalyticsWithCorrelation{
			ASCAnalytics: *a,
		}

		dateKey := a.Date.Format("2006-01-02")

		// Add version if released on this date
		if v, ok := versionByDate[dateKey]; ok {
			corr.NewVersion = &v.Version
		}

		// Add ranking change if significant
		if change, ok := rankingChanges[dateKey]; ok && change != 0 {
			corr.RankingChange = &change
		}

		result[i] = corr
	}

	return result, nil
}

// calculateRankingChanges calculates daily ranking changes
func (s *AnalyticsService) calculateRankingChanges(rankings []*model.RankingWithKeyword) map[string]int {
	changes := make(map[string]int)

	// Group rankings by date and find average change
	type dailyRank struct {
		sum   int
		count int
	}
	dailyRanks := make(map[string]*dailyRank)

	for _, r := range rankings {
		if r.Rank == nil {
			continue
		}
		dateKey := r.RecordedAt.Format("2006-01-02")
		if dr, ok := dailyRanks[dateKey]; ok {
			dr.sum += *r.Rank
			dr.count++
		} else {
			dailyRanks[dateKey] = &dailyRank{sum: *r.Rank, count: 1}
		}
	}

	// Calculate changes between consecutive days
	var prevAvg float64
	var prevDate string
	for dateKey, dr := range dailyRanks {
		avg := float64(dr.sum) / float64(dr.count)
		if prevDate != "" {
			// Ranking improvement is negative change (lower rank is better)
			change := int(prevAvg - avg)
			if change != 0 {
				changes[dateKey] = change
			}
		}
		prevAvg = avg
		prevDate = dateKey
	}

	return changes
}

// GetSummary retrieves analytics summary
func (s *AnalyticsService) GetSummary(ctx context.Context, userID, appID string, days int) (*model.AnalyticsSummary, error) {
	// Verify app belongs to user
	_, err := s.appRepo.Get(ctx, userID, appID)
	if err != nil {
		return nil, err
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	return s.analyticsRepo.GetSummary(ctx, appID, startDate, endDate)
}
