// Package models
package models

import "time"

type Subject struct {
	ID          uint             `gorm:"primaryKey"`
	Title       string           `gorm:"type:varchar(255);not null"`
	Slug        string           `gorm:"type:varchar(255);unique;not null"`
	IsPublished bool             `gorm:"default:false"`
	CreatedAt   time.Time        `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time        `gorm:"default:CURRENT_TIMESTAMP"`
	Sessions    []SubjectSession `gorm:"foreignKey:SubjectID"`
}

type SubjectSession struct {
	ID        uint      `gorm:"primaryKey"`
	SubjectID uint      `gorm:"not null"`
	Title     string    `gorm:"type:varchar(255);not null"`
	OrderIndex int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Lessons   []Lesson  `gorm:"foreignKey:SubjectSessionID"`
}

type Lesson struct {
	ID               uint            `gorm:"primaryKey"`
	SubjectSessionID uint            `gorm:"not null"`
	Title            string          `gorm:"type:varchar(255);not null"`
	OrderIndex       int             `gorm:"not null"`
	CreatedAt        time.Time       `gorm:"default:CURRENT_TIMESTAMP"`
	Problems         []LessonProblem `gorm:"foreignKey:LessonID"`
}

type LessonProblem struct {
	LessonID   uint `gorm:"primaryKey"`
	ProblemID  uint `gorm:"primaryKey"`
	OrderIndex int  `gorm:"not null"`
}
