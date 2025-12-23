package handlers

import (
	"dormitory-management/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	DB *gorm.DB
}

func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	return &PaymentHandler{DB: db}
}

func (h *PaymentHandler) GetPayments(c *gin.Context) {
	var payments []models.Payment
	query := h.DB.Preload("Resident.User").Preload("Room")

	// Filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by month
	if month := c.Query("month"); month != "" {
		query = query.Where("DATE_FORMAT(month, '%Y-%m') = ?", month)
	}

	// Filter by resident
	if residentID := c.Query("resident_id"); residentID != "" {
		query = query.Where("resident_id = ?", residentID)
	}

	if err := query.Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var paymentReq models.PaymentRequest
	if err := c.ShouldBindJSON(&paymentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse month
	month, err := time.Parse("2006-01", paymentReq.Month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month format. Use YYYY-MM"})
		return
	}

	// Get resident to get room ID
	var resident models.Resident
	if err := h.DB.First(&resident, paymentReq.ResidentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resident not found"})
		return
	}

	if resident.RoomID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resident is not assigned to a room"})
		return
	}

	payment := models.Payment{
		ResidentID:    paymentReq.ResidentID,
		RoomID:        *resident.RoomID,
		Month:         month,
		Amount:        paymentReq.Amount,
		Status:        "pending",
		PaymentMethod: paymentReq.PaymentMethod,
		ReceiptNumber: paymentReq.ReceiptNumber,
		Notes:         paymentReq.Notes,
	}

	if err := h.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	var req struct {
		Status        string `json:"status" binding:"required,oneof=paid pending overdue"`
		PaymentDate   string `json:"payment_date,omitempty"`
		PaymentMethod string `json:"payment_method,omitempty"`
		ReceiptNumber string `json:"receipt_number,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payment models.Payment
	if err := h.DB.First(&payment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	updateData := map[string]interface{}{
		"status": req.Status,
	}

	if req.PaymentMethod != "" {
		updateData["payment_method"] = req.PaymentMethod
	}
	if req.ReceiptNumber != "" {
		updateData["receipt_number"] = req.ReceiptNumber
	}

	if req.Status == "paid" && req.PaymentDate != "" {
		paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
		if err == nil {
			updateData["payment_date"] = paymentDate
		}
	} else if req.Status == "paid" && payment.PaymentDate == nil {
		now := time.Now()
		updateData["payment_date"] = now
	}

	if err := h.DB.Model(&payment).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) GetMonthlyReport(c *gin.Context) {
	month := c.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	var result struct {
		TotalPaid     float64 `json:"total_paid"`
		TotalPending  float64 `json:"total_pending"`
		TotalOverdue  float64 `json:"total_overdue"`
		PaymentCounts struct {
			Paid    int64 `json:"paid"`
			Pending int64 `json:"pending"`
			Overdue int64 `json:"overdue"`
		} `json:"payment_counts"`
	}

	// Get payment statistics
	h.DB.Model(&models.Payment{}).
		Where("DATE_FORMAT(month, '%Y-%m') = ?", month).
		Select("SUM(CASE WHEN status = 'paid' THEN amount ELSE 0 END) as total_paid, " +
			"SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END) as total_pending, " +
			"SUM(CASE WHEN status = 'overdue' THEN amount ELSE 0 END) as total_overdue").
		Scan(&result)

	h.DB.Model(&models.Payment{}).
		Where("DATE_FORMAT(month, '%Y-%m') = ?", month).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&result.PaymentCounts)

	c.JSON(http.StatusOK, result)
}

func (h *PaymentHandler) GetMyPayments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var payments []models.Payment
	query := h.DB.Preload("Resident.User").Preload("Room").
		Joins("JOIN residents ON payments.resident_id = residents.id").
		Where("residents.user_id = ?", userID)

	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("payments.status = ?", status)
	}

	if err := query.Order("payments.month DESC").Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}
