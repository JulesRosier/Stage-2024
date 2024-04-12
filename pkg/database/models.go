package database

import (
	"database/sql"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/common"
	"stage2024/pkg/protogen/stations"
	"stage2024/pkg/protogen/users"
	"time"

	"gorm.io/gorm"
)

type Bike struct {
	gorm.Model
	Id             string `gorm:"primaryKey"`
	BikeModel      string
	Lat            float64
	Lon            float64
	IsElectric     sql.NullBool
	IsReserved     sql.NullBool
	PickedUp       sql.NullBool
	IsDefect       sql.NullBool
	IsImmobilized  sql.NullBool
	IsAbandoned    sql.NullBool
	IsInStorage    sql.NullBool
	IsReturned     sql.NullBool
	InUseTimestamp sql.NullTime
}

type User struct {
	gorm.Model
	Id                   string `gorm:"primaryKey"`
	UserName             string
	EmailAddress         string
	IsAvailableTimestamp sql.NullTime
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

type Outbox struct {
	gorm.Model
	EventTimestamp time.Time
	Topic          string
	Payload        []byte
}

type HistoricalStationData struct {
	gorm.Model
	EventTimeStamp time.Time
	Uuid           string
	OpenDataId     string
	Lat            float64
	Lon            float64
	Name           string
	MaxCapacity    int32
	Occupation     int32
	IsActive       sql.NullBool
	TopicName      string
	AmountChanged  int32
	AmountFaked    int32
}

type BikeGenData struct {
	EventTime time.Time
	StationId string
	UserId    string
}

func (s *Station) ToHistoricalStationData() HistoricalStationData {
	return HistoricalStationData{
		Uuid:        s.Id,
		OpenDataId:  s.OpenDataId,
		Lat:         s.Lat,
		Lon:         s.Lon,
		Name:        s.Name,
		MaxCapacity: s.MaxCapacity,
		Occupation:  s.Occupation,
		IsActive:    s.IsActive,
	}
}

func (Outbox) TableName() string {
	return "outbox"
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
