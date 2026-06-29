package repository

import (
	"context"
	"errors"

	"github.com/timurdian/inventory-management/internal/entity"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, cat *entity.Category) error
	GetByID(ctx context.Context, id uint) (*entity.Category, error)
	GetAll(ctx context.Context) ([]entity.Category, error)
	Update(ctx context.Context, cat *entity.Category) error
	Delete(ctx context.Context, id uint) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, cat *entity.Category) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Create(cat).Error
}

func (r *categoryRepository) GetByID(ctx context.Context, id uint) (*entity.Category, error) {
	db := GetDBFromContext(ctx, r.db)
	var cat entity.Category
	err := db.WithContext(ctx).First(&cat, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepository) GetAll(ctx context.Context) ([]entity.Category, error) {
	db := GetDBFromContext(ctx, r.db)
	var categories []entity.Category
	err := db.WithContext(ctx).Order("id ASC").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) Update(ctx context.Context, cat *entity.Category) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Save(cat).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uint) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Delete(&entity.Category{}, id).Error
}
