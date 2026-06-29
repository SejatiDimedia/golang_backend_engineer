package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/timurdian/booking-system/internal/entity"
	"github.com/timurdian/booking-system/internal/repository"
)

var (
	ErrDoubleBooking            = errors.New("the room/desk is already booked for this time range")
	ErrBookingNotFound          = errors.New("booking not found")
	ErrCancellationWindowClosed = errors.New("bookings can only be cancelled at least 2 hours before the start time")
	ErrUnauthorizedAction       = errors.New("unauthorized to perform this action")
	ErrInvalidBookingDuration   = errors.New("start time must be before end time and after the current time")
)

type BookingService interface {
	CreateBooking(ctx context.Context, userID uint, userEmail string, deskID uint, startTime, endTime time.Time) (*entity.Booking, error)
	CancelBooking(ctx context.Context, userID uint, userRole string, bookingID uint) error
	GetAllBookings(ctx context.Context) ([]entity.Booking, error)
	GetBookingsByUserID(ctx context.Context, userID uint) ([]entity.Booking, error)
}

type bookingService struct {
	txManager   repository.TransactionManager
	deskRepo    repository.DeskRepository
	bookingRepo repository.BookingRepository
}

func NewBookingService(
	txManager repository.TransactionManager,
	deskRepo repository.DeskRepository,
	bookingRepo repository.BookingRepository,
) BookingService {
	return &bookingService{
		txManager:   txManager,
		deskRepo:    deskRepo,
		bookingRepo: bookingRepo,
	}
}

func (s *bookingService) CreateBooking(ctx context.Context, userID uint, userEmail string, deskID uint, startTime, endTime time.Time) (*entity.Booking, error) {
	// Konversi input ke UTC secara ketat
	startTime = startTime.UTC()
	endTime = endTime.UTC()
	
	now := time.Now().UTC()
	if startTime.Before(now) || !startTime.Before(endTime) {
		return nil, ErrInvalidBookingDuration
	}

	var booking *entity.Booking

	// Jalankan dalam transaksi database atomis
	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Lock data meja/ruangan (SELECT ... FOR UPDATE) untuk mencegah overlap paralel
		desk, err := s.deskRepo.GetByIDForUpdate(txCtx, deskID)
		if err != nil {
			return err
		}
		if desk == nil {
			return ErrDeskNotFound
		}
		if !desk.IsActive {
			return errors.New("the selected desk/room is currently inactive")
		}

		// 2. Cari data booking yang tumpang tindih
		overlaps, err := s.bookingRepo.GetOverlapBookings(txCtx, deskID, startTime, endTime)
		if err != nil {
			return err
		}
		if len(overlaps) > 0 {
			return ErrDoubleBooking
		}

		// 3. Simpan data booking baru
		booking = &entity.Booking{
			UserID:    userID,
			DeskID:    deskID,
			StartTime: startTime,
			EndTime:   endTime,
			Status:    "CONFIRMED",
		}
		if err := s.bookingRepo.Create(txCtx, booking); err != nil {
			return err
		}

		// 4. Pemicu Notifikasi STUB (Menulis ke Standard Output log)
		log.Printf("STUB NOTIFICATION: Booking #%d confirmed for User %s. Room ID: %d, Time: %s to %s (UTC)",
			booking.ID, userEmail, deskID, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.bookingRepo.GetByID(ctx, booking.ID)
}

func (s *bookingService) CancelBooking(ctx context.Context, userID uint, userRole string, bookingID uint) error {
	// Jalankan dalam transaksi database
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Ambil data booking dan kunci baris (FOR UPDATE)
		booking, err := s.bookingRepo.GetByIDForUpdate(txCtx, bookingID)
		if err != nil {
			return err
		}
		if booking == nil {
			return ErrBookingNotFound
		}

		if booking.Status == "CANCELLED" {
			return errors.New("booking is already cancelled")
		}

		// 2. Jika bukan admin, pastikan pemilik booking adalah user aktif
		if userRole != "admin" && booking.UserID != userID {
			return ErrUnauthorizedAction
		}

		// 3. Jika bukan admin, periksa batasan waktu pembatalan (minimal 2 jam sebelum start_time)
		if userRole != "admin" {
			now := time.Now().UTC()
			if booking.StartTime.Sub(now) < 2*time.Hour {
				return ErrCancellationWindowClosed
			}
		}

		// 4. Update status pemesanan
		booking.Status = "CANCELLED"
		if err := s.bookingRepo.Update(txCtx, booking); err != nil {
			return err
		}

		// 5. Pemicu Notifikasi STUB Pembatalan
		log.Printf("STUB NOTIFICATION: Booking #%d CANCELLED. Room ID: %d, Owner User ID: %d",
			booking.ID, booking.DeskID, booking.UserID)

		return nil
	})
}

func (s *bookingService) GetAllBookings(ctx context.Context) ([]entity.Booking, error) {
	return s.bookingRepo.GetAll(ctx)
}

func (s *bookingService) GetBookingsByUserID(ctx context.Context, userID uint) ([]entity.Booking, error) {
	return s.bookingRepo.GetByUserID(ctx, userID)
}
