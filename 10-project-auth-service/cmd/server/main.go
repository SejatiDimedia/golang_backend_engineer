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
	"github.com/timurdian/auth-service/internal/config"
	"github.com/timurdian/auth-service/internal/entity"
	"github.com/timurdian/auth-service/internal/handler"
	"github.com/timurdian/auth-service/internal/repository"
	"github.com/timurdian/auth-service/internal/service"
	"github.com/timurdian/auth-service/internal/utils"
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

	// 3. Ensure RSA Asymmetric Keypair Exists
	log.Println("Initializing cryptographic keys...")
	err = utils.EnsureRSAKeys(cfg.RSAPrivateKeyPath, cfg.RSAPublicKeyPath)
	if err != nil {
		log.Fatalf("Fatal: failed to ensure RSA keypair: %v", err)
	}
	log.Printf("RSA Keypair verified at: %s & %s", cfg.RSAPrivateKeyPath, cfg.RSAPublicKeyPath)

	// 4. Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fatal: failed to connect to database: %v", err)
	}
	log.Println("PostgreSQL connection established successfully")

	// Auto-Migrate Database Schema
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Role{},
		&entity.Permission{},
		&entity.RefreshToken{},
		&entity.VerificationToken{},
		&entity.ResetToken{},
	)
	if err != nil {
		log.Fatalf("Fatal: failed to run database auto-migrations: %v", err)
	}
	log.Println("PostgreSQL auto-migrations completed successfully")

	// 5. Connect to Redis Client
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

	// 6. Initialize Manual Dependency Injection
	tokenMgr, err := utils.NewTokenManager(cfg.RSAPrivateKeyPath, cfg.RSAPublicKeyPath, cfg.AccessTokenExpiryMinutes)
	if err != nil {
		log.Fatalf("Fatal: failed to initialize TokenManager: %v", err)
	}

	authRepo := repository.NewAuthRepository(db)
	txMgr := repository.NewTransactionManager(db)

	authService := service.NewAuthService(
		authRepo,
		txMgr,
		tokenMgr,
		rdb,
		cfg.RefreshTokenExpiryDays,
		cfg.AccessTokenExpiryMinutes,
	)

	authHandler := handler.NewAuthHandler(authService)
	introspectHandler := handler.NewIntrospectHandler(authService)
	rbacHandler := handler.NewRBACHandler(authService)
	healthHandler := handler.NewHealthHandler(db, rdb)

	// 7. Setup Router
	r := gin.Default()
	r.Use(gin.Recovery())

	// Health Check Route
	r.GET("/health", healthHandler.Check)

	// Public Authentication Routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.Refresh)
		authRoutes.POST("/logout", authHandler.Logout)
		authRoutes.GET("/verify-email", authHandler.VerifyEmail)
		authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
		authRoutes.POST("/reset-password", authHandler.ResetPassword)
		authRoutes.POST("/introspect", introspectHandler.Introspect)

		// Admin dynamic RBAC routing setup
		rbac := authRoutes.Group("/rbac")
		{
			rbac.POST("/roles", rbacHandler.CreateRole)
			rbac.POST("/permissions", rbacHandler.CreatePermission)
			rbac.POST("/users/:id/roles", rbacHandler.AssignRole)
			rbac.POST("/roles/:id/permissions", rbacHandler.AssignPermission)
		}
	}

	// 8. Graceful Shutdown Server setup
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Auth Service starting on :%s in %s mode...", cfg.Port, cfg.Env)
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

	_ = rdb.Close()
	log.Println("Auth Service stopped cleanly. Goodbye!")
}
