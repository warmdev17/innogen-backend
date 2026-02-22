// Package models
package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string `gorm:"unique"`
	Password  string
	Role      string `gorm:"default:'student'"` // student, teacher, admin (tui)
	FullName  string
	CreatedAt time.Time
}
