package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/timurdian/prompt-management/internal/middleware"
	"github.com/timurdian/prompt-management/internal/service"
)

type PromptHandler struct {
	svc service.PromptService
}

func NewPromptHandler(svc service.PromptService) *PromptHandler {
	return &PromptHandler{svc: svc}
}

type CreatePromptReq struct {
	WorkspaceID uint   `json:"workspace_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreatePrompt godoc
// @Summary      Create a new prompt
// @Description  Creates a new prompt inside a workspace. Requires admin/member JWT authorization.
// @Tags         Prompts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request body CreatePromptReq true "Prompt details"
// @Success      201      {object} entity.Prompt
// @Failure      400      {object} map[string]string "Invalid request payload"
// @Failure      401      {object} map[string]string "Unauthorized"
// @Failure      403      {object} map[string]string "Forbidden access to workspace"
// @Failure      500      {object} map[string]string "Internal server error"
// @Router       /prompts [post]
func (h *PromptHandler) CreatePrompt(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req CreatePromptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prompt, err := h.svc.CreatePrompt(c.Request.Context(), req.WorkspaceID, userID, req.Name, req.Description)
	if err != nil {
		if err == service.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, prompt)
}

// GetPrompt godoc
// @Summary      Get prompt details
// @Description  Retrieve prompt metadata by ID. Requires JWT authorization.
// @Tags         Prompts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Prompt ID"
// @Success      200  {object}  entity.Prompt
// @Failure      400  {object}  map[string]string "Invalid ID"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      403  {object}  map[string]string "Forbidden"
// @Failure      404  {object}  map[string]string "Prompt not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /prompts/{id} [get]
func (h *PromptHandler) GetPrompt(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	promptIDStr := c.Param("id")
	promptID, err := strconv.ParseUint(promptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID"})
		return
	}

	prompt, err := h.svc.GetPrompt(c.Request.Context(), uint(promptID), userID)
	if err != nil {
		if err == service.ErrPromptNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == service.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prompt)
}

// GetWorkspacePrompts godoc
// @Summary      List workspace prompts
// @Description  Retrieve all prompts inside a specific workspace. Requires JWT authorization.
// @Tags         Prompts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Workspace ID"
// @Success      200  {array}   entity.Prompt
// @Failure      400  {object}  map[string]string "Invalid workspace ID"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      403  {object}  map[string]string "Forbidden"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /workspaces/{id}/prompts [get]
func (h *PromptHandler) GetWorkspacePrompts(c *gin.Context) {
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

	prompts, err := h.svc.GetWorkspacePrompts(c.Request.Context(), uint(wsID), userID)
	if err != nil {
		if err == service.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prompts)
}

type CreateVersionReq struct {
	PromptText string `json:"prompt_text" binding:"required"`
}

// CreateVersion godoc
// @Summary      Create a new prompt version
// @Description  Creates a new draft version snapshot of a prompt. Requires JWT authorization.
// @Tags         Prompts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int               true  "Prompt ID"
// @Param        request body   CreateVersionReq  true  "Version details"
// @Success      201  {object}  entity.PromptVersion
// @Failure      400  {object}  map[string]string "Invalid ID or payload"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      403  {object}  map[string]string "Forbidden"
// @Failure      404  {object}  map[string]string "Prompt not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /prompts/{id}/versions [post]
func (h *PromptHandler) CreateVersion(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	promptIDStr := c.Param("id")
	promptID, err := strconv.ParseUint(promptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID"})
		return
	}

	var req CreateVersionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pv, err := h.svc.CreatePromptVersion(c.Request.Context(), uint(promptID), userID, req.PromptText)
	if err != nil {
		if err == service.ErrPromptNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == service.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pv)
}

// ActivateVersion godoc
// @Summary      Activate a prompt version
// @Description  Promotes a draft version of a prompt to ACTIVE status. Requires JWT authorization.
// @Tags         Prompts
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id              path      int  true  "Prompt ID"
// @Param        version_number  path      int  true  "Version Number"
// @Success      200  {object}  map[string]string "Prompt version activated successfully"
// @Failure      400  {object}  map[string]string "Invalid input"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      403  {object}  map[string]string "Forbidden"
// @Failure      404  {object}  map[string]string "Prompt not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /prompts/{id}/versions/{version_number}/activate [put]
func (h *PromptHandler) ActivateVersion(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	promptIDStr := c.Param("id")
	promptID, err := strconv.ParseUint(promptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID"})
		return
	}

	versionStr := c.Param("version_number")
	versionNum, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version number"})
		return
	}

	err = h.svc.ActivatePromptVersion(c.Request.Context(), uint(promptID), userID, versionNum)
	if err != nil {
		if err == service.ErrPromptNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == service.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prompt version activated successfully"})
}

type CompilePromptReq struct {
	Variables map[string]string `json:"variables"`
}

// CompilePrompt godoc
// @Summary      Compile prompt template (Client server-to-server)
// @Description  Replaces placeholders {{var}} in the active prompt version with values. Requires API Key authentication.
// @Tags         Client API
// @Accept       json
// @Produce      json
// @Security     ClientApiKeyAuth
// @Param        id      path      int               true  "Prompt ID"
// @Param        request body     CompilePromptReq  true  "Variables mapping"
// @Success      200     {object}  map[string]interface{} "Returns compiled_prompt and token_estimate"
// @Failure      400     {object}  map[string]string "Invalid payload"
// @Failure      403     {object}  map[string]string "Invalid API key or forbidden access"
// @Failure      404     {object}  map[string]string "Prompt or active version not found"
// @Failure      500     {object}  map[string]string "Internal server error"
// @Router       /client/prompts/{id}/compile [post]
func (h *PromptHandler) CompilePrompt(c *gin.Context) {
	hashVal, exists := c.Get("api_key_hash")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
		return
	}
	hash := hashVal.(string)

	promptIDStr := c.Param("id")
	promptID, err := strconv.ParseUint(promptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID"})
		return
	}

	var req CompilePromptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	compiledText, tokenEstimate, err := h.svc.CompilePrompt(c.Request.Context(), hash, uint(promptID), req.Variables)
	if err != nil {
		if err == service.ErrInvalidApiKey || err == service.ErrUnauthorizedAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err == service.ErrPromptNotFound || err == service.ErrNoActiveVersion {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"compiled_prompt": compiledText,
		"token_estimate":  tokenEstimate,
	})
}
