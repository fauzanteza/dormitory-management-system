package models

import (
	"time"
)

type RepairRequest struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	ResidentID      uint       `json:"resident_id" gorm:"not null"`
	RoomID          uint       `json:"room_id" gorm:"not null"`
	Title           string     `json:"title" gorm:"not null"`
	Description     string     `json:"description" gorm:"not null"`
	Priority        string     `json:"priority" gorm:"type:enum('low','medium','high');default:'medium'"`
	Status          string     `json:"status" gorm:"type:enum('pending','in_progress','completed','cancelled');default:'pending'"`
	ReportedAt      time.Time  `json:"reported_at"`
	CompletedAt     *time.Time `json:"completed_at"`
	AssignedTo      *uint      `json:"assigned_to"`
	TechnicianNotes string     `json:"technician_notes"`
	Resident        *Resident  `json:"resident,omitempty" gorm:"foreignKey:ResidentID"`
	Room            *Room      `json:"room,omitempty" gorm:"foreignKey:RoomID"`
	AssignedUser    *User      `json:"assigned_user,omitempty" gorm:"foreignKey:AssignedTo"`
}

type RepairRequestInput struct {
	RoomID      uint   `json:"room_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Priority    string `json:"priority" binding:"required,oneof=low medium high"`
}
