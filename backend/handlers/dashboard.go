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
