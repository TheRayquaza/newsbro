package models

type Article struct {
	Link string `gorm:"type:text;uniqueIndex"`
}
