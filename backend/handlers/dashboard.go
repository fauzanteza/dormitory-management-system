package handlers

import (
	"dormitory-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	DB *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{DB: db}
}

func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	var stats models.DashboardStats

	// Room statistics
	h.DB.Model(&models.Room{}).Count(&stats.TotalRooms)
	h.DB.Model(&models.Room{}).Where("status = ?", "occupied").Count(&stats.OccupiedRooms)
	h.DB.Model(&models.Room{}).Where("status = ?", "available").Count(&stats.AvailableRooms)

	// Resident statistics
	h.DB.Model(&models.Resident{}).Where("status = ?", "active").Count(&stats.TotalResidents)

	// Payment statistics
	h.DB.Model(&models.Payment{}).Where("status = ?", "pending").Count(&stats.PendingPayments)

	var totalRevenue struct{ Total float64 }
	h.DB.Model(&models.Payment{}).Where("status = ?", "paid").
		Select("SUM(amount) as total").Scan(&totalRevenue)
	stats.TotalRevenue = totalRevenue.Total

	// Repair statistics
	h.DB.Model(&models.RepairRequest{}).Where("status = ?", "pending").Count(&stats.PendingRepairs)

	// Booking statistics
	h.DB.Model(&models.Booking{}).Where("status = ?", "pending").Count(&stats.PendingBookings)

	// Monthly revenue for last 6 months
	var monthlyRevenues []struct {
		Month   string  `json:"month"`
		Revenue float64 `json:"revenue"`
	}

	sixMonthsAgo := time.Now().AddDate(0, -5, 0)
	h.DB.Model(&models.Payment{}).
		Where("status = ? AND month >= ?", "paid", sixMonthsAgo).
		Select("DATE_FORMAT(month, '%Y-%m') as month, SUM(amount) as revenue").
		Group("DATE_FORMAT(month, '%Y-%m')").
		Order("month DESC").
		Limit(6).
		Scan(&monthlyRevenues)

	stats.MonthlyRevenue = make([]models.MonthlyRevenue, len(monthlyRevenues))
	for i, mr := range monthlyRevenues {
		stats.MonthlyRevenue[i] = models.MonthlyRevenue{
			Month:   mr.Month,
			Revenue: mr.Revenue,
		}
	}

	// Room status distribution
	var roomStatuses []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	h.DB.Model(&models.Room{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&roomStatuses)

	stats.RoomStatus = make([]models.RoomStatus, len(roomStatuses))
	for i, rs := range roomStatuses {
		stats.RoomStatus[i] = models.RoomStatus{
			Status: rs.Status,
			Count:  rs.Count,
		}
	}

	c.JSON(http.StatusOK, stats)
}

func (h *DashboardHandler) GetStudentDashboardStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var stats struct {
		MyRoom          interface{} `json:"my_room"`
		PendingPayments int64       `json:"pending_payments"`
		ActiveRepairs   int64       `json:"active_repairs"`
		UpcomingPayment interface{} `json:"upcoming_payment"`
	}

	// Get student's room
	var resident models.Resident
	if err := h.DB.Preload("Room").Where("user_id = ?", userID).First(&resident).Error; err == nil {
		stats.MyRoom = resident.Room
	}

	// Get pending payments count
	h.DB.Model(&models.Payment{}).
		Joins("JOIN residents ON payments.resident_id = residents.id").
		Where("residents.user_id = ? AND payments.status = ?", userID, "pending").
		Count(&stats.PendingPayments)

	// Get active repairs count
	h.DB.Model(&models.RepairRequest{}).
		Joins("JOIN residents ON repair_requests.resident_id = residents.id").
		Where("residents.user_id = ? AND repair_requests.status IN ?",
			userID, []string{"pending", "in_progress"}).
		Count(&stats.ActiveRepairs)

	// Get upcoming payment
	var upcomingPayment models.Payment
	h.DB.Preload("Room").
		Joins("JOIN residents ON payments.resident_id = residents.id").
		Where("residents.user_id = ? AND payments.status = ? AND payments.month >= ?",
			userID, "pending", time.Now().Format("2006-01-01")).
		Order("payments.month ASC").
		First(&upcomingPayment)

	stats.UpcomingPayment = upcomingPayment

	c.JSON(http.StatusOK, stats)
}
