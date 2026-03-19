// Package models
package models

import "time"

type SubjectSession struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SubjectID  uint      `gorm:"not null" json:"subjectId"`
	Title      string    `gorm:"type:varchar(255);not null" json:"title"`
	OrderIndex int       `gorm:"not null" json:"orderIndex"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	Lessons    []Lesson  `gorm:"foreignKey:SubjectSessionID" json:"lessons"`
}

type Lesson struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	SubjectSessionID uint            `gorm:"not null" json:"subjectSessionId"`
	Title            string          `gorm:"type:varchar(255);not null" json:"title"`
	OrderIndex       int             `gorm:"not null" json:"orderIndex"`
	CreatedAt        time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	Problems         []LessonProblem `gorm:"foreignKey:LessonID" json:"problems"`
}

type LessonProblem struct {
	LessonID   uint `gorm:"primaryKey" json:"lessonId"`
	ProblemID  uint `gorm:"primaryKey" json:"problemId"`
	OrderIndex int  `gorm:"not null" json:"orderIndex"`
}
