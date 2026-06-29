package repository

import (
	"context"
	"errors"

	"github.com/timurdian/notification-service/internal/entity"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(ctx context.Context, notif *entity.Notification) error
	Update(ctx context.Context, notif *entity.Notification) error
	GetByID(ctx context.Context, id uint) (*entity.Notification, error)
	CreateLog(ctx context.Context, log *entity.NotificationLog) error
}

type gormNotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &gormNotificationRepository{db: db}
}

func (r *gormNotificationRepository) Create(ctx context.Context, notif *entity.Notification) error {
	return GetDB(ctx, r.db).Create(notif).Error
}

func (r *gormNotificationRepository) Update(ctx context.Context, notif *entity.Notification) error {
	return GetDB(ctx, r.db).Save(notif).Error
}

func (r *gormNotificationRepository) GetByID(ctx context.Context, id uint) (*entity.Notification, error) {
	var notif entity.Notification
	err := GetDB(ctx, r.db).Preload("Logs").First(&notif, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &notif, nil
}

func (r *gormNotificationRepository) CreateLog(ctx context.Context, log *entity.NotificationLog) error {
	return GetDB(ctx, r.db).Create(log).Error
}
