// Package models
package models

import "time"

type Submission struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	ProblemID uint

	Code      string `gorm:"type:text"`
	Language  string // cpp, go, py
	Status    string // pending, accepted, wrong_answer, runtime_error
	RuntimeMs int

	CreatedAt time.Time
}
