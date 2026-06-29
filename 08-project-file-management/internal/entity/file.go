package entity

import "time"

type File struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	FileName    string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FileSize    int64     `gorm:"not null" json:"file_size"`
	ContentType string    `gorm:"type:varchar(100);not null" json:"content_type"`
	ObjectKey   string    `gorm:"type:varchar(500);uniqueIndex;not null" json:"object_key"`
	Status      string    `gorm:"type:varchar(20);default:'PENDING';not null" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
