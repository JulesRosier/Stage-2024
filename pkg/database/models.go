package database

import (
	"gorm.io/gorm"
)

type Bike struct {
	gorm.Model
	Id            string `gorm:"primaryKey"`
	BikeModel     string
	IsElectric    string
	Location      string
	IsImmobilized bool
	IsAbandoned   bool
	IsAvalable    bool
	IsInStorage   bool
	IsReserved    bool
	IsDefect      bool
}

type User struct {
	gorm.Model
	Id           string `gorm:"primaryKey"`
	UserName     string
	EmailAddress string
}

type Station struct {
	gorm.Model
	Id         string `gorm:"primaryKey"`
	Location   string
	Name       string
	Occupation int32
	IsActive   bool
}
