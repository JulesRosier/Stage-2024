package database

import (
	"database/sql"
	"math/rand/v2"
	"stage2024/pkg/helper"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const maxUsers = 200
const maxBikes = 444

func CreateUsersBikes(db *gorm.DB) {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := CreateUsers(db); err != nil {
			return err
		}
		if err := CreateBikes(db); err != nil {
			return err
		}
		return nil
	})
	helper.MaybeDie(err, "Transaction failed: Error creating users and bikes")
}

// Creates 'maxUsers' amount of users
func CreateUsers(db *gorm.DB) error {
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	err := db.Transaction(func(tx *gorm.DB) error {
		for userCount < maxUsers {
			user := &User{
				Id:           gofakeit.UUID(),
				UserName:     gofakeit.Username(),
				EmailAddress: gofakeit.Email(),
			}
			if err := db.Create(user).Error; err != nil {
				return err
			}

			if err := createUserEvent(user, db); err != nil {
				return err
			}

			userCount++
		}
		return nil
	})
	return err
}

// Creates 'maxBikes' amount of bikes
func CreateBikes(db *gorm.DB) error {
	var bikeCount int64
	db.Model(&Bike{}).Count(&bikeCount)
	err := db.Transaction(func(tx *gorm.DB) error {
		for bikeCount < maxBikes {
			bike := &Bike{
				Id:             uuid.New().String(),
				BikeModel:      bikeBrands[rand.IntN(len(bikeBrands))],
				Lat:            rand.Float64()*0.1 + 51.0,
				Lon:            rand.Float64()*0.2 + 3.6,
				IsElectric:     sql.NullBool{Bool: gofakeit.Bool(), Valid: true},
				PickedUp:       sql.NullBool{Bool: false, Valid: true},
				IsImmobilized:  sql.NullBool{Bool: false, Valid: true},
				IsAbandoned:    sql.NullBool{Bool: false, Valid: true},
				IsInStorage:    sql.NullBool{Bool: false, Valid: true},
				IsReserved:     sql.NullBool{Bool: false, Valid: true},
				IsDefect:       sql.NullBool{Bool: false, Valid: true},
				IsReturned:     sql.NullBool{Bool: true, Valid: true},
				InUseTimestamp: sql.NullTime{},
			}

			if err := db.Create(bike).Error; err != nil {
				return err
			}

			if err := BikeCreatedEvent(bike, db); err != nil {
				return err
			}

			bikeCount++
		}
		return nil
	})
	return err
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
