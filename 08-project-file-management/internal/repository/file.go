package repository

import (
	"context"
	"errors"

	"github.com/timurdian/file-management/internal/entity"
	"gorm.io/gorm"
)

type FileRepository interface {
	Create(ctx context.Context, file *entity.File) error
	Update(ctx context.Context, file *entity.File) error
	GetByID(ctx context.Context, id uint) (*entity.File, error)
	GetByUserID(ctx context.Context, userID uint) ([]entity.File, error)
	Delete(ctx context.Context, id uint) error
}

type gormFileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &gormFileRepository{db: db}
}

func (r *gormFileRepository) Create(ctx context.Context, file *entity.File) error {
	return GetDB(ctx, r.db).Create(file).Error
}

func (r *gormFileRepository) Update(ctx context.Context, file *entity.File) error {
	return GetDB(ctx, r.db).Save(file).Error
}

func (r *gormFileRepository) GetByID(ctx context.Context, id uint) (*entity.File, error) {
	var file entity.File
	err := GetDB(ctx, r.db).First(&file, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

func (r *gormFileRepository) GetByUserID(ctx context.Context, userID uint) ([]entity.File, error) {
	var files []entity.File
	err := GetDB(ctx, r.db).Where("user_id = ? AND status = ?", userID, "SUCCESS").Order("created_at desc").Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (r *gormFileRepository) Delete(ctx context.Context, id uint) error {
	return GetDB(ctx, r.db).Delete(&entity.File{}, id).Error
}
