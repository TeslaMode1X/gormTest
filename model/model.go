package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null" json:"username"`
	Age      int    `gorm:"not null" json:"age"`
	Job      string `gorm:"not null" json:"job"`
}
