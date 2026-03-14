// Package models
package models

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Language       string    `gorm:"type:varchar(50);not null"`
	Code           string    `gorm:"type:text;not null"`
	Status         string    `gorm:"type:varchar(50);not null"`
	RuntimeMs      int
	MemoryKb       int
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ErrorMessage   string    `gorm:"type:text"`
	PassCount      int
	TotalTestcases int
	UserID         uint
	ProblemID      uint
}
