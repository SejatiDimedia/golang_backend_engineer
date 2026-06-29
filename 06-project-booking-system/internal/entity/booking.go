package entity

import "time"

type Booking struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	DeskID    uint      `gorm:"not null" json:"desk_id"`
	StartTime time.Time `gorm:"type:timestamp;not null" json:"start_time"`
	EndTime   time.Time `gorm:"type:timestamp;not null" json:"end_time"`
	Status    string    `gorm:"type:varchar(20);default:CONFIRMED;not null" json:"status"` // "CONFIRMED" atau "CANCELLED"
	CreatedAt time.Time `json:"created_at"`
	
	// Relasi GORM (Preloaded)
	User      *User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Desk      *Desk     `gorm:"foreignKey:DeskID;constraint:OnDelete:RESTRICT" json:"desk,omitempty"`
}
