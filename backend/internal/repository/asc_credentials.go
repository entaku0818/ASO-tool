package repository

import (
	"context"
	"errors"
	"time"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ASCCredentialsRepository struct {
	pool *pgxpool.Pool
}

func NewASCCredentialsRepository(pool *pgxpool.Pool) *ASCCredentialsRepository {
	return &ASCCredentialsRepository{pool: pool}
}

func (r *ASCCredentialsRepository) Create(ctx context.Context, appID string, issuerID, keyID string, encryptedKey, nonce []byte) (*model.ASCCredentials, error) {
	cred := &model.ASCCredentials{
		ID:                  uuid.New().String(),
		AppID:               appID,
		IssuerID:            issuerID,
		KeyID:               keyID,
		PrivateKeyEncrypted: encryptedKey,
		EncryptionNonce:     nonce,
		IsValid:             true,
	}

	query := `
		INSERT INTO asc_credentials (id, app_id, issuer_id, key_id, private_key_encrypted, encryption_nonce, is_valid, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		cred.ID, cred.AppID, cred.IssuerID, cred.KeyID,
		cred.PrivateKeyEncrypted, cred.EncryptionNonce, cred.IsValid,
	).Scan(&cred.CreatedAt, &cred.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return cred, nil
}

func (r *ASCCredentialsRepository) GetByAppID(ctx context.Context, appID string) (*model.ASCCredentials, error) {
	query := `
		SELECT id, app_id, issuer_id, key_id, private_key_encrypted, encryption_nonce,
		       is_valid, last_validated_at, created_at, updated_at
		FROM asc_credentials
		WHERE app_id = $1
	`

	cred := &model.ASCCredentials{}
	err := r.pool.QueryRow(ctx, query, appID).Scan(
		&cred.ID, &cred.AppID, &cred.IssuerID, &cred.KeyID,
		&cred.PrivateKeyEncrypted, &cred.EncryptionNonce,
		&cred.IsValid, &cred.LastValidatedAt, &cred.CreatedAt, &cred.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrCredentialsNotFound
	}
	if err != nil {
		return nil, err
	}

	return cred, nil
}

func (r *ASCCredentialsRepository) Update(ctx context.Context, appID string, issuerID, keyID string, encryptedKey, nonce []byte) (*model.ASCCredentials, error) {
	query := `
		UPDATE asc_credentials
		SET issuer_id = $2, key_id = $3, private_key_encrypted = $4, encryption_nonce = $5,
		    is_valid = true, updated_at = NOW()
		WHERE app_id = $1
		RETURNING id, app_id, issuer_id, key_id, private_key_encrypted, encryption_nonce,
		          is_valid, last_validated_at, created_at, updated_at
	`

	cred := &model.ASCCredentials{}
	err := r.pool.QueryRow(ctx, query, appID, issuerID, keyID, encryptedKey, nonce).Scan(
		&cred.ID, &cred.AppID, &cred.IssuerID, &cred.KeyID,
		&cred.PrivateKeyEncrypted, &cred.EncryptionNonce,
		&cred.IsValid, &cred.LastValidatedAt, &cred.CreatedAt, &cred.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrCredentialsNotFound
	}
	if err != nil {
		return nil, err
	}

	return cred, nil
}

func (r *ASCCredentialsRepository) Delete(ctx context.Context, appID string) error {
	query := `DELETE FROM asc_credentials WHERE app_id = $1`

	result, err := r.pool.Exec(ctx, query, appID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrCredentialsNotFound
	}

	return nil
}

func (r *ASCCredentialsRepository) SetValidationStatus(ctx context.Context, appID string, isValid bool) error {
	query := `
		UPDATE asc_credentials
		SET is_valid = $2, last_validated_at = $3, updated_at = NOW()
		WHERE app_id = $1
	`

	now := time.Now()
	result, err := r.pool.Exec(ctx, query, appID, isValid, now)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrCredentialsNotFound
	}

	return nil
}

func (r *ASCCredentialsRepository) ListAppsWithCredentials(ctx context.Context) ([]*model.ASCCredentials, error) {
	query := `
		SELECT id, app_id, issuer_id, key_id, private_key_encrypted, encryption_nonce,
		       is_valid, last_validated_at, created_at, updated_at
		FROM asc_credentials
		WHERE is_valid = true
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creds []*model.ASCCredentials
	for rows.Next() {
		cred := &model.ASCCredentials{}
		err := rows.Scan(
			&cred.ID, &cred.AppID, &cred.IssuerID, &cred.KeyID,
			&cred.PrivateKeyEncrypted, &cred.EncryptionNonce,
			&cred.IsValid, &cred.LastValidatedAt, &cred.CreatedAt, &cred.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		creds = append(creds, cred)
	}

	return creds, nil
}
