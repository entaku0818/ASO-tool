package service

import (
	"context"
	"encoding/base64"

	"github.com/entaku0818/aso-tool/backend/internal/crypto"
	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
	"github.com/entaku0818/aso-tool/backend/internal/scraper"
)

type KeywordPopularityService struct {
	credRepo    *repository.SearchAdsCredentialsRepository
	keywordRepo *repository.KeywordRepository
	appRepo     *repository.AppRepository
}

func NewKeywordPopularityService(
	credRepo *repository.SearchAdsCredentialsRepository,
	keywordRepo *repository.KeywordRepository,
	appRepo *repository.AppRepository,
) *KeywordPopularityService {
	return &KeywordPopularityService{
		credRepo:    credRepo,
		keywordRepo: keywordRepo,
		appRepo:     appRepo,
	}
}

// SetCredentials stores Search Ads API credentials for an app
func (s *KeywordPopularityService) SetCredentials(ctx context.Context, userID, appID string, req *model.CreateSearchAdsCredentialsRequest) error {
	// Verify app belongs to user
	_, err := s.appRepo.Get(ctx, userID, appID)
	if err != nil {
		return err
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

	_, err = s.credRepo.Save(ctx, appID, req.ClientID, req.TeamID, req.KeyID, encryptedKey, nonce, req.OrgID, req.AdamID)
	return err
}

// GetCredentials returns credentials metadata (no private key)
func (s *KeywordPopularityService) GetCredentials(ctx context.Context, userID, appID string) (*model.SearchAdsCredentials, error) {
	// Verify app belongs to user
	_, err := s.appRepo.Get(ctx, userID, appID)
	if err != nil {
		return nil, err
	}

	return s.credRepo.GetByAppID(ctx, appID)
}

// DeleteCredentials removes stored credentials
func (s *KeywordPopularityService) DeleteCredentials(ctx context.Context, userID, appID string) error {
	// Verify app belongs to user
	_, err := s.appRepo.Get(ctx, userID, appID)
	if err != nil {
		return err
	}

	return s.credRepo.Delete(ctx, appID)
}

// RefreshScores updates popularity scores for all tracked keywords of the app
func (s *KeywordPopularityService) RefreshScores(ctx context.Context, userID, appID string) (int, error) {
	// Verify app belongs to user
	_, err := s.appRepo.Get(ctx, userID, appID)
	if err != nil {
		return 0, err
	}

	client, err := s.buildClient(ctx, appID)
	if err != nil {
		return 0, err
	}

	creds, err := s.credRepo.GetByAppID(ctx, appID)
	if err != nil {
		return 0, err
	}

	keywords, err := s.keywordRepo.ListByApp(ctx, appID)
	if err != nil {
		return 0, err
	}

	if len(keywords) == 0 {
		return 0, nil
	}

	// Collect keyword texts
	texts := make([]string, len(keywords))
	for i, kw := range keywords {
		texts[i] = kw.Keyword
	}

	// Fetch scores in one API call
	results, err := client.GetKeywordPopularity(ctx, creds.AdamID, texts, len(texts))
	if err != nil {
		_ = s.credRepo.UpdateValidity(ctx, appID, false)
		return 0, err
	}

	// Build score map
	scoreMap := make(map[string]int, len(results))
	for _, r := range results {
		scoreMap[r.Text] = r.PopularityScore
	}

	// Update each keyword
	updated := 0
	for _, kw := range keywords {
		score, ok := scoreMap[kw.Keyword]
		if !ok {
			continue
		}
		if err := s.keywordRepo.UpdatePopularity(ctx, kw.ID, score); err != nil {
			continue
		}
		updated++
	}

	return updated, nil
}

// GetSuggestions returns keyword suggestions with popularity scores for an app
func (s *KeywordPopularityService) GetSuggestions(ctx context.Context, userID, appID string, limit int) ([]scraper.KeywordPopularity, error) {
	// Verify app belongs to user
	_, err := s.appRepo.Get(ctx, userID, appID)
	if err != nil {
		return nil, err
	}

	client, err := s.buildClient(ctx, appID)
	if err != nil {
		return nil, err
	}

	creds, err := s.credRepo.GetByAppID(ctx, appID)
	if err != nil {
		return nil, err
	}

	// Pass empty keywords slice to get suggestions from Apple
	results, err := client.GetKeywordPopularity(ctx, creds.AdamID, nil, limit)
	if err != nil {
		_ = s.credRepo.UpdateValidity(ctx, appID, false)
		return nil, err
	}

	return results, nil
}

// GetCompetitorSuggestions returns keyword suggestions for a competitor app by Adam ID
func (s *KeywordPopularityService) GetCompetitorSuggestions(ctx context.Context, userID, appID string, competitorAdamID int64, limit int) ([]scraper.KeywordPopularity, error) {
	// Verify app belongs to user (credentials are stored under this app)
	_, err := s.appRepo.Get(ctx, userID, appID)
	if err != nil {
		return nil, err
	}

	client, err := s.buildClient(ctx, appID)
	if err != nil {
		return nil, err
	}

	results, err := client.GetKeywordPopularity(ctx, competitorAdamID, nil, limit)
	if err != nil {
		_ = s.credRepo.UpdateValidity(ctx, appID, false)
		return nil, err
	}

	return results, nil
}

// buildClient creates a SearchAdsClient from stored credentials
func (s *KeywordPopularityService) buildClient(ctx context.Context, appID string) (*scraper.SearchAdsClient, error) {
	creds, err := s.credRepo.GetByAppID(ctx, appID)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.Decrypt(creds.PrivateKeyEncrypted, creds.EncryptionNonce)
	if err != nil {
		return nil, err
	}

	return scraper.NewSearchAdsClient(creds.ClientID, creds.TeamID, creds.KeyID, privateKey, creds.OrgID)
}
