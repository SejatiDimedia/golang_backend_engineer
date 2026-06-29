package service

import (
	"context"
	"errors"

	"github.com/timurdian/inventory-management/internal/entity"
	"github.com/timurdian/inventory-management/internal/repository"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock quantity")
)

type MovementService interface {
	StockIn(ctx context.Context, productID uint, quantity int64, reference string) (*entity.StockMovement, error)
	StockOut(ctx context.Context, productID uint, quantity int64, reference string) (*entity.StockMovement, error)
	GetHistory(ctx context.Context, productID uint, movementType string, page, limit int) ([]entity.StockMovement, int64, error)
}

type movementService struct {
	txManager    repository.TransactionManager
	productRepo  repository.ProductRepository
	movementRepo repository.StockMovementRepository
}

func NewMovementService(
	txManager repository.TransactionManager,
	productRepo repository.ProductRepository,
	movementRepo repository.StockMovementRepository,
) MovementService {
	return &movementService{
		txManager:    txManager,
		productRepo:  productRepo,
		movementRepo: movementRepo,
	}
}

func (s *movementService) StockIn(ctx context.Context, productID uint, quantity int64, reference string) (*entity.StockMovement, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	var movement *entity.StockMovement

	// Menjalankan di dalam transaksi database atomis
	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Verifikasi produk eksis
		prod, err := s.productRepo.GetByID(txCtx, productID)
		if err != nil {
			return err
		}
		if prod == nil {
			return ErrProductNotFound
		}

		// 2. Update stock produk
		if err := s.productRepo.UpdateStock(txCtx, productID, quantity); err != nil {
			return err
		}

		// 3. Catat riwayat mutasi
		movement = &entity.StockMovement{
			ProductID: productID,
			Type:      "IN",
			Quantity:  quantity,
			Reference: reference,
		}
		if err := s.movementRepo.Create(txCtx, movement); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return movement, nil
}

func (s *movementService) StockOut(ctx context.Context, productID uint, quantity int64, reference string) (*entity.StockMovement, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	var movement *entity.StockMovement

	// Menjalankan di dalam transaksi database atomis
	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Ambil produk dan kunci baris untuk update (FOR UPDATE)
		prod, err := s.productRepo.GetByIDForUpdate(txCtx, productID)
		if err != nil {
			return err
		}
		if prod == nil {
			return ErrProductNotFound
		}

		// 2. Validasi stok mencukupi
		if prod.StockQuantity < quantity {
			return ErrInsufficientStock
		}

		// 3. Update stock produk (mengurangi kuantitas)
		if err := s.productRepo.UpdateStock(txCtx, productID, -quantity); err != nil {
			return err
		}

		// 4. Catat riwayat mutasi
		movement = &entity.StockMovement{
			ProductID: productID,
			Type:      "OUT",
			Quantity:  quantity,
			Reference: reference,
		}
		if err := s.movementRepo.Create(txCtx, movement); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return movement, nil
}

func (s *movementService) GetHistory(ctx context.Context, productID uint, movementType string, page, limit int) ([]entity.StockMovement, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return s.movementRepo.GetAll(ctx, productID, movementType, page, limit)
}
