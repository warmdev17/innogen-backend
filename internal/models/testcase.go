// Package models
package models

type Testcase struct {
	ID        uint `gorm:"primaryKey"`
	ProblemID uint
	Input     string `gorm:"type:text"`
	Output    string `gorm:"type:text"`
}
