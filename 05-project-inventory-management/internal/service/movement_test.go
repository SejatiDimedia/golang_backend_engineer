package service

import (
	"context"
	"errors"
	"testing"

	"github.com/timurdian/inventory-management/internal/entity"
)

type mockTxManager struct{}

func (m *mockTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockProductRepository struct {
	products map[uint]*entity.Product
	dbUpdateError error
}

func (m *mockProductRepository) Create(ctx context.Context, prod *entity.Product) error {
	m.products[prod.ID] = prod
	return nil
}

func (m *mockProductRepository) GetByID(ctx context.Context, id uint) (*entity.Product, error) {
	p, exists := m.products[id]
	if !exists {
		return nil, nil
	}
	return p, nil
}

func (m *mockProductRepository) GetByIDForUpdate(ctx context.Context, id uint) (*entity.Product, error) {
	return m.GetByID(ctx, id)
}

func (m *mockProductRepository) GetAll(ctx context.Context, search string, categoryID uint, page, limit int) ([]entity.Product, int64, error) {
	return nil, 0, nil
}

func (m *mockProductRepository) Update(ctx context.Context, prod *entity.Product) error {
	m.products[prod.ID] = prod
	return nil
}

func (m *mockProductRepository) UpdateStock(ctx context.Context, id uint, qtyChange int64) error {
	if m.dbUpdateError != nil {
		return m.dbUpdateError
	}
	p, exists := m.products[id]
	if !exists {
		return errors.New("product not found")
	}
	p.StockQuantity += qtyChange
	return nil
}

func (m *mockProductRepository) Delete(ctx context.Context, id uint) error {
	return nil
}

type mockStockMovementRepository struct {
	movements []*entity.StockMovement
}

func (m *mockStockMovementRepository) Create(ctx context.Context, movement *entity.StockMovement) error {
	m.movements = append(m.movements, movement)
	return nil
}

func (m *mockStockMovementRepository) GetAll(ctx context.Context, productID uint, movementType string, page, limit int) ([]entity.StockMovement, int64, error) {
	return nil, 0, nil
}

func TestStockIn_Success(t *testing.T) {
	tx := &mockTxManager{}
	prodRepo := &mockProductRepository{
		products: map[uint]*entity.Product{
			1: {ID: 1, Name: "Product Test", StockQuantity: 10},
		},
	}
	moveRepo := &mockStockMovementRepository{}

	svc := NewMovementService(tx, prodRepo, moveRepo)

	res, err := svc.StockIn(context.Background(), 1, 5, "PO-100")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Type != "IN" {
		t.Errorf("expected type IN, got %s", res.Type)
	}

	if res.Quantity != 5 {
		t.Errorf("expected quantity 5, got %d", res.Quantity)
	}

	prod, _ := prodRepo.GetByID(context.Background(), 1)
	if prod.StockQuantity != 15 {
		t.Errorf("expected product stock to be 15, got %d", prod.StockQuantity)
	}
}

func TestStockOut_Success(t *testing.T) {
	tx := &mockTxManager{}
	prodRepo := &mockProductRepository{
		products: map[uint]*entity.Product{
			1: {ID: 1, Name: "Product Test", StockQuantity: 10},
		},
	}
	moveRepo := &mockStockMovementRepository{}

	svc := NewMovementService(tx, prodRepo, moveRepo)

	res, err := svc.StockOut(context.Background(), 1, 4, "SO-200")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Type != "OUT" {
		t.Errorf("expected type OUT, got %s", res.Type)
	}

	prod, _ := prodRepo.GetByID(context.Background(), 1)
	if prod.StockQuantity != 6 {
		t.Errorf("expected product stock to be 6, got %d", prod.StockQuantity)
	}
}

func TestStockOut_InsufficientStock(t *testing.T) {
	tx := &mockTxManager{}
	prodRepo := &mockProductRepository{
		products: map[uint]*entity.Product{
			1: {ID: 1, Name: "Product Test", StockQuantity: 10},
		},
	}
	moveRepo := &mockStockMovementRepository{}

	svc := NewMovementService(tx, prodRepo, moveRepo)

	// Pengurangan stok sebesar 15 (melebihi stok 10)
	_, err := svc.StockOut(context.Background(), 1, 15, "SO-300")
	if !errors.Is(err, ErrInsufficientStock) {
		t.Errorf("expected ErrInsufficientStock error, got %v", err)
	}

	prod, _ := prodRepo.GetByID(context.Background(), 1)
	if prod.StockQuantity != 10 {
		t.Errorf("expected stock to remain 10 after failed txn, got %d", prod.StockQuantity)
	}
}

func TestStockIn_TransactionRollbackOnDBError(t *testing.T) {
	tx := &mockTxManager{}
	prodRepo := &mockProductRepository{
		products: map[uint]*entity.Product{
			1: {ID: 1, Name: "Product Test", StockQuantity: 10},
		},
		dbUpdateError: errors.New("database update failed"),
	}
	moveRepo := &mockStockMovementRepository{}

	svc := NewMovementService(tx, prodRepo, moveRepo)

	_, err := svc.StockIn(context.Background(), 1, 5, "PO-ERR")
	if err == nil {
		t.Error("expected database update error, got nil")
	}

	// Karena error dilempar di callback, di setup real db ini akan di-rollback.
	// Kita pastikan data mutasi tidak dimasukkan ke log history jika db update error.
	if len(moveRepo.movements) != 0 {
		t.Errorf("expected no movements logs created, got %d", len(moveRepo.movements))
	}
}
