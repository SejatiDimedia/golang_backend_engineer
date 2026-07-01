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

type CreateWorkspaceReq struct {
	Name string `json:"name" binding:"required"`
}

// CreateWorkspace godoc
// @Summary      Create a new workspace
// @Description  Creates a new workspace. The creator automatically becomes the workspace Admin. Requires JWT authorization.
// @Tags         Workspaces
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request body CreateWorkspaceReq true "Workspace name"
// @Success      201      {object} entity.Workspace
// @Failure      400      {object} map[string]string "Invalid request payload"
// @Failure      401      {object} map[string]string "Unauthorized"
// @Failure      500      {object} map[string]string "Internal server error"
// @Router       /workspaces [post]
func (h *WorkspaceHandler) CreateWorkspace(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req CreateWorkspaceReq
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

type CreateApiKeyReq struct {
	Name string `json:"name" binding:"required"`
}

// CreateApiKey godoc
// @Summary      Create API Key for workspace
// @Description  Generates a new API Key for client integration. The key string is only shown once in the response. Requires workspace Admin JWT authorization.
// @Tags         API Keys
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int              true  "Workspace ID"
// @Param        request body   CreateApiKeyReq  true  "Key descriptor name"
// @Success      201  {object}  map[string]interface{} "Returns raw api_key, masked_key, id, and expires_at"
// @Failure      400  {object}  map[string]string "Invalid input"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      403  {object}  map[string]string "Forbidden access to workspace"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /workspaces/{id}/api-keys [post]
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

	var req CreateApiKeyReq
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

// GetApiKeys godoc
// @Summary      Get workspace API Keys
// @Description  Lists all API Keys generated for a workspace (masked for security). Requires workspace Admin JWT authorization.
// @Tags         API Keys
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Workspace ID"
// @Success      200  {array}   entity.ApiKey
// @Failure      400  {object}  map[string]string "Invalid workspace ID"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      403  {object}  map[string]string "Forbidden"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /workspaces/{id}/api-keys [get]
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

// RevokeApiKey godoc
// @Summary      Revoke an API Key
// @Description  Deletes and revokes an API Key, invalidating it from cache instantly. Requires workspace Admin JWT authorization.
// @Tags         API Keys
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id       path      int  true  "Workspace ID"
// @Param        key_id   path      int  true  "API Key ID"
// @Success      200  {object}  map[string]string "API Key revoked successfully"
// @Failure      400  {object}  map[string]string "Invalid input IDs"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      403  {object}  map[string]string "Forbidden"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /workspaces/{id}/api-keys/{key_id} [delete]
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

// GetWorkspaceAnalytics godoc
// @Summary      Get workspace execution analytics logs
// @Description  Retrieves logs of prompt compiler executions (latency, token estimates). Requires workspace Admin/Member JWT authorization.
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Workspace ID"
// @Success      200  {array}   entity.AnalyticsLog
// @Failure      400  {object}  map[string]string "Invalid workspace ID"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      403  {object}  map[string]string "Forbidden"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /workspaces/{id}/analytics [get]
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
