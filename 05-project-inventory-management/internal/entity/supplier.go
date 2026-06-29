package entity

import "time"

type Supplier struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(150);not null" json:"name" binding:"required"`
	ContactName string    `gorm:"type:varchar(100)" json:"contact_name"`
	Email       string    `gorm:"type:varchar(100)" json:"email" binding:"omitempty,email"`
	Phone       string    `gorm:"type:varchar(30)" json:"phone"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
