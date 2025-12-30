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
	"github.com/entaku0818/aso-tool/backend/internal/repository"
	"github.com/entaku0818/aso-tool/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://aso:aso_password@localhost:5432/aso_tool?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	log.Println("Connected to database")

	// Initialize layers
	appRepo := repository.NewAppRepository(pool)
	appService := service.NewAppService(appRepo)
	appHandler := handler.NewAppHandler(appService)

	keywordRepo := repository.NewKeywordRepository(pool)
	keywordService := service.NewKeywordService(keywordRepo)
	keywordHandler := handler.NewKeywordHandler(keywordService)

	rankingRepo := repository.NewRankingRepository(pool)
	rankingService := service.NewRankingService(rankingRepo)
	rankingHandler := handler.NewRankingHandler(rankingService)

	reviewRepo := repository.NewReviewRepository(pool)
	reviewService := service.NewReviewService(reviewRepo)
	reviewHandler := handler.NewReviewHandler(reviewService)

	scraperService := service.NewScraperService(keywordRepo, rankingRepo, reviewRepo, appRepo)
	scraperHandler := handler.NewScraperHandler(scraperService)

	// Router setup
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
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
