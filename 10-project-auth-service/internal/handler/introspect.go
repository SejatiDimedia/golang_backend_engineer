package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/timurdian/auth-service/internal/service"
)

type IntrospectHandler struct {
	authService service.AuthService
}

func NewIntrospectHandler(authService service.AuthService) *IntrospectHandler {
	return &IntrospectHandler{authService: authService}
}

type introspectRequest struct {
	Token string `json:"token" binding:"required"`
}

func (h *IntrospectHandler) Introspect(c *gin.Context) {
	var req introspectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	active, claims, err := h.authService.Introspect(c.Request.Context(), req.Token)
	if err != nil || !active {
		c.JSON(http.StatusOK, gin.H{
			"active": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"active":      true,
		"user_id":     claims.UserID,
		"email":       claims.Email,
		"role":        claims.Role,
		"permissions": claims.Permissions,
	})
}
