package entity

import "time"

type Notification struct {
	ID           uint              `gorm:"primaryKey" json:"id"`
	Type         string            `gorm:"type:varchar(20);not null" json:"type"` // email, webhook, push
	Target       string            `gorm:"type:varchar(255);not null" json:"target"`
	Content      string            `gorm:"type:text;not null" json:"content"`
	Status       string            `gorm:"type:varchar(20);default:'PENDING';not null" json:"status"` // PENDING, PROCESSING, SENT, FAILED
	MaxRetries   int               `gorm:"default:5;not null" json:"max_retries"`
	AttemptCount int               `gorm:"default:0;not null" json:"attempt_count"`
	SendAt       time.Time         `gorm:"not null" json:"send_at"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Logs         []NotificationLog `gorm:"foreignKey:NotificationID;constraint:OnDelete:CASCADE" json:"logs,omitempty"`
}

type NotificationLog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	NotificationID uint      `gorm:"not null" json:"notification_id"`
	Attempt        int       `gorm:"not null" json:"attempt"`
	Status         string    `gorm:"type:varchar(20);not null" json:"status"` // SENT, FAILED
	ErrorMessage   string    `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}
