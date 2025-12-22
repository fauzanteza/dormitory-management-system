package models

import (
	"time"
)

type Payment struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	ResidentID    uint       `json:"resident_id" gorm:"not null"`
	RoomID        uint       `json:"room_id" gorm:"not null"`
	Month         time.Time  `json:"month" gorm:"not null"`
	Amount        float64    `json:"amount" gorm:"type:decimal(10,2);not null"`
	PaymentDate   *time.Time `json:"payment_date"`
	Status        string     `json:"status" gorm:"type:enum('paid','pending','overdue');default:'pending'"`
	PaymentMethod string     `json:"payment_method"`
	ReceiptNumber string     `json:"receipt_number"`
	Notes         string     `json:"notes"`
	Resident      *Resident  `json:"resident,omitempty" gorm:"foreignKey:ResidentID"`
	Room          *Room      `json:"room,omitempty" gorm:"foreignKey:RoomID"`
	CreatedAt     time.Time  `json:"created_at"`
}

type PaymentRequest struct {
	ResidentID    uint    `json:"resident_id" binding:"required"`
	Month         string  `json:"month" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	PaymentMethod string  `json:"payment_method"`
	ReceiptNumber string  `json:"receipt_number"`
	Notes         string  `json:"notes"`
}
