package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/model"
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

func main() {
	ctx := context.Background()
	startTime := time.Now()

	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
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

	rankingChangeService := service.NewRankingChangeService(
		keywordRepo,
		rankingRepo,
		appRepo,
	)

	// Initialize Slack notifier if webhook URL is set
	slackWebhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	var slackNotifier *notification.SlackNotifier
	if slackWebhookURL != "" {
		slackNotifier = notification.NewSlackNotifier(slackWebhookURL)
		log.Println("Slack notifications enabled")
	} else {
		log.Println("SLACK_WEBHOOK_URL not set, notifications disabled")
	}

	// Determine which job to run
	job := "all"
	if len(os.Args) > 1 {
		job = os.Args[1]
	}

	// Initialize batch result for notification
	result := &model.BatchResult{
		JobName:   job,
		StartedAt: startTime,
		Success:   true,
	}

	switch job {
	case "migrate":
		runMigrations(ctx, pool)
	case "seed":
		runSeed(ctx, pool)
	case "rankings":
		apps, keywords, errs := runRankingsUpdateWithResult(ctx, scraperService)
		result.AppsProcessed = apps
		result.KeywordsUpdated = keywords
		result.Errors = errs
		if len(errs) > 0 {
			result.Success = false
		}

		// Detect ranking changes
		changes, err := rankingChangeService.DetectChangesForAllApps(ctx)
		if err != nil {
			log.Printf("Error detecting ranking changes: %v", err)
		} else {
			result.Changes = changes
		}
	case "tracked-keywords":
		keywords, errs := runTrackedKeywordsUpdateWithResult(ctx, scraperService)
		result.KeywordsUpdated = keywords
		result.Errors = errs
		if len(errs) > 0 {
			result.Success = false
		}
	case "all":
		apps, keywords1, errs1 := runRankingsUpdateWithResult(ctx, scraperService)
		keywords2, errs2 := runTrackedKeywordsUpdateWithResult(ctx, scraperService)
		result.AppsProcessed = apps
		result.KeywordsUpdated = keywords1 + keywords2
		result.Errors = append(errs1, errs2...)
		if len(result.Errors) > 0 {
			result.Success = false
		}

		// Detect ranking changes
		changes, err := rankingChangeService.DetectChangesForAllApps(ctx)
		if err != nil {
			log.Printf("Error detecting ranking changes: %v", err)
		} else {
			result.Changes = changes
		}
	default:
		log.Fatalf("Unknown job: %s. Use: migrate, seed, rankings, tracked-keywords, or all", job)
	}

	// Complete the result
	result.CompletedAt = time.Now()
	result.Duration = result.CompletedAt.Sub(result.StartedAt)

	// Send Slack notification
	if slackNotifier != nil && job != "migrate" && job != "seed" {
		if err := slackNotifier.SendBatchResult(ctx, result); err != nil {
			log.Printf("Failed to send Slack notification: %v", err)
		} else {
			log.Println("Slack notification sent successfully")
		}
	}

	log.Println("Batch job completed successfully")
}

func runRankingsUpdate(ctx context.Context, s *service.ScraperService) {
	apps, keywords, _ := runRankingsUpdateWithResult(ctx, s)
	fmt.Printf("Rankings update complete: %d total keywords updated across %d apps\n", keywords, apps)
}

func runRankingsUpdateWithResult(ctx context.Context, s *service.ScraperService) (apps int, keywords int, errors []string) {
	log.Println("Starting rankings update...")

	results, err := s.TriggerAllUpdates(ctx)
	if err != nil {
		log.Printf("Error updating rankings: %v", err)
		errors = append(errors, fmt.Sprintf("Failed to trigger updates: %v", err))
		return 0, 0, errors
	}

	apps = len(results)
	for appID, count := range results {
		if count >= 0 {
			keywords += count
			log.Printf("App %s: %d keywords updated", appID, count)
		} else {
			log.Printf("App %s: failed", appID)
			errors = append(errors, fmt.Sprintf("App %s: update failed", appID))
		}
	}

	log.Printf("Rankings update complete: %d total keywords updated across %d apps", keywords, apps)
	return apps, keywords, errors
}

func runTrackedKeywordsUpdate(ctx context.Context, s *service.ScraperService) {
	keywords, _ := runTrackedKeywordsUpdateWithResult(ctx, s)
	fmt.Printf("Tracked keywords update complete: %d keywords\n", keywords)
}

func runTrackedKeywordsUpdateWithResult(ctx context.Context, s *service.ScraperService) (keywords int, errors []string) {
	log.Println("Starting tracked keywords update...")

	results, err := s.UpdateTrackedKeywordResults(ctx)
	if err != nil {
		log.Printf("Error updating tracked keywords: %v", err)
		errors = append(errors, fmt.Sprintf("Failed to update tracked keywords: %v", err))
		return 0, errors
	}

	total := 0
	for keywordID, count := range results {
		if count >= 0 {
			total += count
			log.Printf("Tracked keyword %s: %d results", keywordID, count)
		} else {
			log.Printf("Tracked keyword %s: failed", keywordID)
			errors = append(errors, fmt.Sprintf("Tracked keyword %s: update failed", keywordID))
		}
	}

	keywords = len(results)
	log.Printf("Tracked keywords update complete: %d total results across %d keywords", total, keywords)
	return keywords, errors
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) {
	log.Println("Starting migrations...")

	migrations := []struct {
		name string
		sql  string
	}{
		{"007_create_asc_credentials", migration007},
		{"008_create_analytics", migration008},
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
