package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/timurdian/digital-wallet/internal/service"
)

type WalletHandler struct {
	svc service.WalletService
}

func NewWalletHandler(svc service.WalletService) *WalletHandler {
	return &WalletHandler{svc: svc}
}

type TopUpRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

type WithdrawRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

type TransferRequest struct {
	DestinationWalletNumber string  `json:"destination_wallet_number" binding:"required"`
	Amount                  float64 `json:"amount" binding:"required,gt=0"`
	Description             string  `json:"description"`
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	balance, err := h.svc.GetBalance(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrWalletNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}

func (h *WalletHandler) TopUp(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	var req TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.svc.TopUp(c.Request.Context(), userID, req.Amount, req.Description)
	if err != nil {
		if errors.Is(err, service.ErrWalletNotFound) || errors.Is(err, service.ErrInvalidAmount) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tx)
}

func (h *WalletHandler) Withdraw(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.svc.Withdraw(c.Request.Context(), userID, req.Amount, req.Description)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientBalance) || errors.Is(err, service.ErrInvalidAmount) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tx)
}

func (h *WalletHandler) Transfer(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.svc.Transfer(c.Request.Context(), userID, req.DestinationWalletNumber, req.Amount, req.Description)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientBalance) ||
			errors.Is(err, service.ErrWalletNotFound) ||
			errors.Is(err, service.ErrInvalidTransfer) ||
			errors.Is(err, service.ErrInvalidAmount) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tx)
}

func (h *WalletHandler) GetTransactions(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	txs, err := h.svc.GetTransactions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, txs)
}
