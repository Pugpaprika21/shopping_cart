package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string `gorm:"type:varchar(200);not null"`
	Description string `gorm:"type:text"`
}
