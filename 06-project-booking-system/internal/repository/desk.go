package repository

import (
	"context"
	"errors"

	"github.com/timurdian/booking-system/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeskRepository interface {
	Create(ctx context.Context, desk *entity.Desk) error
	GetByID(ctx context.Context, id uint) (*entity.Desk, error)
	GetByIDForUpdate(ctx context.Context, id uint) (*entity.Desk, error)
	GetAllActive(ctx context.Context) ([]entity.Desk, error)
	GetAll(ctx context.Context) ([]entity.Desk, error)
	Update(ctx context.Context, desk *entity.Desk) error
	Delete(ctx context.Context, id uint) error
}

type deskRepository struct {
	db *gorm.DB
}

func NewDeskRepository(db *gorm.DB) DeskRepository {
	return &deskRepository{db: db}
}

func (r *deskRepository) Create(ctx context.Context, desk *entity.Desk) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Create(desk).Error
}

func (r *deskRepository) GetByID(ctx context.Context, id uint) (*entity.Desk, error) {
	db := GetDBFromContext(ctx, r.db)
	var desk entity.Desk
	err := db.WithContext(ctx).First(&desk, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &desk, nil
}

func (r *deskRepository) GetByIDForUpdate(ctx context.Context, id uint) (*entity.Desk, error) {
	db := GetDBFromContext(ctx, r.db)
	var desk entity.Desk
	err := db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&desk, id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &desk, nil
}

func (r *deskRepository) GetAllActive(ctx context.Context) ([]entity.Desk, error) {
	db := GetDBFromContext(ctx, r.db)
	var desks []entity.Desk
	err := db.WithContext(ctx).Where("is_active = ?", true).Order("id ASC").Find(&desks).Error
	return desks, err
}

func (r *deskRepository) GetAll(ctx context.Context) ([]entity.Desk, error) {
	db := GetDBFromContext(ctx, r.db)
	var desks []entity.Desk
	err := db.WithContext(ctx).Order("id ASC").Find(&desks).Error
	return desks, err
}

func (r *deskRepository) Update(ctx context.Context, desk *entity.Desk) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Save(desk).Error
}

func (r *deskRepository) Delete(ctx context.Context, id uint) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Delete(&entity.Desk{}, id).Error
}
