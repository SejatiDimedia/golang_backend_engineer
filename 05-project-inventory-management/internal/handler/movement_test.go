package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/timurdian/inventory-management/internal/entity"
	"github.com/timurdian/inventory-management/internal/service"
)

type mockMovementService struct {
	stockInFunc    func(ctx context.Context, productID uint, quantity int64, reference string) (*entity.StockMovement, error)
	stockOutFunc   func(ctx context.Context, productID uint, quantity int64, reference string) (*entity.StockMovement, error)
	getHistoryFunc func(ctx context.Context, productID uint, movementType string, page, limit int) ([]entity.StockMovement, int64, error)
}

func (m *mockMovementService) StockIn(ctx context.Context, productID uint, quantity int64, reference string) (*entity.StockMovement, error) {
	return m.stockInFunc(ctx, productID, quantity, reference)
}

func (m *mockMovementService) StockOut(ctx context.Context, productID uint, quantity int64, reference string) (*entity.StockMovement, error) {
	return m.stockOutFunc(ctx, productID, quantity, reference)
}

func (m *mockMovementService) GetHistory(ctx context.Context, productID uint, movementType string, page, limit int) ([]entity.StockMovement, int64, error) {
	return m.getHistoryFunc(ctx, productID, movementType, page, limit)
}

func TestStockIn_HandlerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockMovementService{
		stockInFunc: func(ctx context.Context, productID uint, quantity int64, reference string) (*entity.StockMovement, error) {
			return &entity.StockMovement{
				ID:        1,
				ProductID: productID,
				Type:      "IN",
				Quantity:  quantity,
				Reference: reference,
			}, nil
		},
	}

	h := NewMovementHandler(mockSvc)
	r := gin.Default()
	r.POST("/products/:id/stock-in", h.StockIn)

	payload := MutationRequest{
		Quantity:  10,
		Reference: "PO-001",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/products/1/stock-in", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", w.Code)
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["type"] != "IN" || resp["quantity"].(float64) != 10 {
		t.Errorf("unexpected response content: %v", resp)
	}
}

func TestStockOut_HandlerInsufficient(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &mockMovementService{
		stockOutFunc: func(ctx context.Context, productID uint, quantity int64, reference string) (*entity.StockMovement, error) {
			return nil, service.ErrInsufficientStock
		},
	}

	h := NewMovementHandler(mockSvc)
	r := gin.Default()
	r.POST("/products/:id/stock-out", h.StockOut)

	payload := MutationRequest{
		Quantity:  50,
		Reference: "SO-002",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/products/1/stock-out", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", w.Code)
	}
}
