package database

import (
	"database/sql"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/common"
	"stage2024/pkg/protogen/stations"
	"stage2024/pkg/protogen/users"

	"gorm.io/gorm"
)

type Bike struct {
	gorm.Model
	Id            string `gorm:"primaryKey"`
	OpenDataId    string
	BikeModel     string
	Lat           float64
	Lon           float64
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
	Id          string `gorm:"primaryKey"`
	OpenDataId  string
	Lat         float64
	Lon         float64
	Name        string
	MaxCapacity int32
	Occupation  int32
	IsActive    sql.NullBool
}

// Converts a Bike into a BikeIdentification proto event struct
func (b Bike) IntoId() *bikes.BikeIdentification {
	return &bikes.BikeIdentification{
		Id:    b.Id,
		Model: b.BikeModel,
	}
}

// Converts a User into a UserIdentification proto event struct
func (u User) IntoId() *users.UserIdentification {
	return &users.UserIdentification{
		Id:           u.Id,
		UserName:     u.UserName,
		EmailAddress: u.EmailAddress,
	}
}

// Converts a Station into a StationIdentification proto event struct
func (s Station) IntoId() *stations.StationIdentification {
	return &stations.StationIdentification{
		Id: s.Id,
		Location: &common.Location{
			Latitude:  s.Lat,
			Longitude: s.Lon,
		},
		Name: s.Name,
	}
}
