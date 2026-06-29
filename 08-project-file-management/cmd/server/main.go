package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/timurdian/file-management/internal/config"
	"github.com/timurdian/file-management/internal/entity"
	"github.com/timurdian/file-management/internal/handler"
	"github.com/timurdian/file-management/internal/middleware"
	"github.com/timurdian/file-management/internal/repository"
	"github.com/timurdian/file-management/internal/service"
	"github.com/timurdian/file-management/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load dotenv (Abaikan jika file tidak ditemukan)
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

	// 3. Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fatal: failed to connect to database: %v", err)
	}
	log.Println("PostgreSQL database connection established successfully")

	// Auto-Migrate Database Schema
	err = db.AutoMigrate(&entity.User{}, &entity.File{})
	if err != nil {
		log.Fatalf("Fatal: failed to run auto-migration: %v", err)
	}
	log.Println("PostgreSQL auto-migration completed successfully")

	// 4. Auto-Create MinIO Bucket Target (Ensure bucket exists)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = utils.EnsureBucketExists(
		ctx,
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioUseSSL,
		cfg.MinioBucketName,
	)
	if err != nil {
		log.Fatalf("Fatal: failed to verify or create MinIO bucket: %v", err)
	}

	// 5. Initialize MinIO Client Adapter
	storageClient, err := utils.NewMinioStorageClient(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioUseSSL,
	)
	if err != nil {
		log.Fatalf("Fatal: failed to initialize MinIO storage adapter: %v", err)
	}
	log.Println("MinIO Client adapter established successfully")

	// 6. Manual Dependency Injection Scaffolding
	userRepo := repository.NewUserRepository(db)
	fileRepo := repository.NewFileRepository(db)

	userService := service.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	fileService := service.NewFileService(
		fileRepo,
		storageClient,
		cfg.MinioBucketName,
		cfg.MinioPresignedExpiryMinutes,
	)

	userHandler := handler.NewUserHandler(userService)
	fileHandler := handler.NewFileHandler(fileService)
	healthHandler := handler.NewHealthHandler(db, storageClient)

	// 7. Initialize Gin Engine Router
	r := gin.Default()

	// Global Middlewares
	r.Use(gin.Recovery())

	// Health Check
	r.GET("/health", healthHandler.Check)

	// Authentication Routes
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// File Management Routes (Protected by AuthMiddleware)
	fileRoutes := r.Group("/files")
	fileRoutes.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		fileRoutes.POST("/upload", fileHandler.Upload)
		fileRoutes.GET("", fileHandler.GetList)
		fileRoutes.GET("/:id/download", fileHandler.GetDownloadURL)
		fileRoutes.GET("/:id/view", fileHandler.ViewStream)
		fileRoutes.DELETE("/:id", fileHandler.Delete)
	}

	// Fallback Route
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
	})

	log.Printf("Server starting on :%s in %s mode...", cfg.Port, cfg.Env)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Fatal: server failed to start: %v", err)
	}
}
