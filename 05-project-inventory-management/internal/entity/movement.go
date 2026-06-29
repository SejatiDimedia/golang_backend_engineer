package entity

import "time"

type StockMovement struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	Type      string    `gorm:"type:varchar(10);not null" json:"type"` // "IN" atau "OUT"
	Quantity  int64     `gorm:"not null" json:"quantity"`
	Reference string    `gorm:"type:varchar(100)" json:"reference,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	
	// Relasi GORM
	Product   *Product  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product,omitempty"`
}
