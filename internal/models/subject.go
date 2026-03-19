package models

import (
	"time"

	"gorm.io/gorm"
)

// Subject represents a subject/course in the competitive programming platform
type Subject struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;size:255"`
	Slug        string         `json:"slug" gorm:"uniqueIndex;not null;size:255"`
	Description string         `json:"description" gorm:"type:text"`
	Color       string           `json:"color" gorm:"size:7;comment:Hex color code for frontend display (e.g., #FF0000)"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt   `json:"-" gorm:"index"`
	Sessions    []SubjectSession `json:"sessions" gorm:"foreignKey:SubjectID"`
}
