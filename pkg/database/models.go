package database

import (
	"gorm.io/gorm"
)

type Test struct {
	gorm.Model
	Code  string
	Price uint
}

type User struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	UserName    string
	EmailAdress string
}
