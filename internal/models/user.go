// Package models
package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password  string    `gorm:"type:text" json:"password"`
	Username  string    `gorm:"type:varchar(255)" json:"username"`
	FullName  string    `gorm:"type:varchar(255)" json:"fullName"`
	Role      string    `gorm:"default:'student'" json:"role"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
}
