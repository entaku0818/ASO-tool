package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os"

	"github.com/entaku0818/aso-tool/backend/internal/model"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
	"github.com/stripe/stripe-go/v76"
	stripeSession "github.com/stripe/stripe-go/v76/checkout/session"
	"golang.org/x/crypto/bcrypt"
)

// charset excludes visually ambiguous characters (0/O, 1/I/l)
const licenseCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

type LicenseService struct {
	repo            *repository.LicenseRepository
	userRepo        *repository.UserRepository
	authService     *AuthService
	emailService    *EmailService
	priceLicense    string
	frontendBaseURL string
}

func NewLicenseService(
	repo *repository.LicenseRepository,
	userRepo *repository.UserRepository,
	authService *AuthService,
) *LicenseService {
	priceLicense := os.Getenv("STRIPE_PRICE_LICENSE")
	if priceLicense == "" {
		log.Println("warn: STRIPE_PRICE_LICENSE is not set — license checkout will fail")
	}
	return &LicenseService{
		repo:            repo,
		userRepo:        userRepo,
		authService:     authService,
		emailService:    NewEmailService(),
		priceLicense:    priceLicense,
		frontendBaseURL: getEnvOrDefault("FRONTEND_BASE_URL", "http://localhost:3000"),
	}
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

func (s *LicenseService) List(ctx context.Context) ([]*model.LicenseKey, error) {
	return s.repo.List(ctx)
}

// CreateCheckoutSession creates a Stripe one-time payment session for a license key.
func (s *LicenseService) CreateCheckoutSession(ctx context.Context, email string) (string, error) {
	if stripe.Key == "" {
		return "", model.ErrStripeKeyRequired
	}
	if s.priceLicense == "" {
		return "", fmt.Errorf("STRIPE_PRICE_LICENSE is not configured")
	}

	params := &stripe.CheckoutSessionParams{
		CustomerEmail: stripe.String(email),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{Price: stripe.String(s.priceLicense), Quantity: stripe.Int64(1)},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(s.frontendBaseURL + "/buy/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(s.frontendBaseURL + "/buy"),
		Metadata:   map[string]string{"purchase_type": "license"},
	}

	sess, err := stripeSession.New(params)
	if err != nil {
		return "", fmt.Errorf("create checkout session: %w", err)
	}
	return sess.URL, nil
}

// GenerateAndSendLicense generates a license key and sends it by email.
// Called from the Stripe webhook after payment is confirmed.
func (s *LicenseService) GenerateAndSendLicense(ctx context.Context, email, stripeSessionID string) (*model.LicenseKey, error) {
	key, err := generateKey()
	if err != nil {
		return nil, err
	}
	sid := stripeSessionID
	lk, err := s.repo.Create(ctx, key, email, &sid)
	if err != nil {
		return nil, err
	}
	if sendErr := s.emailService.SendLicenseKey(email, key); sendErr != nil {
		log.Printf("warn: failed to send license email to %s: %v", email, sendErr)
	}
	return lk, nil
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
