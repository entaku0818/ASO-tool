package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/notification"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

//go:embed migrations/007_create_asc_credentials.up.sql
var migration007 string

//go:embed migrations/008_create_analytics.up.sql
var migration008 string

//go:embed migrations/009_create_store_rankings.up.sql
var migration009 string

func main() {
	ctx := context.Background()
	startTime := time.Now()

	// Initialize Slack notifier early so we can send failure notifications
	slackWebhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	slackNotifier := notification.NewSlackNotifier(slackWebhookURL)
	if slackNotifier.IsConfigured() {
		log.Println("Slack notifications enabled")
	}

	// Determine job type early for error notifications
	job := "all"
	if len(os.Args) > 1 {
		job = os.Args[1]
	}

	// Helper function to send failure notification and exit
	sendFailureAndExit := func(errorMsg string) {
		log.Printf("FATAL: %s", errorMsg)
		if slackNotifier.IsConfigured() {
			result := &notification.BatchResult{
				StartTime: startTime,
				EndTime:   time.Now(),
				JobType:   job,
				Success:   false,
				Errors:    []string{errorMsg},
			}
			if err := slackNotifier.SendBatchResult(result); err != nil {
				log.Printf("Failed to send Slack notification: %v", err)
			} else {
				log.Println("Slack failure notification sent")
			}
		}
		os.Exit(1)
	}

	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		sendFailureAndExit("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		sendFailureAndExit(fmt.Sprintf("Unable to connect to database: %v", err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		sendFailureAndExit(fmt.Sprintf("Unable to ping database: %v", err))
	}
	log.Println("Connected to database")

	// Initialize repositories and services
	keywordRepo := repository.NewKeywordRepository(pool)
	rankingRepo := repository.NewRankingRepository(pool)
	reviewRepo := repository.NewReviewRepository(pool)
	appRepo := repository.NewAppRepository(pool)
	trackedKeywordRepo := repository.NewTrackedKeywordRepository(pool)

	scraperService := service.NewScraperServiceWithTracking(
		keywordRepo,
		rankingRepo,
		reviewRepo,
		appRepo,
		trackedKeywordRepo,
	)

	storeRankingRepo := repository.NewStoreRankingRepository(pool)
	appRankingService := service.NewAppRankingService(storeRankingRepo)

	rankingChangeService := service.NewRankingChangeService(
		rankingRepo,
		keywordRepo,
		appRepo,
	)

	// Initialize batch result for notification
	result := &notification.BatchResult{
		StartTime: startTime,
		Success:   true,
		Errors:    []string{},
	}

	result.JobType = job

	switch job {
	case "migrate":
		runMigrations(ctx, pool)
	case "seed":
		runSeed(ctx, pool)
	case "rankings":
		apps, keywords, errs := runRankingsUpdate(ctx, scraperService)
		result.AppsProcessed = apps
		result.KeywordsUpdated = keywords
		result.Errors = append(result.Errors, errs...)
		// Detect ranking changes after update
		changes, _ := rankingChangeService.DetectChanges(ctx)
		result.RankingChanges = changes
	case "store-rankings":
		saved, errs := runStoreRankingsFetch(ctx, appRankingService)
		result.KeywordsUpdated = saved
		for _, e := range errs {
			result.Errors = append(result.Errors, e.Error())
		}
	case "tracked-keywords":
		tracked, errs := runTrackedKeywordsUpdate(ctx, scraperService)
		result.TrackedKeywords = tracked
		result.Errors = append(result.Errors, errs...)
	case "all":
		apps, keywords, errs := runRankingsUpdate(ctx, scraperService)
		result.AppsProcessed = apps
		result.KeywordsUpdated = keywords
		result.Errors = append(result.Errors, errs...)
		// Detect ranking changes after update
		changes, _ := rankingChangeService.DetectChanges(ctx)
		result.RankingChanges = changes

		tracked, trackedErrs := runTrackedKeywordsUpdate(ctx, scraperService)
		result.TrackedKeywords = tracked
		result.Errors = append(result.Errors, trackedErrs...)

		saved, storeErrs := runStoreRankingsFetch(ctx, appRankingService)
		result.KeywordsUpdated += saved
		for _, e := range storeErrs {
			result.Errors = append(result.Errors, e.Error())
		}
	default:
		sendFailureAndExit(fmt.Sprintf("Unknown job: %s. Use: migrate, seed, rankings, store-rankings, tracked-keywords, or all", job))
	}

	result.EndTime = time.Now()
	if len(result.Errors) > 0 {
		result.Success = false
	}

	// Send Slack notification
	if slackNotifier.IsConfigured() && (job == "all" || job == "rankings" || job == "tracked-keywords") {
		if err := slackNotifier.SendBatchResult(result); err != nil {
			log.Printf("Failed to send Slack notification: %v", err)
		} else {
			log.Println("Slack notification sent successfully")
		}
	}

	log.Println("Batch job completed successfully")
}

func runRankingsUpdate(ctx context.Context, s *service.ScraperService) (int, int, []string) {
	log.Println("Starting rankings update...")
	var errors []string

	results, err := s.TriggerAllUpdates(ctx)
	if err != nil {
		log.Printf("Error updating rankings: %v", err)
		errors = append(errors, fmt.Sprintf("Rankings update error: %v", err))
		return 0, 0, errors
	}

	total := 0
	failedApps := 0
	for appID, count := range results {
		if count >= 0 {
			total += count
			log.Printf("App %s: %d keywords updated", appID, count)
		} else {
			log.Printf("App %s: failed", appID)
			errors = append(errors, fmt.Sprintf("App %s: failed to update", appID))
			failedApps++
		}
	}

	fmt.Printf("Rankings update complete: %d total keywords updated across %d apps\n", total, len(results))
	return len(results) - failedApps, total, errors
}

func runTrackedKeywordsUpdate(ctx context.Context, s *service.ScraperService) (int, []string) {
	log.Println("Starting tracked keywords update...")
	var errors []string

	results, err := s.UpdateTrackedKeywordResults(ctx)
	if err != nil {
		log.Printf("Error updating tracked keywords: %v", err)
		errors = append(errors, fmt.Sprintf("Tracked keywords error: %v", err))
		return 0, errors
	}

	total := 0
	for keywordID, count := range results {
		if count >= 0 {
			total += count
			log.Printf("Tracked keyword %s: %d results", keywordID, count)
		} else {
			log.Printf("Tracked keyword %s: failed", keywordID)
			errors = append(errors, fmt.Sprintf("Tracked keyword %s: failed", keywordID))
		}
	}

	fmt.Printf("Tracked keywords update complete: %d total results across %d keywords\n", total, len(results))
	return total, errors
}

func runStoreRankingsFetch(ctx context.Context, svc *service.AppRankingService) (int, []error) {
	log.Println("Starting store rankings fetch...")
	saved, errs := svc.FetchAndSaveAll(ctx)
	for _, e := range errs {
		log.Printf("Store rankings error: %v", e)
	}
	log.Printf("Store rankings fetch complete: %d entries saved", saved)
	return saved, errs
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) {
	log.Println("Starting migrations...")

	migrations := []struct {
		name string
		sql  string
	}{
		{"007_create_asc_credentials", migration007},
		{"008_create_analytics", migration008},
		{"009_create_store_rankings", migration009},
	}

	for _, m := range migrations {
		log.Printf("Running migration: %s", m.name)
		_, err := pool.Exec(ctx, m.sql)
		if err != nil {
			log.Printf("Migration %s failed: %v", m.name, err)
			// Continue with other migrations as they use IF NOT EXISTS
		} else {
			log.Printf("Migration %s completed", m.name)
		}
	}

	log.Println("All migrations completed")
}

func runSeed(ctx context.Context, pool *pgxpool.Pool) {
	log.Println("Starting seed...")

	// Check if admin user already exists
	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE email = $1", "admin@example.com").Scan(&count)
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		return
	}

	if count > 0 {
		log.Println("Admin user already exists, skipping seed")
		return
	}

	// Create admin user
	password := "admin123" // Default password - should be changed after first login
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return
	}

	id := uuid.New().String()
	_, err = pool.Exec(ctx,
		"INSERT INTO users (id, email, password_hash, name, is_admin) VALUES ($1, $2, $3, $4, $5)",
		id, "admin@example.com", string(hashedPassword), "Admin", true,
	)
	if err != nil {
		log.Printf("Error creating admin user: %v", err)
		return
	}

	log.Println("Admin user created successfully")
	log.Println("Email: admin@example.com")
	log.Println("Password: admin123")
	log.Println("Please change the password after first login!")
}
