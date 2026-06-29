package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/timurdian/digital-wallet/internal/config"
	"github.com/timurdian/digital-wallet/internal/entity"
	"github.com/timurdian/digital-wallet/internal/handler"
	"github.com/timurdian/digital-wallet/internal/middleware"
	"github.com/timurdian/digital-wallet/internal/repository"
	"github.com/timurdian/digital-wallet/internal/service"
	"github.com/timurdian/digital-wallet/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: .env file not found, using default environment variables")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Fatal: failed to load configuration: %v", err)
	}

	// 1. Initialize PostgreSQL
	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fatal: failed to connect to database: %v", err)
	}
	log.Println("PostgreSQL database connection established successfully")

	// Auto Migration
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Wallet{},
		&entity.Transaction{},
	)
	if err != nil {
		log.Fatalf("Fatal: failed to run database auto-migration: %v", err)
	}
	log.Println("PostgreSQL migration completed successfully")

	// 2. Initialize Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.RedisPass,
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Fatal: failed to connect to Redis: %v", err)
	}
	log.Println("Redis client connection established successfully")

	// DI: Initialize Utilities
	lockManager := utils.NewRedisLockManager(rdb)

	// DI: Initialize Repositories
	txManager := repository.NewGormTxManager(db)
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	transRepo := repository.NewTransactionRepository(db)

	// DI: Initialize Services
	userSvc := service.NewUserService(txManager, userRepo, walletRepo)
	walletSvc := service.NewWalletService(txManager, walletRepo, transRepo, rdb, lockManager)

	// DI: Initialize Handlers
	healthHandler := handler.NewHealthHandler(db, rdb)
	userHandler := handler.NewUserHandler(userSvc, cfg.JWTSecret, cfg.JWTExpiryHours)
	walletHandler := handler.NewWalletHandler(walletSvc)

	// Setup Router Gin
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Public Routes
	r.GET("/health", healthHandler.Check)
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// Protected Routes (Gated by JWT AuthMiddleware)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Cek Saldo & Histori mutasi ledger
		protected.GET("/wallet/balance", walletHandler.GetBalance)
		protected.GET("/wallet/transactions", walletHandler.GetTransactions)

		// Idempotency-Protected Routes (Topup, Withdraw, Transfer)
		idempotentRoutes := protected.Group("/")
		idempotentRoutes.Use(middleware.IdempotencyMiddleware(rdb))
		{
			idempotentRoutes.POST("/wallet/top-up", walletHandler.TopUp)
			idempotentRoutes.POST("/wallet/withdraw", walletHandler.Withdraw)
			idempotentRoutes.POST("/wallet/transfer", walletHandler.Transfer)
		}
	}

	serverAddress := ":" + cfg.Port
	log.Printf("Server starting on %s in %s mode...", serverAddress, cfg.Env)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		log.Fatalf("Fatal: server failed to start: %v", err)
	}
}
