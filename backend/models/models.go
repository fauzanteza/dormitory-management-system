package models

type Resident struct {
	ID               uint            `json:"id" gorm:"primaryKey"`
	UserID           uint            `json:"user_id" gorm:"unique;not null"`
	RoomID           *uint           `json:"room_id"`
	StudentID        string          `json:"student_id" gorm:"unique"`
	Faculty          string          `json:"faculty"`
	Major            string          `json:"major"`
	YearOfEntry      int             `json:"year_of_entry"`
	EmergencyContact string          `json:"emergency_contact"`
	Status           string          `json:"status" gorm:"type:enum('active','inactive','graduated');default:'active'"`
	User             *User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Room             *Room           `json:"room,omitempty" gorm:"foreignKey:RoomID"`
	Payments         []Payment       `json:"payments,omitempty" gorm:"foreignKey:ResidentID"`
	RepairRequests   []RepairRequest `json:"repair_requests,omitempty" gorm:"foreignKey:ResidentID"`
}

type DashboardStats struct {
	TotalRooms      int64            `json:"total_rooms"`
	OccupiedRooms   int64            `json:"occupied_rooms"`
	AvailableRooms  int64            `json:"available_rooms"`
	TotalResidents  int64            `json:"total_residents"`
	PendingPayments int64            `json:"pending_payments"`
	TotalRevenue    float64          `json:"total_revenue"`
	PendingRepairs  int64            `json:"pending_repairs"`
	MonthlyRevenue  []MonthlyRevenue `json:"monthly_revenue"`
	RoomStatus      []RoomStatus     `json:"room_status"`
}

type MonthlyRevenue struct {
	Month   string  `json:"month"`
	Revenue float64 `json:"revenue"`
}

type RoomStatus struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}
