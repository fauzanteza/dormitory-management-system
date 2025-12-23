package handlers

import (
	"dormitory-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingHandler struct {
	DB *gorm.DB
}

func NewBookingHandler(db *gorm.DB) *BookingHandler {
	return &BookingHandler{DB: db}
}

// GetBookings lists bookings. Admin sees all, Student sees own.
func (h *BookingHandler) GetBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var bookings []models.Booking
	query := h.DB.Preload("User").Preload("Room").Order("created_at desc")

	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

type CreateBookingInput struct {
	RoomID         uint   `json:"room_id" binding:"required"`
	StartDate      string `json:"start_date" binding:"required"` // YYYY-MM-DD
	DurationMonths int    `json:"duration_months"`
	Notes          string `json:"notes"`
}

// CreateBooking allows a student to book a room
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var input CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Double check room availability
	var room models.Room
	if err := h.DB.First(&room, input.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	if room.Status != "available" || room.CurrentOccupancy >= room.Capacity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kamar tidak tersedia atau penuh"})
		return
	}

	// Check if user already has a pending booking or is already a resident
	var existingBooking models.Booking
	if err := h.DB.Where("user_id = ? AND status = 'pending'", userID).First(&existingBooking).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda masih memiliki pemesanan yang menunggu persetujuan"})
		return
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	booking := models.Booking{
		UserID:         userID.(uint),
		RoomID:         input.RoomID,
		StartDate:      startDate,
		DurationMonths: input.DurationMonths,
		Notes:          input.Notes,
		Status:         "pending",
	}

	if booking.DurationMonths == 0 {
		booking.DurationMonths = 12 // Default 1 year
	}

	if err := h.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pemesanan"})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

type UpdateBookingStatusInput struct {
	Status string `json:"status" binding:"required"` // approved, rejected
}

// UpdateBookingStatus (Admin Only) handles approval/rejection
func (h *BookingHandler) UpdateBookingStatus(c *gin.Context) {
	id := c.Param("id")
	var input UpdateBookingStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Status != "approved" && input.Status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	var booking models.Booking
	if err := h.DB.Preload("User").Preload("Room").First(&booking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if booking.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking is already processed"})
		return
	}

	tx := h.DB.Begin()

	booking.Status = input.Status
	if err := tx.Save(&booking).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}

	if input.Status == "approved" {
		// Logic to approve:
		// 1. Check room capacity again (race condition check)
		var room models.Room
		if err := tx.First(&room, booking.RoomID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room not found"})
			return
		}

		if room.CurrentOccupancy >= room.Capacity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kamar sudah penuh, tidak bisa menyetujui"})
			return
		}

		// 2. Increment Room Occupancy
		room.CurrentOccupancy++
		if room.CurrentOccupancy >= room.Capacity {
			room.Status = "occupied"
		}
		if err := tx.Save(&room).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room"})
			return
		}

		// 3. Create Resident Record
		// Note: We don't have all resident details like StudentID, Faculty, Major in bookings table usually.
		// For MVP, we can leave them blank or assume User Profile has them?
		// Or we just create a basic Resident record linking User <-> Room.
		// Let's create basic record. Admin can edit details later.
		resident := models.Resident{
			UserID:      booking.UserID,
			RoomID:      &booking.RoomID,
			Status:      "active",
			StudentID:   "-", // Placeholder
			YearOfEntry: time.Now().Year(),
		}

		// Check if resident already exists for this user (e.g. re-booking?)
		var existingResident models.Resident
		if err := tx.Where("user_id = ?", booking.UserID).First(&existingResident).Error; err == nil {
			// Update existing
			existingResident.RoomID = &booking.RoomID
			existingResident.Status = "active"
			if err := tx.Save(&existingResident).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resident"})
				return
			}
			resident = existingResident // Set for payment creation
		} else {
			// Create new
			if err := tx.Create(&resident).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Proses gagal: User mungkin sudah menjadi penghuni"})
				return
			}
		}

		// 4. Create Initial Payment (Tagihan Bulan Pertama)
		// Booking Start Date determines the month
		paymentMonth := booking.StartDate
		// Ideally set day to 1
		paymentMonth = time.Date(paymentMonth.Year(), paymentMonth.Month(), 1, 0, 0, 0, 0, paymentMonth.Location())

		payment := models.Payment{
			ResidentID: resident.ID,
			RoomID:     room.ID,
			Month:      paymentMonth,
			Amount:     room.MonthlyRate,
			Status:     "pending",
			Notes:      "Tagihan awal dari Booking",
		}
		// Use Create, ignore error if unique constraint (already billed)
		tx.Create(&payment)
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Booking status updated to " + input.Status,
		"booking": booking,
	})
}
