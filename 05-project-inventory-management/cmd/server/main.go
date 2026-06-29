package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/timurdian/inventory-management/internal/config"
	"github.com/timurdian/inventory-management/internal/entity"
	"github.com/timurdian/inventory-management/internal/handler"
	"github.com/timurdian/inventory-management/internal/repository"
	"github.com/timurdian/inventory-management/internal/service"
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

	// Auto Migration (urutan penting untuk menjaga dependensi FK)
	err = db.AutoMigrate(
		&entity.Category{},
		&entity.Supplier{},
		&entity.Product{},
		&entity.StockMovement{},
	)
	if err != nil {
		log.Fatalf("Fatal: failed to run database auto-migration: %v", err)
	}
	log.Println("Database migration completed successfully")

	// DI: Initialize Repositories
	txManager := repository.NewGormTxManager(db)
	categoryRepo := repository.NewCategoryRepository(db)
	supplierRepo := repository.NewSupplierRepository(db)
	productRepo := repository.NewProductRepository(db)
	movementRepo := repository.NewStockMovementRepository(db)

	// DI: Initialize Services
	categorySvc := service.NewCategoryService(categoryRepo)
	supplierSvc := service.NewSupplierService(supplierRepo)
	productSvc := service.NewProductService(productRepo, categoryRepo, supplierRepo)
	movementSvc := service.NewMovementService(txManager, productRepo, movementRepo)

	// DI: Initialize Handlers
	healthHandler := handler.NewHealthHandler(db)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	supplierHandler := handler.NewSupplierHandler(supplierSvc)
	productHandler := handler.NewProductHandler(productSvc)
	movementHandler := handler.NewMovementHandler(movementSvc)

	// Setup Router Gin
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Define Routes
	r.GET("/health", healthHandler.Check)

	// Categories
	r.POST("/categories", categoryHandler.Create)
	r.GET("/categories", categoryHandler.GetAll)
	r.GET("/categories/:id", categoryHandler.GetByID)
	r.PUT("/categories/:id", categoryHandler.Update)
	r.DELETE("/categories/:id", categoryHandler.Delete)

	// Suppliers
	r.POST("/suppliers", supplierHandler.Create)
	r.GET("/suppliers", supplierHandler.GetAll)
	r.GET("/suppliers/:id", supplierHandler.GetByID)
	r.PUT("/suppliers/:id", supplierHandler.Update)
	r.DELETE("/suppliers/:id", supplierHandler.Delete)

	// Products
	r.POST("/products", productHandler.Create)
	r.GET("/products", productHandler.GetAll)
	r.GET("/products/export", productHandler.ExportCSV) // Letakkan di atas GET /products/:id agar tidak terbentur route path parameter
	r.GET("/products/:id", productHandler.GetByID)
	r.PUT("/products/:id", productHandler.Update)
	r.DELETE("/products/:id", productHandler.Delete)

	// Stock Mutations
	r.POST("/products/:id/stock-in", movementHandler.StockIn)
	r.POST("/products/:id/stock-out", movementHandler.StockOut)
	r.GET("/stock-movements", movementHandler.GetHistory)

	serverAddress := ":" + cfg.Port
	log.Printf("Server starting on %s in %s mode...", serverAddress, cfg.Env)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		log.Fatalf("Fatal: server failed to start: %v", err)
	}
}
