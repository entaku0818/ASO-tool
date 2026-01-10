package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/entaku0818/aso-tool/backend/internal/handler"
	appMiddleware "github.com/entaku0818/aso-tool/backend/internal/middleware"
	"github.com/entaku0818/aso-tool/backend/internal/repository"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func connectWithRetry(ctx context.Context, dbURL string, maxRetries int, retryDelay time.Duration) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	for i := 0; i < maxRetries; i++ {
		pool, err = pgxpool.New(ctx, dbURL)
		if err != nil {
			log.Printf("Attempt %d/%d: Failed to create pool: %v", i+1, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		// Try to ping the database
		if pingErr := pool.Ping(ctx); pingErr != nil {
			log.Printf("Attempt %d/%d: Failed to ping database: %v", i+1, maxRetries, pingErr)
			pool.Close()
			time.Sleep(retryDelay)
			continue
		}

		log.Println("Connected to database")
		return pool, nil
	}

	return nil, err
}

func main() {
	ctx := context.Background()

	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://aso:aso_password@localhost:5432/aso_tool?sslmode=disable"
	}

	// Connect with retry (20 retries, 10 second delay = 200 seconds max)
	// Allows time for Cloud SQL cold start
	pool, err := connectWithRetry(ctx, dbURL, 20, 10*time.Second)
	if err != nil {
		log.Fatalf("Unable to connect to database after retries: %v", err)
	}
	defer pool.Close()

	// Initialize layers
	// Auth
	userRepo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo)
	authMiddleware := appMiddleware.NewAuthMiddleware(authService)
	authHandler := handler.NewAuthHandler(authService)

	appRepo := repository.NewAppRepository(pool)
	appService := service.NewAppService(appRepo)
	appHandler := handler.NewAppHandler(appService)

	trackedKeywordRepo := repository.NewTrackedKeywordRepository(pool)
	trackedKeywordService := service.NewTrackedKeywordService(trackedKeywordRepo)
	trackedKeywordHandler := handler.NewTrackedKeywordHandler(trackedKeywordService)

	keywordRepo := repository.NewKeywordRepository(pool)
	keywordService := service.NewKeywordServiceWithTracking(keywordRepo, appRepo, trackedKeywordRepo)
	keywordHandler := handler.NewKeywordHandler(keywordService)

	rankingRepo := repository.NewRankingRepository(pool)
	rankingService := service.NewRankingService(rankingRepo)
	rankingHandler := handler.NewRankingHandler(rankingService)

	reviewRepo := repository.NewReviewRepository(pool)
	reviewService := service.NewReviewService(reviewRepo)
	reviewHandler := handler.NewReviewHandler(reviewService)

	scraperService := service.NewScraperServiceWithTracking(keywordRepo, rankingRepo, reviewRepo, appRepo, trackedKeywordRepo)
	scraperHandler := handler.NewScraperHandler(scraperService)

	competitorRepo := repository.NewCompetitorRepository(pool)
	competitorService := service.NewCompetitorService(competitorRepo, keywordRepo)
	competitorHandler := handler.NewCompetitorHandler(competitorService)

	screenshotRepo := repository.NewScreenshotRepository(pool)
	screenshotService := service.NewScreenshotService(screenshotRepo)
	screenshotHandler := handler.NewScreenshotHandler(screenshotService)

	// Analytics
	ascCredentialsRepo := repository.NewASCCredentialsRepository(pool)
	analyticsRepo := repository.NewAnalyticsRepository(pool)
	appVersionRepo := repository.NewAppVersionRepository(pool)
	analyticsService := service.NewAnalyticsService(
		ascCredentialsRepo, analyticsRepo, appVersionRepo,
		rankingRepo, reviewRepo, appRepo,
	)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	// Router setup
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check with DB status
	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		dbStatus := "ok"
		if err := pool.Ping(req.Context()); err != nil {
			dbStatus = "error: " + err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","database":"` + dbStatus + `"}`))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Public routes
		r.Post("/auth/login", authHandler.Login)

		// Protected routes (require authentication)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Verify)

			// Auth routes
			r.Get("/auth/me", authHandler.Me)

			// Admin-only routes
			r.With(authMiddleware.RequireAdmin).Post("/users", authHandler.CreateUser)

			r.Route("/apps", func(r chi.Router) {
				r.Get("/", appHandler.List)
				r.Post("/", appHandler.Create)
				r.Get("/{id}", appHandler.Get)
				r.Put("/{id}", appHandler.Update)
				r.Delete("/{id}", appHandler.Delete)

				// Nested routes for app-specific resources
				r.Route("/{appID}/keywords", func(r chi.Router) {
					r.Get("/", keywordHandler.ListByApp)
					r.Post("/", keywordHandler.Create)
					r.Get("/{keywordID}", keywordHandler.Get)
					r.Delete("/{keywordID}", keywordHandler.Delete)
					r.Get("/{keywordID}/rankings", rankingHandler.ListByKeyword)
					r.Get("/{keywordID}/rankings/latest", rankingHandler.GetLatest)
				})

				r.Route("/{appID}/rankings", func(r chi.Router) {
					r.Get("/", rankingHandler.ListByApp)
				})

				r.Route("/{appID}/reviews", func(r chi.Router) {
					r.Get("/", reviewHandler.ListByApp)
					r.Post("/", reviewHandler.Create)
					r.Get("/stats", reviewHandler.GetStats)
					r.Get("/{reviewID}", reviewHandler.Get)
					r.Delete("/{reviewID}", reviewHandler.Delete)
				})

				r.Route("/{appID}/competitors", func(r chi.Router) {
					r.Get("/", competitorHandler.List)
					r.Post("/", competitorHandler.Create)
					r.Get("/{competitorID}", competitorHandler.Get)
					r.Delete("/{competitorID}", competitorHandler.Delete)
					r.Post("/update-rankings", competitorHandler.UpdateRankings)
				})

				r.Route("/{appID}/screenshots", func(r chi.Router) {
					r.Get("/", screenshotHandler.List)
					r.Post("/", screenshotHandler.Create)
					r.Post("/reorder", screenshotHandler.Reorder)
					r.Get("/{screenshotID}", screenshotHandler.Get)
					r.Put("/{screenshotID}", screenshotHandler.Update)
					r.Delete("/{screenshotID}", screenshotHandler.Delete)
				})

				// App Store Connect credentials
				r.Route("/{appID}/asc-credentials", func(r chi.Router) {
					r.Post("/", analyticsHandler.SetCredentials)
					r.Get("/", analyticsHandler.GetCredentials)
					r.Delete("/", analyticsHandler.DeleteCredentials)
					r.Post("/validate", analyticsHandler.ValidateCredentials)
				})

				// Analytics
				r.Route("/{appID}/analytics", func(r chi.Router) {
					r.Get("/", analyticsHandler.GetAnalytics)
					r.Get("/correlation", analyticsHandler.GetAnalyticsWithCorrelation)
					r.Get("/summary", analyticsHandler.GetSummary)
					r.Post("/fetch", analyticsHandler.TriggerFetch)
				})

				r.Get("/{appID}/keywords/{keywordID}/comparison", competitorHandler.GetComparison)
			})

			// Rankings endpoint for creating new rankings
			r.Post("/rankings", rankingHandler.Create)

			// Scraper endpoints
			r.Route("/scraper", func(r chi.Router) {
				r.Get("/app-info", scraperHandler.FetchAppInfo)
				r.Get("/search", scraperHandler.SearchApps)
				r.Post("/trigger", scraperHandler.TriggerAllUpdates)
			})

			// App-specific scraper actions
			r.Post("/apps/{appID}/scrape/rankings", scraperHandler.UpdateRankings)
			r.Post("/apps/{appID}/scrape/reviews", scraperHandler.FetchReviews)

			// Tracked keywords endpoints
			r.Route("/tracked-keywords", func(r chi.Router) {
				r.Get("/", trackedKeywordHandler.List)
				r.Post("/trigger", scraperHandler.UpdateTrackedKeywords)
				r.Get("/{id}", trackedKeywordHandler.Get)
				r.Get("/{id}/results", trackedKeywordHandler.GetSearchResults)
				r.Post("/{id}/trigger", scraperHandler.UpdateSingleTrackedKeyword)
			})
		})
	})

	// Server setup
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
