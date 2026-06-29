package repository

import (
	"context"
	"errors"

	"github.com/timurdian/inventory-management/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository interface {
	Create(ctx context.Context, prod *entity.Product) error
	GetByID(ctx context.Context, id uint) (*entity.Product, error)
	GetByIDForUpdate(ctx context.Context, id uint) (*entity.Product, error)
	GetAll(ctx context.Context, search string, categoryID uint, page, limit int) ([]entity.Product, int64, error)
	Update(ctx context.Context, prod *entity.Product) error
	UpdateStock(ctx context.Context, id uint, qtyChange int64) error
	Delete(ctx context.Context, id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, prod *entity.Product) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Create(prod).Error
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*entity.Product, error) {
	db := GetDBFromContext(ctx, r.db)
	var prod entity.Product
	err := db.WithContext(ctx).
		Preload("Category").
		Preload("Supplier").
		First(&prod, id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &prod, nil
}

func (r *productRepository) GetByIDForUpdate(ctx context.Context, id uint) (*entity.Product, error) {
	db := GetDBFromContext(ctx, r.db)
	var prod entity.Product
	// Menggunakan SELECT ... FOR UPDATE untuk mengunci baris agar tidak termutasi thread lain
	err := db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&prod, id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &prod, nil
}

func (r *productRepository) GetAll(ctx context.Context, search string, categoryID uint, page, limit int) ([]entity.Product, int64, error) {
	db := GetDBFromContext(ctx, r.db)
	var products []entity.Product
	var total int64

	query := db.WithContext(ctx).Model(&entity.Product{})

	if search != "" {
		query = query.Where("name ILIKE ? OR sku ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	// Hitung total data
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginasi
	offset := (page - 1) * limit
	err := query.Order("id DESC").
		Limit(limit).
		Offset(offset).
		Preload("Category").
		Preload("Supplier").
		Find(&products).
		Error

	return products, total, err
}

func (r *productRepository) Update(ctx context.Context, prod *entity.Product) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Save(prod).Error
}

func (r *productRepository) UpdateStock(ctx context.Context, id uint, qtyChange int64) error {
	db := GetDBFromContext(ctx, r.db)
	// Update atomic
	return db.WithContext(ctx).
		Model(&entity.Product{}).
		Where("id = ?", id).
		Update("stock_quantity", gorm.Expr("stock_quantity + ?", qtyChange)).
		Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Delete(&entity.Product{}, id).Error
}
