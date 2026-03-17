// Package models
package models

import "time"

type Testcase struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	InputData      string    `gorm:"type:text;not null" json:"inputData"`
	ExpectedOutput string    `gorm:"type:text;not null" json:"expectedOutput"`
	IsHidden       bool      `gorm:"default:true" json:"isHidden"`
	Role           string    `gorm:"type:varchar(50)" json:"role"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	ProblemID      uint      `json:"problemId"`
}
