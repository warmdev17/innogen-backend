// Package models
package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"type:varchar(255);unique;not null"`
	Password  string    `gorm:"type:text"`
	Username  string    `gorm:"type:varchar(255)"`
	FullName  string    `gorm:"type:varchar(255)"`
	Role      string    `gorm:"default:'student'"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
