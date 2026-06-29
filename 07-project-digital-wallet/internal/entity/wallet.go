package entity

import "time"

type Wallet struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	WalletNumber string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"wallet_number"`
	Balance      float64   `gorm:"type:numeric(15,2);default:0.00;not null" json:"balance"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relasi GORM (Preloaded)
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}
