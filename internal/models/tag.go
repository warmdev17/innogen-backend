// Package models
package models

type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(255);unique;not null"`
}

type ProblemTag struct {
	ProblemID uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`
}
