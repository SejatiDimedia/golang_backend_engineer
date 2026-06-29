package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/timurdian/booking-system/internal/config"
	"github.com/timurdian/booking-system/internal/entity"
	"github.com/timurdian/booking-system/internal/handler"
	"github.com/timurdian/booking-system/internal/middleware"
	"github.com/timurdian/booking-system/internal/repository"
	"github.com/timurdian/booking-system/internal/service"
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

	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fatal: failed to connect to database: %v", err)
	}
	log.Println("Database connection established successfully")

	// Auto Migration (urutan penting)
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Desk{},
		&entity.Booking{},
	)
	if err != nil {
		log.Fatalf("Fatal: failed to run database auto-migration: %v", err)
	}
	log.Println("Database migration completed successfully")

	// DI: Initialize Repositories
	txManager := repository.NewGormTxManager(db)
	userRepo := repository.NewUserRepository(db)
	deskRepo := repository.NewDeskRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	// DI: Initialize Services
	userSvc := service.NewUserService(userRepo)
	deskSvc := service.NewDeskService(deskRepo)
	bookingSvc := service.NewBookingService(txManager, deskRepo, bookingRepo)

	// DI: Initialize Handlers
	healthHandler := handler.NewHealthHandler(db)
	userHandler := handler.NewUserHandler(userSvc, cfg.JWTSecret, cfg.JWTExpiryHours)
	deskHandler := handler.NewDeskHandler(deskSvc)
	bookingHandler := handler.NewBookingHandler(bookingSvc)

	// Setup Router Gin
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 1. Public Routes
	r.GET("/health", healthHandler.Check)
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// 2. Protected Routes (Gated by JWT AuthMiddleware)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Desks (Catalog) - Customer can see active desks
		protected.GET("/desks", deskHandler.GetAllActive)

		// Bookings
		protected.POST("/bookings", bookingHandler.Create)
		protected.GET("/bookings", bookingHandler.List)
		protected.POST("/bookings/:id/cancel", bookingHandler.Cancel)

		// Admin-Only Routes
		adminOnly := protected.Group("/")
		adminOnly.Use(middleware.RequireRole("admin"))
		{
			adminOnly.POST("/desks", deskHandler.Create)
			adminOnly.GET("/admin/desks", deskHandler.GetAll)
			adminOnly.PUT("/desks/:id", deskHandler.Update)
			adminOnly.DELETE("/desks/:id", deskHandler.Delete)
		}
	}

	serverAddress := ":" + cfg.Port
	log.Printf("Server starting on %s in %s mode...", serverAddress, cfg.Env)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		log.Fatalf("Fatal: server failed to start: %v", err)
	}
}
