package models

import (
	"time"
)

type Booking struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	UserID         uint      `json:"user_id" gorm:"not null"`
	RoomID         uint      `json:"room_id" gorm:"not null"`
	StartDate      time.Time `json:"start_date" gorm:"type:date;not null"`
	DurationMonths int       `json:"duration_months" gorm:"default:12"`
	Notes          string    `json:"notes"`
	Status         string    `json:"status" gorm:"type:enum('pending','approved','rejected','cancelled');default:'pending'"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	User User `json:"user" gorm:"foreignKey:UserID"`
	Room Room `json:"room" gorm:"foreignKey:RoomID"`
}
