package database

import (
	"database/sql"
	"log/slog"
	"math/rand/v2"
	"stage2024/pkg/helper"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const maxUsers = 200
const maxBikes = 333

// Creates 'maxUsers' amount of users
func CreateUsers(db *gorm.DB) {
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	users := []*User{}
	for userCount < maxUsers {
		users = append(users, &User{
			Id:           gofakeit.UUID(),
			UserName:     gofakeit.Username(),
			EmailAddress: gofakeit.Email(),
		})

		userCount++
	}
	UpdateUser(users, db)
}

// Creates 'maxBikes' amount of bikes
func CreateBikes(db *gorm.DB) {
	slog.Debug("Creating bikes...")
	var bikeCount int64
	db.Model(&Bike{}).Count(&bikeCount)
	bikes := []*Bike{}

	for bikeCount < maxBikes {
		bikes = append(bikes, &Bike{
			Id:             uuid.New().String(),
			BikeModel:      bikeBrands[rand.IntN(len(bikeBrands))],
			Lat:            rand.Float64()*1.0 + 51.0,
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
		})

		bikeCount++
	}
	UpdateBike(bikes, db, helper.Change{})
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
