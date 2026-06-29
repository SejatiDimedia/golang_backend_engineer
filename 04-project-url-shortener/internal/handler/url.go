package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/timurdian/url-shortener/internal/service"
)

type URLHandler struct {
	svc service.URLService
}

func NewURLHandler(svc service.URLService) *URLHandler {
	return &URLHandler{svc: svc}
}

type ShortenRequest struct {
	LongURL          string `json:"long_url" binding:"required,url"`
	CustomAlias      string `json:"custom_alias,omitempty"`
	ExpiresInSeconds *int64 `json:"expires_in_seconds,omitempty"`
}

func (h *URLHandler) Shorten(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body or malformed URL: " + err.Error()})
		return
	}

	var expiresAt *time.Time
	if req.ExpiresInSeconds != nil && *req.ExpiresInSeconds > 0 {
		t := time.Now().Add(time.Duration(*req.ExpiresInSeconds) * time.Second)
		expiresAt = &t
	}

	urlObj, err := h.svc.Shorten(c.Request.Context(), req.LongURL, req.CustomAlias, expiresAt)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrAliasConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to shorten URL: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"short_code": urlObj.ShortCode,
		"long_url":   urlObj.LongURL,
		"expires_at": urlObj.ExpiresAt,
	})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	code := c.Param("short_code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "short code is required"})
		return
	}

	urlObj, err := h.svc.GetAndRecordClick(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrURLExpired) {
			c.JSON(http.StatusGone, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process redirect: " + err.Error()})
		return
	}

	// Lakukan pengalihan HTTP 302 (Found / Temporary Redirect)
	c.Redirect(http.StatusFound, urlObj.LongURL)
}

func (h *URLHandler) Stats(c *gin.Context) {
	code := c.Param("short_code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "short code is required"})
		return
	}

	urlObj, err := h.svc.GetStats(c.Request.Context(), code)
	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stats: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          urlObj.ID,
		"short_code":  urlObj.ShortCode,
		"long_url":    urlObj.LongURL,
		"click_count": urlObj.ClickCount,
		"created_at":  urlObj.CreatedAt,
		"expires_at":  urlObj.ExpiresAt,
	})
}
