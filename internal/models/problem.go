// Package models
package models

import "time"

type Problem struct {
	ID            uint   `gorm:"primaryKey"`
	Title         string `gorm:"not null"`
	Description   string `gorm:"type:text"`
	Difficulty    string `gorm:"default:'easy'"`
	TimeLimitMs   int    `gorm:"default:1000"`
	MemoryLimitMB int    `gorm:"default:256"`
	CreatedBy     uint
	CreatedAt     time.Time
}
