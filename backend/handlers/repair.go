package handlers

import (
	"dormitory-management/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RepairHandler struct {
	DB *gorm.DB
}

func NewRepairHandler(db *gorm.DB) *RepairHandler {
	return &RepairHandler{DB: db}
}

func (h *RepairHandler) GetRepairRequests(c *gin.Context) {
	var repairs []models.RepairRequest
	query := h.DB.Preload("Resident.User").Preload("Room").Preload("AssignedUser")

	// Filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by priority
	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
	}

	// Filter by resident (for student view)
	role, _ := c.Get("role")
	if role == "student" {
		userID, _ := c.Get("user_id")
		query = query.Joins("JOIN residents ON repair_requests.resident_id = residents.id").
			Where("residents.user_id = ?", userID)
	}

	if err := query.Order("created_at DESC").Find(&repairs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repairs)
}

func (h *RepairHandler) CreateRepairRequest(c *gin.Context) {
	var repairReq models.RepairRequestInput
	if err := c.ShouldBindJSON(&repairReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For students, get resident ID from user
	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")

	var resident models.Resident
	if role == "student" {
		if err := h.DB.Where("user_id = ?", userID).First(&resident).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Resident profile not found"})
			return
		}
	} else {
		// For admin, use the provided resident ID from room
		if err := h.DB.Where("room_id = ?", repairReq.RoomID).
			First(&resident).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No resident found in this room"})
			return
		}
	}

	repair := models.RepairRequest{
		ResidentID:  resident.ID,
		RoomID:      repairReq.RoomID,
		Title:       repairReq.Title,
		Description: repairReq.Description,
		Priority:    repairReq.Priority,
		Status:      "pending",
		ReportedAt:  time.Now(),
	}

	if err := h.DB.Create(&repair).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, repair)
}

func (h *RepairHandler) UpdateRepairStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repair ID"})
		return
	}

	var req struct {
		Status          string `json:"status" binding:"required,oneof=pending in_progress completed cancelled"`
		AssignedTo      *uint  `json:"assigned_to"`
		TechnicianNotes string `json:"technician_notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var repair models.RepairRequest
	if err := h.DB.First(&repair, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repair request not found"})
		return
	}

	updateData := map[string]interface{}{
		"status":           req.Status,
		"technician_notes": req.TechnicianNotes,
	}

	if req.AssignedTo != nil {
		updateData["assigned_to"] = req.AssignedTo
	}

	if req.Status == "completed" && repair.CompletedAt == nil {
		now := time.Now()
		updateData["completed_at"] = now
	}

	if err := h.DB.Model(&repair).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repair)
}

func (h *RepairHandler) GetRepairStats(c *gin.Context) {
	var stats struct {
		Total        int64 `json:"total"`
		Pending      int64 `json:"pending"`
		InProgress   int64 `json:"in_progress"`
		Completed    int64 `json:"completed"`
		HighPriority int64 `json:"high_priority"`
	}

	h.DB.Model(&models.RepairRequest{}).Count(&stats.Total)
	h.DB.Model(&models.RepairRequest{}).Where("status = ?", "pending").Count(&stats.Pending)
	h.DB.Model(&models.RepairRequest{}).Where("status = ?", "in_progress").Count(&stats.InProgress)
	h.DB.Model(&models.RepairRequest{}).Where("status = ?", "completed").Count(&stats.Completed)
	h.DB.Model(&models.RepairRequest{}).Where("priority = ?", "high").Count(&stats.HighPriority)

	c.JSON(http.StatusOK, stats)
}
