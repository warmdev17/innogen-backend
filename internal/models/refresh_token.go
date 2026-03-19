package models

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a refresh token record
type RefreshToken struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"userId"`
	TokenHash  string    `gorm:"type:text;not null" json:"-"`
	ExpiresAt  time.Time `gorm:"not null;index" json:"expiresAt"`
	Revoked    bool      `gorm:"default:false" json:"revoked"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

// TableName specifies the table name for RefreshToken
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
