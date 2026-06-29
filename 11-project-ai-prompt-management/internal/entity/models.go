package entity

import (
	"time"
)

type Workspace struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Prompts   []Prompt  `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"-"`
	ApiKeys   []ApiKey  `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"-"`
}

type WorkspaceMember struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	WorkspaceID uint      `gorm:"not null" json:"workspace_id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	Role        string    `gorm:"type:varchar(50);default:'member';not null" json:"role"` // 'admin' / 'member'
	CreatedAt   time.Time `json:"created_at"`

	Workspace   Workspace `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"-"`
}

type Prompt struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	WorkspaceID uint      `gorm:"not null" json:"workspace_id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Workspace   Workspace       `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"-"`
	Versions    []PromptVersion `gorm:"foreignKey:PromptID;constraint:OnDelete:CASCADE" json:"versions,omitempty"`
}

type PromptVersion struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	PromptID      uint      `gorm:"not null;index" json:"prompt_id"`
	VersionNumber int       `gorm:"not null" json:"version_number"`
	PromptText    string    `gorm:"type:text;not null" json:"prompt_text"`
	Status        string    `gorm:"type:varchar(20);default:'DRAFT';not null" json:"status"` // 'DRAFT', 'ACTIVE'
	CreatedAt     time.Time `json:"created_at"`

	Prompt        Prompt    `gorm:"foreignKey:PromptID;constraint:OnDelete:CASCADE" json:"-"`
}

type ApiKey struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	WorkspaceID uint      `gorm:"not null;index" json:"workspace_id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	KeyHash     string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"-"` // SHA-256 hash
	MaskedKey   string    `gorm:"type:varchar(50);not null" json:"masked_key"`   // prompt_live_xxxx...1234
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`

	Workspace   Workspace `gorm:"foreignKey:WorkspaceID;constraint:OnDelete:CASCADE" json:"-"`
}

type AnalyticsLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ApiKeyID      uint      `gorm:"not null;index" json:"api_key_id"`
	PromptID      uint      `gorm:"not null;index" json:"prompt_id"`
	LatencyMs     int64     `gorm:"not null" json:"latency_ms"`
	TokenEstimate int       `gorm:"not null" json:"token_estimate"`
	ResponseCode  int       `gorm:"not null" json:"response_code"`
	CreatedAt     time.Time `json:"created_at"`

	ApiKey        ApiKey    `gorm:"foreignKey:ApiKeyID;constraint:OnDelete:CASCADE" json:"-"`
	Prompt        Prompt    `gorm:"foreignKey:PromptID;constraint:OnDelete:CASCADE" json:"-"`
}
