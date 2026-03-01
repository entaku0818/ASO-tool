package repository

import (
	"context"
	"errors"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SearchAdsCredentialsRepository struct {
	pool *pgxpool.Pool
}

func NewSearchAdsCredentialsRepository(pool *pgxpool.Pool) *SearchAdsCredentialsRepository {
	return &SearchAdsCredentialsRepository{pool: pool}
}

func (r *SearchAdsCredentialsRepository) Save(ctx context.Context, appID, clientID, teamID, keyID string, encryptedKey, nonce []byte, orgID string, adamID int64) (*model.SearchAdsCredentials, error) {
	// Upsert pattern: insert or update on conflict
	id := uuid.New().String()
	query := `
		INSERT INTO search_ads_credentials
			(id, app_id, client_id, team_id, key_id, private_key_encrypted, encryption_nonce, org_id, adam_id, is_valid, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, true, NOW(), NOW())
		ON CONFLICT (app_id) DO UPDATE
		SET client_id = EXCLUDED.client_id,
		    team_id = EXCLUDED.team_id,
		    key_id = EXCLUDED.key_id,
		    private_key_encrypted = EXCLUDED.private_key_encrypted,
		    encryption_nonce = EXCLUDED.encryption_nonce,
		    org_id = EXCLUDED.org_id,
		    adam_id = EXCLUDED.adam_id,
		    is_valid = true,
		    updated_at = NOW()
		RETURNING id, app_id, client_id, team_id, key_id, private_key_encrypted, encryption_nonce,
		          org_id, adam_id, is_valid, created_at, updated_at
	`

	cred := &model.SearchAdsCredentials{}
	err := r.pool.QueryRow(ctx, query,
		id, appID, clientID, teamID, keyID, encryptedKey, nonce, orgID, adamID,
	).Scan(
		&cred.ID, &cred.AppID, &cred.ClientID, &cred.TeamID, &cred.KeyID,
		&cred.PrivateKeyEncrypted, &cred.EncryptionNonce,
		&cred.OrgID, &cred.AdamID, &cred.IsValid, &cred.CreatedAt, &cred.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return cred, nil
}

func (r *SearchAdsCredentialsRepository) GetByAppID(ctx context.Context, appID string) (*model.SearchAdsCredentials, error) {
	query := `
		SELECT id, app_id, client_id, team_id, key_id, private_key_encrypted, encryption_nonce,
		       org_id, adam_id, is_valid, created_at, updated_at
		FROM search_ads_credentials
		WHERE app_id = $1
	`

	cred := &model.SearchAdsCredentials{}
	err := r.pool.QueryRow(ctx, query, appID).Scan(
		&cred.ID, &cred.AppID, &cred.ClientID, &cred.TeamID, &cred.KeyID,
		&cred.PrivateKeyEncrypted, &cred.EncryptionNonce,
		&cred.OrgID, &cred.AdamID, &cred.IsValid, &cred.CreatedAt, &cred.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrSearchAdsCredentialsNotFound
	}
	if err != nil {
		return nil, err
	}

	return cred, nil
}

func (r *SearchAdsCredentialsRepository) Delete(ctx context.Context, appID string) error {
	query := `DELETE FROM search_ads_credentials WHERE app_id = $1`

	result, err := r.pool.Exec(ctx, query, appID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrSearchAdsCredentialsNotFound
	}

	return nil
}

func (r *SearchAdsCredentialsRepository) UpdateValidity(ctx context.Context, appID string, isValid bool) error {
	query := `
		UPDATE search_ads_credentials
		SET is_valid = $2, updated_at = NOW()
		WHERE app_id = $1
	`

	result, err := r.pool.Exec(ctx, query, appID, isValid)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrSearchAdsCredentialsNotFound
	}

	return nil
}
