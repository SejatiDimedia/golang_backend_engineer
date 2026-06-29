package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// IdempotencyMiddleware memastikan request yang sama tidak dieksekusi berulang kali
func IdempotencyMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idempotencyKey := c.GetHeader("X-Idempotency-Key")
		if idempotencyKey == "" {
			c.Next()
			return
		}

		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized for idempotency check"})
			c.Abort()
			return
		}
		userID := userIDVal.(uint)

		// Namespace key di Redis: idempotency:user:<user_id>:<key_hash>
		hash := sha256.Sum256([]byte(idempotencyKey))
		redisKey := fmt.Sprintf("idempotency:user:%d:%s", userID, hex.EncodeToString(hash[:]))

		// 1. Cek apakah request sudah pernah sukses dijalankan sebelumnya
		val, err := rdb.Get(c.Request.Context(), redisKey).Result()
		if err == nil && val != "" {
			c.Header("X-Cache-Lookup", "HIT")
			c.Data(http.StatusOK, "application/json", []byte(val))
			c.Abort()
			return
		}

		// 2. Intercept response body
		bw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bw

		c.Next()

		// 3. Simpan response sukses (2xx) ke Redis selama 1 jam (3600 detik)
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			rdb.Set(c.Request.Context(), redisKey, bw.body.String(), 1*time.Hour)
		}
	}
}
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}
