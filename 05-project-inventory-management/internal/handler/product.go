package handler

import (
	"encoding/csv"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/timurdian/inventory-management/internal/entity"
	"github.com/timurdian/inventory-management/internal/service"
)

type ProductHandler struct {
	svc service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var prod entity.Product
	if err := c.ShouldBindJSON(&prod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.Create(c.Request.Context(), &prod)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	prod, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prod)
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	search := c.Query("search")
	
	categoryIDStr := c.Query("category_id")
	var categoryID uint
	if categoryIDStr != "" {
		if cid, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			categoryID = uint(cid)
		}
	}

	pageStr := c.Query("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limitStr := c.Query("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	products, total, err := h.svc.GetAll(c.Request.Context(), search, categoryID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": products,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *ProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var prod entity.Product
	if err := c.ShouldBindJSON(&prod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.Update(c.Request.Context(), uint(id), &prod)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	err = h.svc.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

func (h *ProductHandler) ExportCSV(c *gin.Context) {
	// Ambil semua produk tanpa limit paginasi (untuk ekspor data penuh)
	products, _, err := h.svc.GetAll(c.Request.Context(), "", 0, 1, 100000)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve products: " + err.Error()})
		return
	}

	// Atur HTTP Headers agar browser mendeteksinya sebagai file download CSV
	c.Header("Content-Disposition", "attachment; filename=inventory_report.csv")
	c.Header("Content-Type", "text/csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Tulis Header
	_ = writer.Write([]string{"ID", "SKU", "Name", "Category", "Supplier", "Price", "Stock Quantity"})

	for _, p := range products {
		catName := "N/A"
		if p.Category != nil {
			catName = p.Category.Name
		}
		supName := "N/A"
		if p.Supplier != nil {
			supName = p.Supplier.Name
		}
		row := []string{
			strconv.FormatUint(uint64(p.ID), 10),
			p.SKU,
			p.Name,
			catName,
			supName,
			strconv.FormatFloat(p.Price, 'f', 2, 64),
			strconv.FormatInt(p.StockQuantity, 10),
		}
		_ = writer.Write(row)
	}
}
