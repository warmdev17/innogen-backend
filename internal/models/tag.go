// Package models
package models

type Tag struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(255);unique;not null" json:"name"`
}

type ProblemTag struct {
	ProblemID uint `gorm:"primaryKey" json:"problemId"`
	TagID     uint `gorm:"primaryKey" json:"tagId"`
}
