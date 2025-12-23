package handlers

import (
	"dormitory-management/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ResidentHandler struct {
	DB *gorm.DB
}

func NewResidentHandler(db *gorm.DB) *ResidentHandler {
	return &ResidentHandler{DB: db}
}

func (h *ResidentHandler) GetResidents(c *gin.Context) {
	var residents []models.Resident
	query := h.DB.Preload("User").Preload("Room")

	if name := c.Query("name"); name != "" {
		query = query.Joins("JOIN users ON users.id = residents.user_id").Where("users.name LIKE ?", "%"+name+"%")
	}

	if err := query.Find(&residents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch residents"})
		return
	}

	c.JSON(http.StatusOK, residents)
}

func (h *ResidentHandler) GetResident(c *gin.Context) {
	id := c.Param("id")
	var resident models.Resident
	if err := h.DB.Preload("User").Preload("Room").First(&resident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resident not found"})
		return
	}
	c.JSON(http.StatusOK, resident)
}

type CreateResidentInput struct {
	Name             string `json:"name" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required"`
	Phone            string `json:"phone"`
	StudentID        string `json:"student_id" binding:"required"`
	Faculty          string `json:"faculty"`
	Major            string `json:"major"`
	YearOfEntry      int    `json:"year_of_entry"`
	EmergencyContact string `json:"emergency_contact"`
	RoomID           *uint  `json:"room_id"`
}

func (h *ResidentHandler) CreateResident(c *gin.Context) {
	var input CreateResidentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()

	// 1. Create User
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "student",
		Phone:    input.Phone,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	// 2. Create Resident
	resident := models.Resident{
		UserID:           user.ID,
		StudentID:        input.StudentID,
		Faculty:          input.Faculty,
		Major:            input.Major,
		YearOfEntry:      input.YearOfEntry,
		EmergencyContact: input.EmergencyContact,
		RoomID:           input.RoomID,
		Status:           "active",
	}

	if err := tx.Create(&resident).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create resident: " + err.Error()})
		return
	}

	// 3. Update Room Occupancy
	if input.RoomID != nil {
		var room models.Room
		if err := tx.First(&room, *input.RoomID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room not found"})
			return
		}

		if room.CurrentOccupancy >= room.Capacity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room is full"})
			return
		}

		room.CurrentOccupancy++
		if room.CurrentOccupancy >= room.Capacity {
			room.Status = "occupied"
		}

		if err := tx.Save(&room).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusCreated, resident)
}

type UpdateResidentInput struct {
	Name             string `json:"name"`
	Phone            string `json:"phone"`
	Faculty          string `json:"faculty"`
	Major            string `json:"major"`
	EmergencyContact string `json:"emergency_contact"`
	Status           string `json:"status"`
	RoomID           *uint  `json:"room_id"`
}

func (h *ResidentHandler) UpdateResident(c *gin.Context) {
	id := c.Param("id")
	var resident models.Resident
	if err := h.DB.Preload("User").First(&resident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resident not found"})
		return
	}

	var input UpdateResidentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()

	// Update User info
	if input.Name != "" || input.Phone != "" {
		if input.Name != "" {
			resident.User.Name = input.Name
		}
		if input.Phone != "" {
			resident.User.Phone = input.Phone
		}
		if err := tx.Save(&resident.User).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user info"})
			return
		}
	}

	// Update Resident info
	resident.Faculty = input.Faculty
	resident.Major = input.Major
	resident.EmergencyContact = input.EmergencyContact
	if input.Status != "" {
		resident.Status = input.Status
	}

	// Move Room Logic
	if input.RoomID != nil && (resident.RoomID == nil || *input.RoomID != *resident.RoomID) {
		// Decrement old room
		if resident.RoomID != nil {
			var oldRoom models.Room
			if err := tx.First(&oldRoom, *resident.RoomID).Error; err == nil {
				if oldRoom.CurrentOccupancy > 0 {
					oldRoom.CurrentOccupancy--
				}
				if oldRoom.CurrentOccupancy < oldRoom.Capacity {
					oldRoom.Status = "available"
				}
				tx.Save(&oldRoom)
			}
		}

		// Increment new room
		var newRoom models.Room
		if err := tx.First(&newRoom, *input.RoomID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "New Room not found"})
			return
		}

		if newRoom.CurrentOccupancy >= newRoom.Capacity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "New Room is full"})
			return
		}

		newRoom.CurrentOccupancy++
		if newRoom.CurrentOccupancy >= newRoom.Capacity {
			newRoom.Status = "occupied"
		}
		tx.Save(&newRoom)

		resident.RoomID = input.RoomID
	} else if input.RoomID == nil && resident.RoomID != nil {
		// Remove from room
		var oldRoom models.Room
		if err := tx.First(&oldRoom, *resident.RoomID).Error; err == nil {
			if oldRoom.CurrentOccupancy > 0 {
				oldRoom.CurrentOccupancy--
			}
			if oldRoom.CurrentOccupancy < oldRoom.Capacity {
				oldRoom.Status = "available"
			}
			tx.Save(&oldRoom)
		}
		resident.RoomID = nil
	}

	if err := tx.Save(&resident).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resident"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, resident)
}

func (h *ResidentHandler) DeleteResident(c *gin.Context) {
	id := c.Param("id")
	var resident models.Resident
	if err := h.DB.First(&resident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resident not found"})
		return
	}

	tx := h.DB.Begin()

	// Update room occupancy
	if resident.RoomID != nil {
		var room models.Room
		if err := tx.First(&room, *resident.RoomID).Error; err == nil {
			if room.CurrentOccupancy > 0 {
				room.CurrentOccupancy--
			}
			if room.CurrentOccupancy < room.Capacity {
				room.Status = "available"
			}
			tx.Save(&room)
		}
	}

	// Delete user (cascade will handle resident if configured, but let's be safe)
	if err := tx.Delete(&models.User{}, resident.UserID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	// Note: GORM might soft delete if User has DeletedAt. If not, hard delete.
	// Residents table has ON DELETE CASCADE on user_id, so deleting user deletes resident.

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Resident deleted successfully"})
}
