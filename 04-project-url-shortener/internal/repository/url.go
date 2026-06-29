package repository

import (
	"context"
	"errors"

	"github.com/timurdian/url-shortener/internal/entity"
	"gorm.io/gorm"
)

type URLRepository interface {
	Create(ctx context.Context, url *entity.URL) error
	GetByShortCode(ctx context.Context, code string) (*entity.URL, error)
	IncrementClick(ctx context.Context, code string) error
}

type urlRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Create(ctx context.Context, url *entity.URL) error {
	return r.db.WithContext(ctx).Create(url).Error
}

func (r *urlRepository) GetByShortCode(ctx context.Context, code string) (*entity.URL, error) {
	var url entity.URL
	err := r.db.WithContext(ctx).Where("short_code = ?", code).First(&url).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // return nil, nil jika tidak ditemukan (idiomatik, atau bisa dengan custom error)
		}
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) IncrementClick(ctx context.Context, code string) error {
	// Menggunakan update atomic untuk menghindari race condition
	return r.db.WithContext(ctx).
		Model(&entity.URL{}).
		Where("short_code = ?", code).
		Update("click_count", gorm.Expr("click_count + 1")).
		Error
}
