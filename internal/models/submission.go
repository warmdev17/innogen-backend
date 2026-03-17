// Package models
package models

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Language       string    `gorm:"type:varchar(50);not null" json:"language"`
	Code           string    `gorm:"type:text;not null" json:"code"`
	Status         string    `gorm:"type:varchar(50);not null" json:"status"`
	RuntimeMs      int       `json:"runtimeMs"`
	MemoryKb       int       `json:"memoryKb"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	ErrorMessage   string    `gorm:"type:text" json:"errorMessage"`
	PassCount      int       `json:"passCount"`
	TotalTestcases int       `json:"totalTestcases"`
	UserID         uint      `json:"userId"`
	ProblemID      uint      `json:"problemId"`
}
