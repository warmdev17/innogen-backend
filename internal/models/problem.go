// Package models
package models

import "time"

type Problem struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	Slug           string     `gorm:"type:varchar(255);unique;not null" json:"slug"`
	AcceptanceRate float64    `gorm:"type:decimal(5,2)" json:"acceptanceRate"`
	Title          string     `gorm:"type:varchar(255);not null" json:"title"`
	Difficulty     string     `gorm:"type:varchar(50);not null" json:"difficulty"`
	ProblemMd      string     `gorm:"type:text;not null" json:"problemMd"`
	TimeLimitMs    int        `json:"timeLimitMs"`
	MemoryLimitKb  int        `json:"memoryLimitKb"`
	IsPublished    bool       `gorm:"default:false" json:"isPublished"`
	CreatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
	CreatedBy      uint       `json:"createdBy"`
	Testcases      []Testcase `gorm:"foreignKey:ProblemID" json:"testcases,omitempty"`
}

