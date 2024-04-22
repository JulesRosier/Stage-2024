package database

import (
	"database/sql"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const maxUsers = 500
const maxBikes = 500

func CreateUsersBikes(db *gorm.DB) {
	slog.Debug("Creating users and bikes")
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := CreateUsers(db); err != nil {
			return err
		}
		if err := CreateBikes(db); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		slog.Warn("Transaction failed: Error creating users and bikes", "error", err)
	}
}

// Creates 'maxUsers' amount of users
func CreateUsers(db *gorm.DB) error {
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	if userCount < maxUsers {
		slog.Debug("Creating users...", "amount", maxUsers-userCount)
	}
	for userCount < maxUsers {
		if _, err := CreateUser(db, time.Now()); err != nil {
			return err
		}

		userCount++
	}
	return nil

}

// Creates 'maxBikes' amount of bikes
func CreateBikes(db *gorm.DB) error {
	var bikeCount int64
	db.Model(&Bike{}).Count(&bikeCount)
	if bikeCount < maxBikes {
		slog.Debug("Creating bikes...", "amount", maxBikes-bikeCount)
	}
	for bikeCount < maxBikes {
		if _, err := CreateBike(db, time.Now()); err != nil {
			return err
		}
		bikeCount++
	}
	return nil
}

// creates a single user
func CreateUser(db *gorm.DB, createTime time.Time) (*User, error) {
	user := &User{
		Id:           gofakeit.UUID(),
		UserName:     gofakeit.Username(),
		EmailAddress: gofakeit.Email(),
	}
	if err := db.Create(user).Error; err != nil {
		return &User{}, err
	}

	if err := createUserEvent(user, createTime, db); err != nil {
		return &User{}, err
	}
	return user, nil
}

// creates a single bike
func CreateBike(db *gorm.DB, createTime time.Time) (*Bike, error) {
	bike := &Bike{
		Id:             uuid.New().String(),
		BikeModel:      bikeBrands[rand.IntN(len(bikeBrands))],
		Lat:            rand.Float64()*0.1 + 51.0,
		Lon:            rand.Float64()*0.2 + 3.6,
		IsElectric:     sql.NullBool{Bool: gofakeit.Bool(), Valid: true},
		IsImmobilized:  sql.NullBool{Bool: false, Valid: true},
		IsAbandoned:    sql.NullBool{Bool: false, Valid: true},
		IsInStorage:    sql.NullBool{Bool: false, Valid: true},
		IsDefect:       sql.NullBool{Bool: false, Valid: true},
		IsReturned:     sql.NullBool{Bool: true, Valid: true},
		InUseTimestamp: sql.NullTime{},
	}

	if err := db.Create(bike).Error; err != nil {
		return &Bike{}, err
	}

	slog.Debug("Bike created", "createTime", createTime, "bike", bike.Id)
	if err := BikeCreatedEvent(bike, db, createTime); err != nil {
		return &Bike{}, err
	}

	return bike, nil
}

var bikeBrands = []string{
	"Spoke-y Dokey",
	"Ride-a-licious",
	"Wheely Good Bikes",
	"Bike-a-boo",
	"ZoomZoom Bikes",
	"Handlebar Hilarity",
	"Spoke-tacular Rides",
	"Silly Spokes",
	"Whimsical Wheels",
	"Chuckling Chains",
}

func MakeFakeStation(db *gorm.DB) {
	// Create bike factory station
	station := &Station{
		Id:          "ab448be3-5c90-43ce-8c37-74f929ec016f",
		OpenDataId:  "Bike-factory123",
		Name:        "Bike factory",
		Lat:         51.037580,
		Lon:         3.735660,
		MaxCapacity: 999999999,
		Occupation:  999999999,
		IsActive:    sql.NullBool{Bool: true, Valid: true},
	}
	UpdateStation([]*Station{station}, db)

	slog.Info("updated", "station", station.Id, "occupation", station.Occupation)
}
