package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/timurdian/file-management/internal/utils"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db      *gorm.DB
	storage utils.StorageClient
}

func NewHealthHandler(db *gorm.DB, storage utils.StorageClient) *HealthHandler {
	return &HealthHandler{db: db, storage: storage}
}

func (h *HealthHandler) Check(c *gin.Context) {
	dbStatus := "connected"
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "disconnected"
	}

	minioStatus := "connected"
	if err := h.storage.Ping(c.Request.Context()); err != nil {
		minioStatus = "disconnected"
	}

	status := http.StatusOK
	overall := "healthy"
	if dbStatus == "disconnected" || minioStatus == "disconnected" {
		status = http.StatusServiceUnavailable
		overall = "unhealthy"
	}

	c.JSON(status, gin.H{
		"status":   overall,
		"database": dbStatus,
		"minio":    minioStatus,
	})
}
