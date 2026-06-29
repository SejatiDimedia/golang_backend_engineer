package entity

import "time"

type Transaction struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	SourceWalletID       *uint     `gorm:"index" json:"source_wallet_id,omitempty"`
	DestinationWalletID  *uint     `gorm:"index" json:"destination_wallet_id,omitempty"`
	Amount               float64   `gorm:"type:numeric(15,2);not null" json:"amount"`
	Type                 string    `gorm:"type:varchar(20);not null" json:"type"` // "top-up", "withdraw", "transfer"
	Description          string    `gorm:"type:text" json:"description"`
	CreatedAt            time.Time `json:"created_at"`

	// Relasi Optional
	SourceWallet      *Wallet `gorm:"foreignKey:SourceWalletID;constraint:OnDelete:SET NULL" json:"source_wallet,omitempty"`
	DestinationWallet *Wallet `gorm:"foreignKey:DestinationWalletID;constraint:OnDelete:SET NULL" json:"destination_wallet,omitempty"`
}
