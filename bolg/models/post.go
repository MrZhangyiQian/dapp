package models

import (
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
	UserID  uint
}