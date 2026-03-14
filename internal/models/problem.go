// Package models
package models

import "time"

type Problem struct {
	ID             uint      `gorm:"primaryKey"`
	Slug           string    `gorm:"type:varchar(255);unique;not null"`
	AcceptanceRate float64   `gorm:"type:decimal(5,2)"`
	Title          string    `gorm:"type:varchar(255);not null"`
	Difficulty     string    `gorm:"type:varchar(50);not null"`
	ProblemMd      string    `gorm:"type:text;not null"`
	TimeLimitMs    int
	MemoryLimitKb  int
	IsPublished    bool      `gorm:"default:false"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy      uint
}
