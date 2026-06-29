package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/timurdian/notification-service/internal/service"
)

type NotificationHandler struct {
	notifService service.NotificationService
}

func NewNotificationHandler(notifService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifService: notifService}
}

type createNotificationRequest struct {
	Type    string `json:"type" binding:"required"`
	Target  string `json:"target" binding:"required"`
	Content string `json:"content" binding:"required"`
	SendAt  string `json:"send_at,omitempty"`
}

func (h *NotificationHandler) Create(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sendAt *time.Time
	if req.SendAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.SendAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid send_at timestamp format, use RFC3339 (e.g. 2026-06-29T15:00:00Z)"})
			return
		}
		sendAt = &parsedTime
	}

	notif, err := h.notifService.Create(c.Request.Context(), req.Type, req.Target, req.Content, sendAt)
	if err != nil {
		if err == service.ErrInvalidNotificationType || err == service.ErrInvalidTarget || err == service.ErrInvalidContent {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Sesuai PRD, notifikasi asinkron mengembalikan 202 Accepted
	c.JSON(http.StatusAccepted, gin.H{
		"message":         "notification queued successfully",
		"notification_id": notif.ID,
		"status":          notif.Status,
		"send_at":         notif.SendAt,
	})
}

func (h *NotificationHandler) GetStatus(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID format"})
		return
	}

	notif, err := h.notifService.GetStatus(c.Request.Context(), uint(id))
	if err != nil {
		if err == service.ErrNotificationNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notif)
}
