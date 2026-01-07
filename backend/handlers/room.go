package handlers

import (
	"dormitory-management/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomHandler struct {
	DB *gorm.DB
}

func NewRoomHandler(db *gorm.DB) *RoomHandler {
	return &RoomHandler{DB: db}
}

func (h *RoomHandler) GetRooms(c *gin.Context) {
	var rooms []models.Room
	query := h.DB

	// Filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by building
	if building := c.Query("building"); building != "" {
		query = query.Where("building = ?", building)
	}

	if err := query.Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) GetRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	var room models.Room
	if err := h.DB.Preload("Residents.User").
		First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	c.JSON(http.StatusOK, room)
}

// CreateRoom - Menambahkan kamar baru
func (h *RoomHandler) CreateRoom(c *gin.Context) {
    var roomReq models.RoomRequest
    if err := c.ShouldBindJSON(&roomReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid: " + err.Error()})
        return
    }

    // Inisialisasi model Room dari request
    room := models.Room{
        RoomNumber:      roomReq.RoomNumber,
        Building:        roomReq.Building,
        Floor:           roomReq.Floor,
        Capacity:        roomReq.Capacity,
        CurrentOccupancy: 0, // Default kamar baru adalah kosong
        MonthlyRate:     roomReq.MonthlyRate,
        Description:     roomReq.Description,
        Status:          "available", // Default status
    }

    // Simpan ke database
    if err := h.DB.Create(&room).Error; err != nil {
        // Cek jika nomor kamar duplikat (karena unique constraint di schema)
        c.JSON(http.StatusConflict, gin.H{"error": "Gagal menambahkan kamar. Nomor kamar mungkin sudah ada."})
        return
    }

    c.JSON(http.StatusCreated, room)
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	var room models.Room
	if err := h.DB.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Model(&room).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID Kamar tidak valid"})
        return
    }

    var room models.Room
    // 1. Cari kamar terlebih dahulu
    if err := h.DB.First(&room, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Kamar tidak ditemukan"})
        return
    }

    // 2. Validasi: Jika current_occupancy > 0, jangan izinkan hapus
    if room.CurrentOccupancy > 0 {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "Kamar tidak dapat dihapus karena masih memiliki penghuni (" + 
                     strconv.Itoa(room.CurrentOccupancy) + " orang).",
        })
        return
    }

    // 3. Jalankan penghapusan
    if err := h.DB.Delete(&room).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kamar: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Kamar berhasil dihapus"})
}

func (h *RoomHandler) GetAvailableRooms(c *gin.Context) {
	var rooms []models.Room
	if err := h.DB.Where("status = ? AND current_occupancy < capacity", "available").
		Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) AssignResident(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	var req struct {
		ResidentID uint `json:"resident_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx := h.DB.Begin()

	// Check room availability
	var room models.Room
	if err := tx.First(&room, roomID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	if room.Status != "available" || room.CurrentOccupancy >= room.Capacity {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room is not available"})
		return
	}

	// Check resident
	var resident models.Resident
	if err := tx.First(&resident, req.ResidentID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Resident not found"})
		return
	}

	// Update resident's room
	if err := tx.Model(&resident).Update("room_id", roomID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update room occupancy
	if err := tx.Model(&room).Updates(map[string]interface{}{
		"current_occupancy": room.CurrentOccupancy + 1,
		"status":            "occupied",
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Resident assigned successfully"})
}

func (h *RoomHandler) GetMyRoom(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var resident models.Resident
	if err := h.DB.Preload("Room").Preload("Room.Residents.User").
		Where("user_id = ?", userID).First(&resident).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	if resident.RoomID == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "You are not assigned to any room yet"})
		return
	}

	c.JSON(http.StatusOK, resident.Room)
}
