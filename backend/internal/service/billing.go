package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/entaku0818/aso-compass/backend/internal/repository"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/webhook"
)

type BillingService struct {
	userRepo        *repository.UserRepository
	licenseService  *LicenseService
	webhookSecret   string
	priceIDMonthly  string
	priceIDYearly   string
	frontendBaseURL string
}

func NewBillingService(userRepo *repository.UserRepository) *BillingService {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		log.Println("warn: STRIPE_SECRET_KEY is not set — billing endpoints will return 503")
	}

	svc := &BillingService{
		userRepo:        userRepo,
		webhookSecret:   os.Getenv("STRIPE_WEBHOOK_SECRET"),
		priceIDMonthly:  os.Getenv("STRIPE_PRICE_MONTHLY"),
		priceIDYearly:   os.Getenv("STRIPE_PRICE_YEARLY"),
		frontendBaseURL: getEnvOrDefault("FRONTEND_BASE_URL", "http://localhost:3000"),
	}

	if svc.webhookSecret == "" {
		if isRunningOnCloudRun() {
			log.Fatal("fatal: STRIPE_WEBHOOK_SECRET is not set in production — refusing to start with an unverifiable billing webhook")
		}
		log.Println("warn: STRIPE_WEBHOOK_SECRET is not set — webhook signature verification will fail")
	}
	if svc.priceIDMonthly == "" || svc.priceIDYearly == "" {
		log.Println("warn: STRIPE_PRICE_MONTHLY or STRIPE_PRICE_YEARLY is not set — checkout will fail")
	}

	return svc
}

// SetLicenseService wires the license service for webhook handling.
func (s *BillingService) SetLicenseService(ls *LicenseService) {
	s.licenseService = ls
}

// CreateCheckoutSession creates a Stripe Checkout session and returns the URL.
func (s *BillingService) CreateCheckoutSession(ctx context.Context, userID, planType string) (string, error) {
	if stripe.Key == "" {
		return "", model.ErrStripeKeyRequired
	}

	priceID := s.priceIDMonthly
	if planType == "yearly" {
		priceID = s.priceIDYearly
	}
	if priceID == "" {
		return "", fmt.Errorf("price ID for plan type %q is not configured", planType)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	// Create Stripe customer if one doesn't exist yet
	stripeCustomerID := user.StripeCustomerID
	if stripeCustomerID == "" {
		c, err := customer.New(&stripe.CustomerParams{
			Email: stripe.String(user.Email),
			Metadata: map[string]string{
				"user_id": userID,
			},
		})
		if err != nil {
			return "", fmt.Errorf("create stripe customer: %w", err)
		}
		stripeCustomerID = c.ID
		if err := s.userRepo.SetStripeCustomerID(ctx, userID, stripeCustomerID); err != nil {
			log.Printf("warn: failed to save stripe_customer_id for user %s: %v", userID, err)
		}
	}

	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(stripeCustomerID),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(s.frontendBaseURL + "/billing/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(s.frontendBaseURL + "/billing/cancel"),
	}

	sess, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("create checkout session: %w", err)
	}

	return sess.URL, nil
}

// HandleWebhook verifies the Stripe signature and processes the event.
// When STRIPE_WEBHOOK_SECRET is empty (local dev), signature verification is skipped
// and the raw JSON payload is parsed directly — never use in production.
func (s *BillingService) HandleWebhook(ctx context.Context, body io.Reader, signature string) error {
	payload, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	var event stripe.Event
	if s.webhookSecret == "" {
		log.Println("warn: skipping webhook signature verification (STRIPE_WEBHOOK_SECRET not set)")
		if err := json.Unmarshal(payload, &event); err != nil {
			return fmt.Errorf("parse webhook payload: %w", err)
		}
	} else {
		event, err = webhook.ConstructEventWithOptions(payload, signature, s.webhookSecret, webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		})
		if err != nil {
			return fmt.Errorf("webhook signature verification failed: %w", err)
		}
	}

	switch event.Type {
	case "checkout.session.completed":
		return s.handleCheckoutCompleted(ctx, event)
	case "invoice.payment_succeeded":
		return s.handleInvoicePaymentSucceeded(ctx, event)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(ctx, event)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, event)
	default:
		// Unhandled event types are ignored
	}

	return nil
}

func (s *BillingService) handleCheckoutCompleted(ctx context.Context, event stripe.Event) error {
	var sess stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
		return fmt.Errorf("parse checkout.session: %w", err)
	}

	// License key subscription purchase
	subMeta := map[string]string{}
	if sess.Subscription != nil {
		subMeta = sess.Subscription.Metadata
	}
	if sess.Metadata["purchase_type"] == "license" || subMeta["purchase_type"] == "license" {
		email := sess.CustomerEmail
		if email == "" && sess.CustomerDetails != nil {
			email = sess.CustomerDetails.Email
		}
		if email == "" {
			return fmt.Errorf("no email in license checkout session %s", sess.ID)
		}
		if s.licenseService == nil {
			return fmt.Errorf("license service not wired")
		}
		subID := ""
		if sess.Subscription != nil {
			subID = sess.Subscription.ID
		}
		lk, err := s.licenseService.GenerateAndSendLicense(ctx, email, sess.ID, subID)
		if err != nil {
			return fmt.Errorf("generate license for %s: %w", email, err)
		}
		log.Printf("info: license key %s generated and sent to %s (session: %s)", lk.Key, email, sess.ID)
		return nil
	}

	// Subscription upgrade (existing flow)
	user, err := s.userRepo.GetByStripeCustomerID(ctx, sess.Customer.ID)
	if err != nil {
		return fmt.Errorf("get user by stripe customer %s: %w", sess.Customer.ID, err)
	}

	subscriptionID := ""
	if sess.Subscription != nil {
		subscriptionID = sess.Subscription.ID
	}

	if err := s.userRepo.UpdatePlan(ctx, user.ID, "pro", sess.Customer.ID, subscriptionID, nil); err != nil {
		return fmt.Errorf("update plan to pro: %w", err)
	}

	log.Printf("info: user %s upgraded to pro (subscription: %s)", user.ID, subscriptionID)
	return nil
}

// handleInvoicePaymentSucceeded extends license expiry by 1 year on renewal.
func (s *BillingService) handleInvoicePaymentSucceeded(ctx context.Context, event stripe.Event) error {
	var inv stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &inv); err != nil {
		return fmt.Errorf("parse invoice: %w", err)
	}
	// Only handle license subscription renewals (not first payment — that is handled by checkout.session.completed)
	if inv.BillingReason != stripe.InvoiceBillingReasonSubscriptionCycle {
		return nil
	}
	if inv.Subscription == nil {
		return nil
	}
	if inv.Subscription.Metadata["purchase_type"] != "license" {
		return nil
	}
	if s.licenseService == nil {
		return fmt.Errorf("license service not wired")
	}
	if err := s.licenseService.ExtendLicense(ctx, inv.Subscription.ID); err != nil {
		return fmt.Errorf("extend license for subscription %s: %w", inv.Subscription.ID, err)
	}
	log.Printf("info: license extended for subscription %s", inv.Subscription.ID)
	return nil
}

func (s *BillingService) handleSubscriptionUpdated(ctx context.Context, event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return fmt.Errorf("parse subscription: %w", err)
	}

	user, err := s.userRepo.GetByStripeCustomerID(ctx, sub.Customer.ID)
	if err != nil {
		return fmt.Errorf("get user by stripe customer %s: %w", sub.Customer.ID, err)
	}

	plan := "free"
	if sub.Status == stripe.SubscriptionStatusActive || sub.Status == stripe.SubscriptionStatusTrialing {
		plan = "pro"
	}

	if err := s.userRepo.UpdatePlan(ctx, user.ID, plan, sub.Customer.ID, sub.ID, nil); err != nil {
		return fmt.Errorf("update plan on subscription update: %w", err)
	}

	log.Printf("info: user %s plan updated to %s (subscription: %s, status: %s)", user.ID, plan, sub.ID, sub.Status)
	return nil
}

func (s *BillingService) handleSubscriptionDeleted(ctx context.Context, event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return fmt.Errorf("parse subscription: %w", err)
	}

	// License subscription cancelled — deactivate license
	if sub.Metadata["purchase_type"] == "license" {
		if s.licenseService != nil {
			if err := s.licenseService.DeactivateBySubscription(ctx, sub.ID); err != nil {
				log.Printf("warn: failed to deactivate license for subscription %s: %v", sub.ID, err)
			}
		}
		log.Printf("info: license deactivated for subscription %s", sub.ID)
		return nil
	}

	user, err := s.userRepo.GetByStripeCustomerID(ctx, sub.Customer.ID)
	if err != nil {
		return fmt.Errorf("get user by stripe customer %s: %w", sub.Customer.ID, err)
	}

	if err := s.userRepo.UpdatePlan(ctx, user.ID, "free", sub.Customer.ID, "", nil); err != nil {
		return fmt.Errorf("update plan to free: %w", err)
	}

	log.Printf("info: user %s downgraded to free (subscription cancelled: %s)", user.ID, sub.ID)
	return nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// isRunningOnCloudRun reports whether the process is running as a Cloud Run
// service. K_SERVICE is injected automatically by the Cloud Run runtime and
// is never set in local development, so it doubles as a production signal
// without requiring any extra configuration.
func isRunningOnCloudRun() bool {
	return os.Getenv("K_SERVICE") != ""
}
