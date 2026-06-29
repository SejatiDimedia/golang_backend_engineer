package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/timurdian/booking-system/internal/entity"
)

type mockTxManager struct{}

func (m *mockTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockDeskRepository struct {
	desks map[uint]*entity.Desk
}

func (m *mockDeskRepository) Create(ctx context.Context, desk *entity.Desk) error {
	m.desks[desk.ID] = desk
	return nil
}

func (m *mockDeskRepository) GetByID(ctx context.Context, id uint) (*entity.Desk, error) {
	d, exists := m.desks[id]
	if !exists {
		return nil, nil
	}
	return d, nil
}

func (m *mockDeskRepository) GetByIDForUpdate(ctx context.Context, id uint) (*entity.Desk, error) {
	return m.GetByID(ctx, id)
}

func (m *mockDeskRepository) GetAllActive(ctx context.Context) ([]entity.Desk, error) {
	return nil, nil
}

func (m *mockDeskRepository) GetAll(ctx context.Context) ([]entity.Desk, error) {
	return nil, nil
}

func (m *mockDeskRepository) Update(ctx context.Context, desk *entity.Desk) error {
	m.desks[desk.ID] = desk
	return nil
}

func (m *mockDeskRepository) Delete(ctx context.Context, id uint) error {
	return nil
}

type mockBookingRepository struct {
	bookings map[uint]*entity.Booking
	idCounter uint
}

func (m *mockBookingRepository) Create(ctx context.Context, booking *entity.Booking) error {
	m.idCounter++
	booking.ID = m.idCounter
	m.bookings[booking.ID] = booking
	return nil
}

func (m *mockBookingRepository) GetByID(ctx context.Context, id uint) (*entity.Booking, error) {
	b, exists := m.bookings[id]
	if !exists {
		return nil, nil
	}
	return b, nil
}

func (m *mockBookingRepository) GetByIDForUpdate(ctx context.Context, id uint) (*entity.Booking, error) {
	return m.GetByID(ctx, id)
}

func (m *mockBookingRepository) GetOverlapBookings(ctx context.Context, deskID uint, startTime, endTime time.Time) ([]entity.Booking, error) {
	var overlaps []entity.Booking
	for _, b := range m.bookings {
		if b.DeskID == deskID && b.Status == "CONFIRMED" {
			// Kriteria overlap check: b.StartTime < endTime AND b.EndTime > startTime
			if b.StartTime.Before(endTime) && b.EndTime.After(startTime) {
				overlaps = append(overlaps, *b)
			}
		}
	}
	return overlaps, nil
}

func (m *mockBookingRepository) GetAll(ctx context.Context) ([]entity.Booking, error) {
	return nil, nil
}

func (m *mockBookingRepository) GetByUserID(ctx context.Context, userID uint) ([]entity.Booking, error) {
	return nil, nil
}

func (m *mockBookingRepository) Update(ctx context.Context, booking *entity.Booking) error {
	m.bookings[booking.ID] = booking
	return nil
}

func TestCreateBooking_Success(t *testing.T) {
	tx := &mockTxManager{}
	deskRepo := &mockDeskRepository{
		desks: map[uint]*entity.Desk{
			1: {ID: 1, Name: "Meja 1", Type: "hot-desk", IsActive: true},
		},
	}
	bookingRepo := &mockBookingRepository{
		bookings: make(map[uint]*entity.Booking),
	}

	svc := NewBookingService(tx, deskRepo, bookingRepo)

	start := time.Now().Add(1 * time.Hour)
	end := start.Add(2 * time.Hour)

	res, err := svc.CreateBooking(context.Background(), 10, "user@email.com", 1, start, end)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Status != "CONFIRMED" {
		t.Errorf("expected status CONFIRMED, got %s", res.Status)
	}

	if res.UserID != 10 {
		t.Errorf("expected userID 10, got %d", res.UserID)
	}
}

func TestCreateBooking_OverlapConflict(t *testing.T) {
	tx := &mockTxManager{}
	deskRepo := &mockDeskRepository{
		desks: map[uint]*entity.Desk{
			1: {ID: 1, Name: "Meja 1", Type: "hot-desk", IsActive: true},
		},
	}
	
	// Siapkan data booking terdaftar: 14:00 s/d 16:00
	registeredStart := time.Now().Add(2 * time.Hour)
	registeredEnd := registeredStart.Add(2 * time.Hour)
	
	bookingRepo := &mockBookingRepository{
		bookings: map[uint]*entity.Booking{
			1: {ID: 1, UserID: 5, DeskID: 1, StartTime: registeredStart, EndTime: registeredEnd, Status: "CONFIRMED"},
		},
	}

	svc := NewBookingService(tx, deskRepo, bookingRepo)

	// Uji overlap: user lain mencoba memesan 15:00 s/d 17:00 (overlap di jam 15:00-16:00)
	testStart := registeredStart.Add(1 * time.Hour)
	testEnd := testStart.Add(2 * time.Hour)

	_, err := svc.CreateBooking(context.Background(), 10, "spammer@email.com", 1, testStart, testEnd)
	if !errors.Is(err, ErrDoubleBooking) {
		t.Errorf("expected ErrDoubleBooking error, got %v", err)
	}
}

func TestCancelBooking_CancellationWindowClosed(t *testing.T) {
	tx := &mockTxManager{}
	deskRepo := &mockDeskRepository{}
	
	// Siapkan booking yang akan dimulai dalam 1 jam (kurang dari batasan 2 jam pembatalan)
	bookingStart := time.Now().Add(1 * time.Hour)
	bookingEnd := bookingStart.Add(2 * time.Hour)
	
	bookingRepo := &mockBookingRepository{
		bookings: map[uint]*entity.Booking{
			1: {ID: 1, UserID: 10, DeskID: 1, StartTime: bookingStart, EndTime: bookingEnd, Status: "CONFIRMED"},
		},
	}

	svc := NewBookingService(tx, deskRepo, bookingRepo)

	// Customer mencoba membatalkan -> harus ditolak
	err := svc.CancelBooking(context.Background(), 10, "customer", 1)
	if !errors.Is(err, ErrCancellationWindowClosed) {
		t.Errorf("expected ErrCancellationWindowClosed error, got %v", err)
	}

	// Admin mencoba membatalkan -> harus diizinkan (bypass window check)
	err = svc.CancelBooking(context.Background(), 99, "admin", 1)
	if err != nil {
		t.Errorf("expected admin to bypass cancellation constraint, but got error: %v", err)
	}

	b, _ := bookingRepo.GetByID(context.Background(), 1)
	if b.Status != "CANCELLED" {
		t.Errorf("expected status to be CANCELLED after admin cancellation, got %s", b.Status)
	}
}
