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
	"github.com/timurdian/notification-service/internal/config"
	"github.com/timurdian/notification-service/internal/entity"
	"github.com/timurdian/notification-service/internal/handler"
	"github.com/timurdian/notification-service/internal/middleware"
	"github.com/timurdian/notification-service/internal/provider"
	"github.com/timurdian/notification-service/internal/queue"
	"github.com/timurdian/notification-service/internal/repository"
	"github.com/timurdian/notification-service/internal/service"
	"github.com/timurdian/notification-service/internal/worker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load dotenv (Abaikan jika file tidak ada)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 2. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Fatal: failed to load configuration: %v", err)
	}

	// Set Gin Mode
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Context global untuk daemon gracefully shutdown
	rootCtx, cancelRootCtx := context.WithCancel(context.Background())
	defer cancelRootCtx()

	// 3. Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fatal: failed to connect to database: %v", err)
	}
	log.Println("PostgreSQL database connection established successfully")

	// Auto-Migrate Database Schema
	err = db.AutoMigrate(&entity.User{}, &entity.Notification{}, &entity.NotificationLog{})
	if err != nil {
		log.Fatalf("Fatal: failed to run database auto-migration: %v", err)
	}
	log.Println("PostgreSQL auto-migration completed successfully")

	// 4. Connect to Redis Client
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Verify Redis Connection
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer pingCancel()
	if err := rdb.Ping(pingCtx).Err(); err != nil {
		log.Fatalf("Fatal: failed to connect to Redis: %v", err)
	}
	log.Println("Redis client connection established successfully")

	// 5. Initialize Manual Dependency Injection
	queueMgr := queue.NewQueueManager(rdb)

	userRepo := repository.NewUserRepository(db)
	notifRepo := repository.NewNotificationRepository(db)

	userService := service.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	notifService := service.NewNotificationService(notifRepo, queueMgr)

	mockProvider := provider.NewMockNotificationProvider(cfg.ProviderFailureRate)

	// 6. Start Background Daemons
	// Worker Pool Concurrency
	workerPool := worker.NewWorkerPool(queueMgr, notifService, mockProvider, cfg.WorkerConcurrency)
	workerPool.Start(rootCtx)

	// Scheduler Ticker Poller (Move scheduled tasks to ready)
	poller := worker.NewSchedulerPoller(queueMgr)
	poller.Start(rootCtx)

	// 7. Initialize Gin Engine Router
	userHandler := handler.NewUserHandler(userService)
	notifHandler := handler.NewNotificationHandler(notifService)
	healthHandler := handler.NewHealthHandler(db, rdb)

	r := gin.Default()
	r.Use(gin.Recovery())

	// Health Check
	r.GET("/health", healthHandler.Check)

	// Auth Endpoints
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// Protected Notification Endpoints
	notifRoutes := r.Group("/notifications")
	notifRoutes.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		notifRoutes.POST("", notifHandler.Create)
		notifRoutes.GET("/:id", notifHandler.GetStatus)
	}

	// Fallback Route
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
	})

	// 8. Graceful Shutdown Server setup
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on :%s in %s mode...", cfg.Port, cfg.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Fatal: server failed to start: %v", err)
		}
	}()

	// Menunggu OS interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received, shutting down gracefully...")

	// 1. Batalkan context global daemon (menghentikan workers & scheduler poller loop)
	cancelRootCtx()

	// 2. Shut down HTTP Server dengan timeout 5 detik
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Fatal: server forced to shutdown: %v", err)
	}

	// 3. Close Redis connection pool
	_ = rdb.Close()

	log.Println("Server stopped cleanly. Goodbye!")
}
