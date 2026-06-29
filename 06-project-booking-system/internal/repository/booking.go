package repository

import (
	"context"
	"errors"
	"time"

	"github.com/timurdian/booking-system/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *entity.Booking) error
	GetByID(ctx context.Context, id uint) (*entity.Booking, error)
	GetByIDForUpdate(ctx context.Context, id uint) (*entity.Booking, error)
	GetOverlapBookings(ctx context.Context, deskID uint, startTime, endTime time.Time) ([]entity.Booking, error)
	GetAll(ctx context.Context) ([]entity.Booking, error)
	GetByUserID(ctx context.Context, userID uint) ([]entity.Booking, error)
	Update(ctx context.Context, booking *entity.Booking) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *entity.Booking) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Create(booking).Error
}

func (r *bookingRepository) GetByID(ctx context.Context, id uint) (*entity.Booking, error) {
	db := GetDBFromContext(ctx, r.db)
	var booking entity.Booking
	err := db.WithContext(ctx).
		Preload("User").
		Preload("Desk").
		First(&booking, id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetByIDForUpdate(ctx context.Context, id uint) (*entity.Booking, error) {
	db := GetDBFromContext(ctx, r.db)
	var booking entity.Booking
	err := db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&booking, id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetOverlapBookings(ctx context.Context, deskID uint, startTime, endTime time.Time) ([]entity.Booking, error) {
	db := GetDBFromContext(ctx, r.db)
	var bookings []entity.Booking
	
	// Query: start_time < requested_end_time AND end_time > requested_start_time AND status = 'CONFIRMED'
	err := db.WithContext(ctx).
		Where("desk_id = ? AND status = ? AND start_time < ? AND end_time > ?", deskID, "CONFIRMED", endTime, startTime).
		Find(&bookings).
		Error
	
	return bookings, err
}

func (r *bookingRepository) GetAll(ctx context.Context) ([]entity.Booking, error) {
	db := GetDBFromContext(ctx, r.db)
	var bookings []entity.Booking
	err := db.WithContext(ctx).
		Preload("User").
		Preload("Desk").
		Order("id DESC").
		Find(&bookings).
		Error
	return bookings, err
}

func (r *bookingRepository) GetByUserID(ctx context.Context, userID uint) ([]entity.Booking, error) {
	db := GetDBFromContext(ctx, r.db)
	var bookings []entity.Booking
	err := db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User").
		Preload("Desk").
		Order("id DESC").
		Find(&bookings).
		Error
	return bookings, err
}

func (r *bookingRepository) Update(ctx context.Context, booking *entity.Booking) error {
	db := GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Save(booking).Error
}
