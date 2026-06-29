package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/timurdian/booking-system/internal/entity"
	"github.com/timurdian/booking-system/internal/service"
)

type BookingHandler struct {
	svc service.BookingService
}

func NewBookingHandler(svc service.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

type BookingRequest struct {
	DeskID    uint      `json:"desk_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	EndTime   time.Time `json:"end_time" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
}

func (h *BookingHandler) Create(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	userEmailVal, emailExists := c.Get("email")
	if !exists || !emailExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized"})
		return
	}
	
	userID := userIDVal.(uint)
	userEmail := userEmailVal.(string)

	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload or date-time format: " + err.Error()})
		return
	}

	booking, err := h.svc.CreateBooking(c.Request.Context(), userID, userEmail, req.DeskID, req.StartTime, req.EndTime)
	if err != nil {
		if errors.Is(err, service.ErrInvalidBookingDuration) || errors.Is(err, service.ErrDoubleBooking) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrDeskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func (h *BookingHandler) List(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	userRoleVal, roleExists := c.Get("role")
	if !exists || !roleExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized"})
		return
	}
	
	userID := userIDVal.(uint)
	userRole := userRoleVal.(string)

	var bookings []entity.Booking
	var err error

	// Jika admin, ambil seluruh booking. Jika customer, filter miliknya saja.
	if userRole == "admin" {
		bookings, err = h.svc.GetAllBookings(c.Request.Context())
	} else {
		bookings, err = h.svc.GetBookingsByUserID(c.Request.Context(), userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) Cancel(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	userRoleVal, roleExists := c.Get("role")
	if !exists || !roleExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user unauthorized"})
		return
	}
	
	userID := userIDVal.(uint)
	userRole := userRoleVal.(string)

	idStr := c.Param("id")
	bookingID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	err = h.svc.CancelBooking(c.Request.Context(), userID, userRole, uint(bookingID))
	if err != nil {
		if errors.Is(err, service.ErrBookingNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrUnauthorizedAction) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrCancellationWindowClosed) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled successfully"})
}
