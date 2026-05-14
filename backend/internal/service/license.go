package service

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// charset excludes visually ambiguous characters (0/O, 1/I/l)
const licenseCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

type LicenseService struct {
	repo        *repository.LicenseRepository
	userRepo    *repository.UserRepository
	authService *AuthService
}

func NewLicenseService(
	repo *repository.LicenseRepository,
	userRepo *repository.UserRepository,
	authService *AuthService,
) *LicenseService {
	return &LicenseService{repo: repo, userRepo: userRepo, authService: authService}
}

func (s *LicenseService) Generate(ctx context.Context, req *model.GenerateLicenseRequest) (*model.LicenseKey, error) {
	key, err := generateKey()
	if err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, key, req.Email, req.StripeSessionID)
}

func (s *LicenseService) Activate(ctx context.Context, req *model.ActivateLicenseRequest) (*model.ActivateLicenseResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	lk, err := s.repo.GetByKey(ctx, req.Key)
	if err != nil {
		return nil, model.ErrLicenseNotFound
	}

	// If already activated by a different email, reject
	if lk.IsActive && lk.UserID != nil {
		existing, err := s.userRepo.GetByID(ctx, *lk.UserID)
		if err == nil && existing.Email != req.Email {
			return nil, model.ErrLicenseAlreadyUsed
		}
	}

	// Find or create user
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Auto-create account with random password (user can reset via web)
		randomPass, genErr := generateRandomString(24)
		if genErr != nil {
			return nil, genErr
		}
		hashed, genErr := bcrypt.GenerateFromPassword([]byte(randomPass), bcrypt.DefaultCost)
		if genErr != nil {
			return nil, genErr
		}
		user, err = s.userRepo.Create(ctx, req.Email, string(hashed), req.Email, false)
		if err != nil {
			return nil, err
		}
	}

	// Mark license active
	if err := s.repo.Activate(ctx, lk.ID, user.ID); err != nil {
		return nil, err
	}

	token, err := s.authService.GenerateTokenForUser(user)
	if err != nil {
		return nil, err
	}

	return &model.ActivateLicenseResponse{
		Token: token,
		User:  user,
		Key:   lk.Key,
	}, nil
}

func (s *LicenseService) GetByUserID(ctx context.Context, userID string) (*model.LicenseKey, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func generateKey() (string, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	chars := make([]byte, 12)
	for i, b := range buf {
		chars[i] = licenseCharset[int(b)%len(licenseCharset)]
	}
	return fmt.Sprintf("ASOT-%c%c%c%c-%c%c%c%c-%c%c%c%c",
		chars[0], chars[1], chars[2], chars[3],
		chars[4], chars[5], chars[6], chars[7],
		chars[8], chars[9], chars[10], chars[11],
	), nil
}

func generateRandomString(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	chars := make([]byte, n)
	for i, b := range buf {
		chars[i] = licenseCharset[int(b)%len(licenseCharset)]
	}
	return string(chars), nil
}
