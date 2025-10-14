package models

import (
	"time"
)

type Article struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Category    string     `json:"category" gorm:"not null;index"`
	Subcategory string     `json:"subcategory" gorm:"not null;index"`
	Title       string     `json:"title" gorm:"not null"`
	Abstract    string     `json:"abstract" gorm:"type:text"`
	Link        string     `json:"link" gorm:"not null;uniqueIndex"`
	RSSName     *string    `json:"rss_name" gorm:"index,not null"`
	RSS         RSSSource  `json:"rss" gorm:"foreignKey:RSSName;references:Name;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PublishedAt time.Time  `json:"published_at" gorm:"not null;index"`
	Feedbacks   []Feedback `json:"feedbacks,omitempty" gorm:"foreignKey:NewsID;references:ID"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	//DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"` // no soft delete
}
