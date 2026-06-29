package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/timurdian/prompt-management/internal/middleware"
	"github.com/timurdian/prompt-management/internal/service"
)

type WorkspaceHandler struct {
	svc service.PromptService
}

func NewWorkspaceHandler(svc service.PromptService) *WorkspaceHandler {
	return &WorkspaceHandler{svc: svc}
}

func (h *WorkspaceHandler) CreateWorkspace(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ws, err := h.svc.CreateWorkspace(c.Request.Context(), req.Name, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ws)
}

func (h *WorkspaceHandler) CreateApiKey(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	wsIDStr := c.Param("id")
	wsID, err := strconv.ParseUint(wsIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rawKey, apiKey, err := h.svc.CreateApiKey(c.Request.Context(), uint(wsID), userID, req.Name)
	if err != nil {
		if err == service.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "API Key created successfully. Save it now, it won't be shown again.",
		"api_key":    rawKey,
		"masked_key": apiKey.MaskedKey,
		"id":         apiKey.ID,
		"expires_at": apiKey.ExpiresAt,
	})
}

func (h *WorkspaceHandler) GetApiKeys(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	wsIDStr := c.Param("id")
	wsID, err := strconv.ParseUint(wsIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	keys, err := h.svc.GetWorkspaceApiKeys(c.Request.Context(), uint(wsID), userID)
	if err != nil {
		if err == service.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, keys)
}

func (h *WorkspaceHandler) RevokeApiKey(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	wsIDStr := c.Param("id")
	wsID, err := strconv.ParseUint(wsIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	keyIDStr := c.Param("key_id")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key ID"})
		return
	}

	err = h.svc.RevokeApiKey(c.Request.Context(), uint(wsID), userID, uint(keyID))
	if err != nil {
		if err == service.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API Key revoked successfully"})
}

func (h *WorkspaceHandler) GetWorkspaceAnalytics(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	wsIDStr := c.Param("id")
	wsID, err := strconv.ParseUint(wsIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
		return
	}

	logs, err := h.svc.GetWorkspaceAnalytics(c.Request.Context(), uint(wsID), userID)
	if err != nil {
		if err == service.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
