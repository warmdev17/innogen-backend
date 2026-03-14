// Package models
package models

import "time"

type Testcase struct {
	ID             uint      `gorm:"primaryKey"`
	InputData      string    `gorm:"type:text;not null"`
	ExpectedOutput string    `gorm:"type:text;not null"`
	IsHidden       bool      `gorm:"default:true"`
	Role           string    `gorm:"type:varchar(50)"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ProblemID      uint
}
