package handler

import (
	"net/http"

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
	ctx := c.Request.Context()

	// 1. Check PostgreSQL
	sqlDB, err := h.db.DB()
	dbStatus := "connected"
	if err != nil {
		dbStatus = "error getting db instance"
	} else if err := sqlDB.PingContext(ctx); err != nil {
		dbStatus = "ping failed: " + err.Error()
	}

	// 2. Check Redis
	redisStatus := "connected"
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		redisStatus = "ping failed: " + err.Error()
	}

	status := http.StatusOK
	if dbStatus != "connected" || redisStatus != "connected" {
		status = http.StatusInternalServerError
	}

	c.JSON(status, gin.H{
		"status":   "healthy",
		"database": dbStatus,
		"redis":    redisStatus,
	})
}
