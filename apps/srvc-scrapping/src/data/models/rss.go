package models

type RSS struct {
	Link string `gorm:"type:text;uniqueIndex"`
}
