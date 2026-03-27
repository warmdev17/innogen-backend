// Package models
package models

import "time"

// Subject represents a subject/course in the competitive programming platform
type Subject struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	Name        string           `gorm:"type:varchar(255)" json:"name"`
	Title       string           `gorm:"type:varchar(255);not null" json:"title"`
	Slug        string           `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description string           `gorm:"type:text" json:"description"`
	Color       string           `gorm:"type:varchar(7);comment:Hex color code for frontend display (e.g., #FF0000)" json:"color"`
	Language    string           `gorm:"type:varchar(50);default:'en'" json:"language"`
	IconName    string           `gorm:"type:varchar(50)" json:"iconName"`
	CreatedAt   time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	Sessions    []SubjectSession `gorm:"foreignKey:SubjectID" json:"sessions,omitempty"`
}

type SubjectSession struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	SubjectID     uint      `gorm:"not null" json:"subjectId"`
	Title         string    `gorm:"type:varchar(255);not null" json:"title"`
	OrderIndex    int       `gorm:"not null" json:"orderIndex"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	Lessons       []Lesson  `gorm:"foreignKey:SubjectSessionID" json:"lessons"`
	SubjectId     *uint     `gorm:"-" json:"-"` // Alias for subjectId field name compatibility
	SessionId     *uint     `gorm:"-" json:"-"` // Alias for session ID
}

type Lesson struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	SubjectSessionID uint            `gorm:"not null" json:"subjectSessionId"`
	SessionID        uint            `gorm:"-" json:"sessionId"` // Alias for frontend compatibility
	Title            string          `gorm:"type:varchar(255);not null" json:"title"`
	ContentMd        string          `gorm:"type:text" json:"contentMd"`
	OrderIndex       int             `gorm:"not null" json:"orderIndex"`
	CreatedAt        time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	Problems         []LessonProblem `gorm:"foreignKey:LessonID" json:"problems"`
}

type LessonProblem struct {
	LessonID   uint     `gorm:"primaryKey" json:"lessonId"`
	ProblemID  uint     `gorm:"primaryKey" json:"problemId"`
	OrderIndex int      `gorm:"not null" json:"orderIndex"`
	Problem    *Problem `gorm:"foreignKey:ProblemID" json:"problem,omitempty"`
}

// ─── Response DTOs (for frontend display) ────────────────────────────────────

// SubjectListItem is used for GET /subjects (list view)
type SubjectListItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name,omitempty"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Language    string `json:"language"`
	IconName    string `json:"iconName"`
}

// LessonItem is used inside session/subject detail responses
type LessonItem struct {
	ID         uint   `json:"id"`
	SessionID  uint   `json:"sessionId,omitempty"`
	Title      string `json:"title"`
	ContentMd  string `json:"contentMd,omitempty"`
	OrderIndex int    `json:"orderIndex"`
}

// SessionItem is used inside subject detail response
type SessionItem struct {
	ID         uint         `json:"id"`
	SubjectID  uint         `json:"subjectId,omitempty"`
	Title      string       `json:"title"`
	OrderIndex int          `json:"orderIndex"`
	Lessons    []LessonItem `json:"lessons"`
}

// SubjectDetailResponse is used for GET /subjects/:slug
type SubjectDetailResponse struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name,omitempty"`
	Title       string        `json:"title"`
	Slug        string        `json:"slug"`
	Description string        `json:"description"`
	Color       string        `json:"color"`
	Language    string        `json:"language"`
	IconName    string        `json:"iconName"`
	Sessions    []SessionItem `json:"sessions"`
}

// SessionResponse is used for GET /sessions/:id
type SessionResponse struct {
	ID         uint         `json:"id"`
	SubjectID  uint         `json:"subjectId"`
	Title      string       `json:"title"`
	OrderIndex int          `json:"orderIndex"`
	Lessons    []LessonItem `json:"lessons"`
}

// LessonProblemItem is a lightweight problem entry inside a lesson
type LessonProblemItem struct {
	ID             uint    `json:"id"`
	Slug           string  `json:"slug"`
	Title          string  `json:"title"`
	Difficulty     string  `json:"difficulty"`
	AcceptanceRate float64 `json:"acceptanceRate"`
	OrderIndex     int     `json:"orderIndex"`
}

// LessonResponse is used for GET /lessons/:id
type LessonResponse struct {
	ID         uint                `json:"id"`
	SessionID  uint                `json:"sessionId"`
	Title      string              `json:"title"`
	ContentMd  string              `json:"contentMd"`
	OrderIndex int                 `json:"orderIndex"`
	Problems   []LessonProblemItem `json:"problems"`
}
