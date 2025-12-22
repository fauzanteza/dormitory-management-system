package models

import (
	"time"
)

type Room struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	RoomNumber       string     `json:"room_number" gorm:"unique;not null"`
	Building         string     `json:"building"`
	Floor            int        `json:"floor"`
	Capacity         int        `json:"capacity" gorm:"default:2"`
	CurrentOccupancy int        `json:"current_occupancy" gorm:"default:0"`
	Status           string     `json:"status" gorm:"type:enum('available','occupied','maintenance');default:'available'"`
	MonthlyRate      float64    `json:"monthly_rate" gorm:"type:decimal(10,2);not null"`
	Description      string     `json:"description"`
	Residents        []Resident `json:"residents,omitempty" gorm:"foreignKey:RoomID"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type RoomRequest struct {
	RoomNumber  string  `json:"room_number" binding:"required"`
	Building    string  `json:"building" binding:"required"`
	Floor       int     `json:"floor" binding:"required"`
	Capacity    int     `json:"capacity" binding:"required"`
	MonthlyRate float64 `json:"monthly_rate" binding:"required"`
	Description string  `json:"description"`
}
