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

func (h *PromptHandler) CreatePrompt(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		WorkspaceID uint   `json:"workspace_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
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

	var req struct {
		PromptText string `json:"prompt_text" binding:"required"`
	}
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

	var req struct {
		Variables map[string]string `json:"variables"`
	}
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
