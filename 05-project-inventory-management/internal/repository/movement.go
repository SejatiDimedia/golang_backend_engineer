package repository

import (
	"context"

	"github.com/timurdian/inventory-management/internal/entity"
	"gorm.io/gorm"
)

type StockMovementRepository interface {
	Create(ctx context.Context, movement *entity.StockMovement) error
	GetAll(ctx context.Context, productID uint, movementType string, page, limit int) ([]entity.StockMovement, int64, error)
}

type stockMovementRepository struct {
	db *gorm.DB
}

func NewStockMovementRepository(db *gorm.DB) StockMovementRepository {
	return &stockMovementRepository{db: db}
}

func (r *stockMovementRepository) Create(ctx context.Context, movement *entity.StockMovement) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Create(movement).Error
}

func (r *stockMovementRepository) GetAll(ctx context.Context, productID uint, movementType string, page, limit int) ([]entity.StockMovement, int64, error) {
	db := GetDBFromContext(ctx, r.db)
	var movements []entity.StockMovement
	var total int64

	query := db.WithContext(ctx).Model(&entity.StockMovement{})

	if productID > 0 {
		query = query.Where("product_id = ?", productID)
	}

	if movementType != "" {
		query = query.Where("type = ?", movementType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Order("id DESC").
		Limit(limit).
		Offset(offset).
		Preload("Product").
		Find(&movements).
		Error

	return movements, total, err
}
