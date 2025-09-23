package models

import (
	"gorm.io/gorm"
	"time"
)

type Article struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Category    string         `json:"category" gorm:"not null;index"`
	Subcategory string         `json:"subcategory" gorm:"not null;index"`
	Title       string         `json:"title" gorm:"not null"`
	Abstract    string         `json:"abstract" gorm:"type:text"`
	Link        string         `json:"link" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
