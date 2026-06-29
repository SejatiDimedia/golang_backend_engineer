package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/timurdian/prompt-management/internal/repository"
	"github.com/timurdian/prompt-management/internal/utils"
)

type APIKeyMiddleware struct {
	repo repository.PromptRepository
	rdb  *redis.Client
}

func NewAPIKeyMiddleware(repo repository.PromptRepository, rdb *redis.Client) *APIKeyMiddleware {
	return &APIKeyMiddleware{repo: repo, rdb: rdb}
}

func (m *APIKeyMiddleware) APIKeyRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be Bearer key"})
			c.Abort()
			return
		}

		rawKey := parts[1]
		if !strings.HasPrefix(rawKey, utils.Prefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key prefix"})
			c.Abort()
			return
		}

		// Hitung SHA-256 Hash
		hash := utils.HashAPIKey(rawKey)
		cacheKey := fmt.Sprintf("apikey:%s", hash)

		// 1. Coba cari di Redis Cache
		ctx := context.Background()
		wsIDStr, err := m.rdb.Get(ctx, cacheKey).Result()
		if err == nil && wsIDStr != "" {
			// Cache Hit! Loloskan otentikasi
			c.Set("api_key_hash", hash)
			c.Next()
			return
		}

		// 2. Cache Miss: Cek PostgreSQL Database
		apiKey, err := m.repo.GetApiKeyByHash(ctx, hash)
		if err != nil || apiKey == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired API key"})
			c.Abort()
			return
		}

		// Simpan hasil ke cache Redis (TTL 1 Jam)
		_ = m.rdb.Set(ctx, cacheKey, fmt.Sprintf("%d", apiKey.WorkspaceID), 1*time.Hour).Err()

		c.Set("api_key_hash", hash)
		c.Next()
	}
}
