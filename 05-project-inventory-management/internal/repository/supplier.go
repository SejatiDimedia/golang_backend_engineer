package repository

import (
	"context"
	"errors"

	"github.com/timurdian/inventory-management/internal/entity"
	"gorm.io/gorm"
)

type SupplierRepository interface {
	Create(ctx context.Context, sup *entity.Supplier) error
	GetByID(ctx context.Context, id uint) (*entity.Supplier, error)
	GetAll(ctx context.Context) ([]entity.Supplier, error)
	Update(ctx context.Context, sup *entity.Supplier) error
	Delete(ctx context.Context, id uint) error
}

type supplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &supplierRepository{db: db}
}

func (r *supplierRepository) Create(ctx context.Context, sup *entity.Supplier) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Create(sup).Error
}

func (r *supplierRepository) GetByID(ctx context.Context, id uint) (*entity.Supplier, error) {
	db := GetDBFromContext(ctx, r.db)
	var sup entity.Supplier
	err := db.WithContext(ctx).First(&sup, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sup, nil
}

func (r *supplierRepository) GetAll(ctx context.Context) ([]entity.Supplier, error) {
	db := GetDBFromContext(ctx, r.db)
	var suppliers []entity.Supplier
	err := db.WithContext(ctx).Order("id ASC").Find(&suppliers).Error
	return suppliers, err
}

func (r *supplierRepository) Update(ctx context.Context, sup *entity.Supplier) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Save(sup).Error
}

func (r *supplierRepository) Delete(ctx context.Context, id uint) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Delete(&entity.Supplier{}, id).Error
}
