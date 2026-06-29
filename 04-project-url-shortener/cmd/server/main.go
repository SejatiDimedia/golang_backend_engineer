package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/timurdian/url-shortener/internal/config"
	"github.com/timurdian/url-shortener/internal/entity"
	"github.com/timurdian/url-shortener/internal/handler"
	"github.com/timurdian/url-shortener/internal/repository"
	"github.com/timurdian/url-shortener/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load .env file (jika ada, terutama untuk environment lokal)
	if err := godotenv.Load(); err != nil {
		log.Println("Info: .env file not found, using default environment variables")
	}

	// 2. Load configurations
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Fatal: failed to load configuration: %v", err)
	}

	// 3. Connect to database
	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fatal: failed to connect to database: %v", err)
	}
	log.Println("Database connection established successfully")

	// 4. Auto Migration
	if err := db.AutoMigrate(&entity.URL{}); err != nil {
		log.Fatalf("Fatal: failed to run database auto-migration: %v", err)
	}
	log.Println("Database migration completed successfully")

	// 5. Initialize Layers (Clean Architecture Dependency Injection)
	urlRepo := repository.NewURLRepository(db)
	urlSvc := service.NewURLService(urlRepo)

	healthHandler := handler.NewHealthHandler(db)
	urlHandler := handler.NewURLHandler(urlSvc)

	// 6. Setup Router Gin
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 7. Define Routes
	r.GET("/health", healthHandler.Check)
	r.POST("/shorten", urlHandler.Shorten)
	r.GET("/r/:short_code", urlHandler.Redirect)
	r.GET("/stats/:short_code", urlHandler.Stats)

	// 8. Start HTTP Server
	serverAddress := ":" + cfg.Port
	log.Printf("Server starting on %s in %s mode...", serverAddress, cfg.Env)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		log.Fatalf("Fatal: server failed to start: %v", err)
	}
}
