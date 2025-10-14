package models

import "time"

type RSSSource struct {
	Name        string      `json:"name" gorm:"primaryKey"`
	DisplayName string      `json:"display_name"`
	Description string      `json:"description"`
	Link        string      `json:"link"`
	Active      bool        `json:"active" gorm:"not null;default:true"`
	Parents     []RSSSource `json:"rss_parents,omitempty" gorm:"many2many:rss_relations;joinForeignKey:ChildName;JoinReferences:ParentName"`
	Children    []RSSSource `json:"rss_children,omitempty" gorm:"many2many:rss_relations;joinForeignKey:ParentName;JoinReferences:ChildName"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (rss *RSSSource) TableName() string {
	return "rss"
}

func (rss *RSSSource) IsMeta() bool {
	return rss.Link == ""
}
