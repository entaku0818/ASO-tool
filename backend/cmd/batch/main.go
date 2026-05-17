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
	"github.com/entaku0818/aso-tool/backend/internal/scraper"
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

//go:embed migrations/010_search_ads.up.sql
var migration010 string

//go:embed migrations/015_public_keyword_cache.up.sql
var migration015 string

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

	keywordCacheRepo := repository.NewPublicKeywordCacheRepository(pool)

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
	case "keyword-cache":
		cached, errs := runKeywordCacheUpdate(ctx, keywordCacheRepo)
		result.KeywordsUpdated = cached
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

		cached, cacheErrs := runKeywordCacheUpdate(ctx, keywordCacheRepo)
		result.KeywordsUpdated += cached
		result.Errors = append(result.Errors, cacheErrs...)
	default:
		sendFailureAndExit(fmt.Sprintf("Unknown job: %s. Use: migrate, seed, rankings, store-rankings, tracked-keywords, keyword-cache, or all", job))
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
		{"010_search_ads", migration010},
		{"015_public_keyword_cache", migration015},
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

// seedApps are popular apps used to generate keyword suggestions from Search Ads API.
var seedApps = []struct {
	adamID int64
	genre  string
}{
	{1130383480, "social"},    // LINE
	{835599320, "social"},     // TikTok
	{389801252, "video"},      // YouTube
	{525895022, "navigation"}, // Google Maps
	{310633997, "social"},     // WhatsApp
	{284882215, "social"},     // Facebook
	{544007664, "photo"},      // CamScanner
	{1480637954, "shopping"},  // Shein
}

var seedCountries = []string{"jp", "us", "gb", "de", "fr", "kr", "tw", "au"}

func runKeywordCacheUpdate(ctx context.Context, cacheRepo *repository.PublicKeywordCacheRepository) (int, []string) {
	log.Println("Starting keyword cache update...")
	var errors []string

	clientID := os.Getenv("ADMIN_ASA_CLIENT_ID")
	teamID := os.Getenv("ADMIN_ASA_TEAM_ID")
	keyID := os.Getenv("ADMIN_ASA_KEY_ID")
	privateKeyPEM := os.Getenv("ADMIN_ASA_PRIVATE_KEY")
	orgID := os.Getenv("ADMIN_ASA_ORG_ID")

	if clientID == "" || teamID == "" || keyID == "" || privateKeyPEM == "" {
		log.Println("Keyword cache: ADMIN_ASA_* env vars not set, skipping")
		return 0, nil
	}

	client, err := scraper.NewSearchAdsClient(clientID, teamID, keyID, []byte(privateKeyPEM), orgID)
	if err != nil {
		errors = append(errors, fmt.Sprintf("keyword cache: failed to create Search Ads client: %v", err))
		return 0, errors
	}

	total := 0
	for _, country := range seedCountries {
		for _, app := range seedApps {
			results, err := client.GetKeywordPopularity(ctx, app.adamID, nil, 100)
			if err != nil {
				log.Printf("Keyword cache: country=%s adamID=%d error: %v", country, app.adamID, err)
				errors = append(errors, fmt.Sprintf("keyword cache %s/%d: %v", country, app.adamID, err))
				continue
			}

			entries := make([]repository.PublicKeywordCacheEntry, 0, len(results))
			for _, kp := range results {
				entries = append(entries, repository.PublicKeywordCacheEntry{
					Keyword:    kp.Text,
					Country:    country,
					Genre:      app.genre,
					Popularity: kp.PopularityScore * 20, // 0–5 → 0–100
				})
			}

			if err := cacheRepo.UpsertMany(ctx, entries); err != nil {
				log.Printf("Keyword cache upsert error: %v", err)
				errors = append(errors, fmt.Sprintf("keyword cache upsert %s/%d: %v", country, app.adamID, err))
				continue
			}
			total += len(entries)
			log.Printf("Keyword cache: country=%s genre=%s adamID=%d saved %d keywords", country, app.genre, app.adamID, len(entries))
		}
	}

	log.Printf("Keyword cache update complete: %d entries saved", total)
	return total, errors
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
