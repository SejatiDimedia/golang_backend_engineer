package entity

import "time"

type Product struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"type:varchar(150);not null" json:"name" binding:"required"`
	SKU           string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"sku" binding:"required"`
	Description   string    `gorm:"type:text" json:"description"`
	Price         float64   `gorm:"type:decimal(12,2);not null" json:"price" binding:"required,gt=0"`
	StockQuantity int64     `gorm:"default:0;not null" json:"stock_quantity"`
	CategoryID    uint      `gorm:"not null" json:"category_id" binding:"required"`
	SupplierID    uint      `gorm:"not null" json:"supplier_id" binding:"required"`
	
	// Relasi GORM (Preload)
	Category      *Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT" json:"category,omitempty"`
	Supplier      *Supplier `gorm:"foreignKey:SupplierID;constraint:OnDelete:RESTRICT" json:"supplier,omitempty"`
	
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
