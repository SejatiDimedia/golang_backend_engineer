package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/timurdian/inventory-management/internal/service"
)

type SupplierHandler struct {
	svc service.SupplierService
}

func NewSupplierHandler(svc service.SupplierService) *SupplierHandler {
	return &SupplierHandler{svc: svc}
}

type SupplierRequest struct {
	Name        string `json:"name" binding:"required"`
	ContactName string `json:"contact_name"`
	Email       string `json:"email" binding:"omitempty,email"`
	Phone       string `json:"phone"`
}

func (h *SupplierHandler) Create(c *gin.Context) {
	var req SupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sup, err := h.svc.Create(c.Request.Context(), req.Name, req.ContactName, req.Email, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sup)
}

func (h *SupplierHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	sup, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrSupplierNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sup)
}

func (h *SupplierHandler) GetAll(c *gin.Context) {
	suppliers, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, suppliers)
}

func (h *SupplierHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var req SupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sup, err := h.svc.Update(c.Request.Context(), uint(id), req.Name, req.ContactName, req.Email, req.Phone)
	if err != nil {
		if errors.Is(err, service.ErrSupplierNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sup)
}

func (h *SupplierHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	err = h.svc.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrSupplierNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete supplier: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "supplier deleted successfully"})
}
