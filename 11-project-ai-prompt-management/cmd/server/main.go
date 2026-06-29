package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/timurdian/prompt-management/internal/config"
	"github.com/timurdian/prompt-management/internal/entity"
	"github.com/timurdian/prompt-management/internal/handler"
	"github.com/timurdian/prompt-management/internal/middleware"
	"github.com/timurdian/prompt-management/internal/repository"
	"github.com/timurdian/prompt-management/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load Dotenv
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 2. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Fatal: failed to load configuration: %v", err)
	}

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 3. Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fatal: failed to connect to database: %v", err)
	}
	log.Println("PostgreSQL connection established successfully")

	// Auto-Migrate Models
	err = db.AutoMigrate(
		&entity.Workspace{},
		&entity.WorkspaceMember{},
		&entity.Prompt{},
		&entity.PromptVersion{},
		&entity.ApiKey{},
		&entity.AnalyticsLog{},
	)
	if err != nil {
		log.Fatalf("Fatal: failed to run database auto-migrations: %v", err)
	}
	log.Println("PostgreSQL auto-migrations completed successfully")

	// 4. Connect to Redis Client
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer pingCancel()
	if err := rdb.Ping(pingCtx).Err(); err != nil {
		log.Fatalf("Fatal: failed to connect to Redis: %v", err)
	}
	log.Println("Redis connection established successfully")

	// 5. Initialize Services & Repositories (DI)
	promptRepo := repository.NewPromptRepository(db)
	analyticsSvc := service.NewAnalyticsService(promptRepo)

	// Start Background Analytics Daemon
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	analyticsSvc.StartWorker(ctx)

	promptSvc := service.NewPromptService(promptRepo, rdb, analyticsSvc)

	workspaceHandler := handler.NewWorkspaceHandler(promptSvc)
	promptHandler := handler.NewPromptHandler(promptSvc)

	// 6. Setup Middlewares
	jwtMiddleware, err := middleware.NewJWTMiddleware(cfg.RSAPublicKeyPath)
	if err != nil {
		log.Fatalf("Fatal: failed to initialize JWT Middleware: %v", err)
	}

	apiKeyMiddleware := middleware.NewAPIKeyMiddleware(promptRepo, rdb)

	// 7. Setup Router
	r := gin.Default()
	r.Use(gin.Recovery())

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// User / Admin dashboard endpoints (JWT Protected - Offline RS256 Verification)
	adminRoutes := r.Group("/api/v1")
	adminRoutes.Use(jwtMiddleware.AuthRequired())
	{
		adminRoutes.POST("/workspaces", workspaceHandler.CreateWorkspace)
		adminRoutes.POST("/workspaces/:id/api-keys", workspaceHandler.CreateApiKey)
		adminRoutes.GET("/workspaces/:id/api-keys", workspaceHandler.GetApiKeys)
		adminRoutes.DELETE("/workspaces/:id/api-keys/:key_id", workspaceHandler.RevokeApiKey)
		adminRoutes.GET("/workspaces/:id/analytics", workspaceHandler.GetWorkspaceAnalytics)

		adminRoutes.POST("/prompts", promptHandler.CreatePrompt)
		adminRoutes.GET("/prompts/:id", promptHandler.GetPrompt)
		adminRoutes.GET("/workspaces/:id/prompts", promptHandler.GetWorkspacePrompts)
		adminRoutes.POST("/prompts/:id/versions", promptHandler.CreateVersion)
		adminRoutes.PUT("/prompts/:id/versions/:version_number/activate", promptHandler.ActivateVersion)
	}

	// Server-to-Server Client integration endpoints (API Key Protected - Redis Lookup)
	clientRoutes := r.Group("/api/v1/client")
	clientRoutes.Use(apiKeyMiddleware.APIKeyRequired())
	{
		clientRoutes.POST("/prompts/:id/compile", promptHandler.CompilePrompt)
	}

	// 8. Graceful Shutdown setup
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("AI Prompt Management API starting on :%s in %s mode...", cfg.Port, cfg.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Fatal: server failed to start: %v", err)
		}
	}()

	// Wait for OS interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received, shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Fatal: server forced to shutdown: %v", err)
	}

	// Stop background workers
	cancel()
	_ = rdb.Close()
	log.Println("AI Prompt Management API stopped cleanly. Goodbye!")
}
