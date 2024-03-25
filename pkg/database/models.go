package database

import (
	"database/sql"

	"gorm.io/gorm"
)

type Bike struct {
	gorm.Model
	Id            string `gorm:"primaryKey"`
	BikeModel     string
	Location      string
	IsElectric    sql.NullBool
	IsImmobilized sql.NullBool
	IsAbandoned   sql.NullBool
	IsAvailable   sql.NullBool
	IsInStorage   sql.NullBool
	IsReserved    sql.NullBool
	IsDefect      sql.NullBool
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
	IsActive   sql.NullBool
}
