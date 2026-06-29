package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewHealthHandler(db *gorm.DB, rdb *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, rdb: rdb}
}

func (h *HealthHandler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "connected"
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "disconnected"
	}

	redisStatus := "connected"
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		redisStatus = "disconnected"
	}

	status := http.StatusOK
	overall := "healthy"
	if dbStatus == "disconnected" || redisStatus == "disconnected" {
		status = http.StatusServiceUnavailable
		overall = "unhealthy"
	}

	c.JSON(status, gin.H{
		"status":   overall,
		"database": dbStatus,
		"redis":    redisStatus,
	})
}
