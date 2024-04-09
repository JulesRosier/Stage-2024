package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const maxBikes = 333
const minDuration = time.Minute * 10
const windowSize = time.Minute * 30

// BikeEventGen generates bike events
func BikeEventGen(db *gorm.DB, frequency int) {
	createBikes(db)
	nowUtc := time.Now().UTC().Format("2006-01-02 15:04:05.999999-07")

	topicname := "station_occupation_decreased"

	decreases := []database.HistoricalStationData{}
	//get all unchecked decreases in the last ... minutes
	db.Where("extract(epoch from ? - created_at)/60 <= ? and topic_name = ? and amount_changed > amount_faked", nowUtc, frequency+20, topicname).Order("updated_at asc").Find(&decreases)

	for _, decrease := range decreases {
		increase := database.HistoricalStationData{}
		topicname := "station_occupation_increased"
		//get all increases between decrease+minDuration and decrease+minDuration+windowSize
		result := db.Where("created_at between ? and ? and topic_name = ? and amount_changed > amount_faked", decrease.CreatedAt.Add(minDuration), decrease.CreatedAt.Add(minDuration+windowSize), topicname).Order("random()").Limit(1).Find(&increase)
		if result.RowsAffected == 0 {
			slog.Info("No increase found for decrease", "increase", decrease.OpenDataId)
			break
		}
		// TODO: generate amount of sequences for amount decreased
		//TODO: check if increase station has multiple increases, alhoewel dat boeit hier eigenlijk niet, als er dan wordt
		//TODO: vragen aan Jules om timestamp op redpanda zelf te veranderen naar timestamp van event
		//gekeken naar de volgende decrease, dan is dat gewoon een andere increase waarbij die station opniew kan worden gebruikt
		// en de aantal fakes erin wordt bijgehouden dus ideaal
		// start sequence with decrease
		generate(db, increase, decrease)
	}
}

// Creates const maxBikes amount of bikes
func createBikes(db *gorm.DB) {
	slog.Debug("Creating bikes...")
	var bikeCount int64
	db.Model(&database.Bike{}).Count(&bikeCount)

	for bikeCount < maxBikes {
		bike := &database.Bike{
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
			InUseTimestamp: sql.NullTime{},
		}
		db.Create(bike)
		bikeCount++
	}
}

// Generates events based on increase and decrease in station occupation
func generate(db *gorm.DB, increase database.HistoricalStationData, decrease database.HistoricalStationData) {
	slog.Info("Generating event sequence", "increase", increase.OpenDataId)

	startoffset := helper.RandMinutes(60*5, 60)
	startTime := decrease.CreatedAt.Add(-startoffset).UTC()
	startTimeString := startTime.Format("2006-01-02 15:04:05.999999-07")
	endTime := increase.CreatedAt.UTC()

	//get available bike not in use
	bike := database.Bike{}
	result := db.Where("is_reserved = false AND in_use_timestamp is null OR extract(epoch from ? - in_use_timestamp)/60 > 0", startTimeString).Order("random()").Limit(1).Find(&bike)
	if result.RowsAffected == 0 {
		slog.Warn("No available bike found")
		helper.Die(fmt.Errorf("no available bike found"))
	}
	// get available user
	user := database.User{}
	result = db.Where("is_available_timestamp is null OR extract(epoch from ? - is_available_timestamp)/60 > 0", startTimeString).Order("random()").Limit(1).Find(&user)
	if result.RowsAffected == 0 {
		slog.Warn("No available user found")
		helper.Die(fmt.Errorf("no available user found"))
	}

	slog.Info("Bike selected", "bike", bike.Id)
	slog.Info("User selected", "user", user.Id)

	// before station capacity decrease
	//bike reserved
	// Need to find a way to update the bike, but change the timestamp in the event...
	// option: change updateBike to take a timestamp to send to event..
	bike.InUseTimestamp = sql.NullTime{Time: endTime, Valid: true}
	bike.IsReserved = sql.NullBool{Bool: true, Valid: true}

	change := helper.Change{
		EventTime: startTime,
		StationId: increase.Uuid,
		UserId:    user.Id,
	}
	database.UpdateBike([]*database.Bike{&bike}, db, change)

	// same time as capacity decrease
	// bike picked up
	change.EventTime = decrease.CreatedAt.UTC()
	bike.PickedUp = sql.NullBool{Bool: true, Valid: true}
	database.UpdateBike([]*database.Bike{&bike}, db, change)

	// after capacity decrease
	// chance bike abandoned
	// chance bike defect
	// chance bike immobilized
	// chance bike in storage

	// same time as capacity increase
	// bike dropped off

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

//TODO dingen die niet werken in ui, zelfde event als twee verschillende, met andere structuren, ligt het aan mij?
// station id : b510863b-e6e1-4edd-a4cb-28995e1ed455
// bike_id : ec502b37-5fa3-4c10-ab86-97ba2b763b52
