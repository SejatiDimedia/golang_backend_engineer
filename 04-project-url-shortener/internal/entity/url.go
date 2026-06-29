package entity

import "time"

type URL struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	LongURL     string     `gorm:"type:text;not null" json:"long_url"`
	ShortCode   string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"short_code"`
	ClickCount  int64      `gorm:"default:0;not null" json:"click_count"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// IsExpired mengecek apakah URL telah melewati masa berlakunya
func (u *URL) IsExpired() bool {
	if u.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*u.ExpiresAt)
}
