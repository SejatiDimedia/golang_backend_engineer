package entity

import "time"

type Desk struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(50);not null" json:"name" binding:"required"`
	Type      string    `gorm:"type:varchar(50);default:hot-desk;not null" json:"type" binding:"required,oneof=hot-desk meeting-room"`
	IsActive  bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
