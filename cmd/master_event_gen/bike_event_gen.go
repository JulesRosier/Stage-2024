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

const chanceAbandoned = 0.2
const chanceDefect = 0.1
const chanceImmobilized = 0.5
const chanceInStorage = 0.2

// BikeEventGen generates bike events
func BikeEventGen(db *gorm.DB, frequency int) {
	createBikes(db)
	nowUtc := time.Now().UTC().Format("2006-01-02 15:04:05.999999-07")

	decreases := []database.HistoricalStationData{}
	//get all unchecked decreases in the last ... minutes
	db.Where("extract(epoch from ? - event_time_stamp)/60 <= ? and topic_name = 'station_occupation_decreased' and amount_changed > amount_faked", nowUtc, frequency+20).Order("updated_at asc").Find(&decreases)

	for _, decrease := range decreases {

		// generate amount of sequences for amount decreased/increased
		// start sequence with decrease
		slog.Debug("Decrease found", "decrease", decrease.OpenDataId)
		slog.Debug("Amount changed", "amount", decrease.AmountChanged)
		for range decrease.AmountChanged {
			decrease.AmountFaked++

			//get increase for decrease
			increase := database.HistoricalStationData{}
			topicname := "station_occupation_increased"
			//get all increases between decrease+minDuration and decrease+minDuration+windowSize
			result := db.Where("event_time_stamp between ? and ? AND topic_name = ? AND amount_changed > amount_faked", decrease.EventTimeStamp.Add(minDuration), decrease.EventTimeStamp.Add(minDuration+windowSize), topicname).Order("random()").Limit(1).Find(&increase)
			if result.RowsAffected == 0 {
				slog.Debug("No increase found for decrease", "increase", decrease.OpenDataId)
				// get bike abandoned/ immobilized events here
				generateNotReturned(db, decrease)
				break
			}
			increase.AmountFaked++
			generate(db, increase, decrease)
		}
	}
}

// TODO: start transaction for amount changed/ amount faked????,
// TODO: some stations amount changes are negative numbers....
// Creates 'maxBikes' amount of bikes
func createBikes(db *gorm.DB) {
	slog.Debug("Creating bikes...")
	var bikeCount int64
	db.Model(&database.Bike{}).Count(&bikeCount)
	bikes := []*database.Bike{}

	for bikeCount < maxBikes {
		bikes = append(bikes, &database.Bike{
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
	database.UpdateBike(bikes, db, helper.Change{})
}

// Generates events based on increase and decrease in station occupation
func generate(db *gorm.DB, increase database.HistoricalStationData, decrease database.HistoricalStationData) {
	slog.Info("Generating event sequence", "station", decrease.OpenDataId)

	// number of minutes bike is reserved before it is picked up
	startoffset := helper.RandMinutes(60*5, 60)

	startTime := decrease.EventTimeStamp.Add(-startoffset).UTC()
	startTimeString := startTime.Format("2006-01-02 15:04:05.999999-07")
	endTime := increase.EventTimeStamp.UTC()
	delta := startTime.Sub(endTime)

	bike, user := getBikeAndUser(db, startTimeString)

	// before station capacity decrease
	// bike reserved
	user.IsAvailableTimestamp = sql.NullTime{Time: endTime, Valid: true}
	bike.InUseTimestamp = sql.NullTime{Time: endTime, Valid: true}
	bike.IsReserved = sql.NullBool{Bool: true, Valid: true}

	change := helper.Change{
		EventTime: startTime,
		StationId: decrease.Uuid,
		UserId:    user.Id,
	}
	database.UpdateBike([]*database.Bike{&bike}, db, change)
	database.UpdateUser([]*database.User{&user}, db)

	// same time as capacity decrease
	// bike picked up
	change.EventTime = decrease.EventTimeStamp.UTC()
	bike.PickedUp = sql.NullBool{Bool: true, Valid: true}
	bike.IsReturned = sql.NullBool{Bool: false, Valid: true}
	database.UpdateBike([]*database.Bike{&bike}, db, change)

	// after capacity decrease
	minutes := rand.Float64() * delta.Minutes()

	// chance bike defect
	if rand.Float32() < chanceDefect {
		change.EventTime = startTime.Add(time.Minute * time.Duration(minutes))
		change.Defect = defects[rand.IntN(len(defects))]
		bike.IsDefect = sql.NullBool{Bool: true, Valid: true}
		database.UpdateBike([]*database.Bike{&bike}, db, change)
	}

	// same time as capacity increase
	// bike returned

	bike.PickedUp = sql.NullBool{Bool: false, Valid: true}
	bike.IsReserved = sql.NullBool{Bool: false, Valid: true}
	bike.IsReturned = sql.NullBool{Bool: true, Valid: true}
	change.StationId = increase.Uuid
	change.EventTime = endTime

	database.UpdateBike([]*database.Bike{&bike}, db, change)

	db.Model(&decrease).Update("amount_faked", decrease.AmountFaked)
	db.Model(&increase).Update("amount_faked", increase.AmountFaked)
}

func generateNotReturned(db *gorm.DB, decrease database.HistoricalStationData) {
	slog.Info("Generating NOT returned event sequence", "decrease", decrease.OpenDataId)

	// number of minutes bike is reserved before it is picked up
	startoffset := helper.RandMinutes(60*5, 60)

	startTime := decrease.EventTimeStamp.Add(-startoffset).UTC()
	startTimeString := startTime.Format("2006-01-02 15:04:05.999999-07")
	pickedUpTime := decrease.EventTimeStamp.UTC()
	defectTime := pickedUpTime.Add(helper.RandMinutes(60*5, 60)).UTC()
	immobilizedTime := defectTime.Add(helper.RandMinutes(2*5, 2)).UTC()
	inStorageTime := immobilizedTime.Add(helper.RandMinutes(2*5, 2)).UTC()
	endTime := inStorageTime.Add(helper.RandMinutes(60*5, 60)).UTC()

	bike, user := getBikeAndUser(db, startTimeString)

	// before station capacity decrease
	// bike reserved
	user.IsAvailableTimestamp = sql.NullTime{Time: endTime, Valid: true}
	bike.InUseTimestamp = sql.NullTime{Time: endTime, Valid: true}
	bike.IsReserved = sql.NullBool{Bool: true, Valid: true}

	change := helper.Change{
		EventTime: startTime,
		StationId: decrease.Uuid,
		UserId:    user.Id,
	}
	database.UpdateBike([]*database.Bike{&bike}, db, change)
	database.UpdateUser([]*database.User{&user}, db)

	// same time as capacity decrease
	// bike picked up
	change.EventTime = pickedUpTime
	bike.PickedUp = sql.NullBool{Bool: true, Valid: true}
	bike.IsReturned = sql.NullBool{Bool: false, Valid: true}
	database.UpdateBike([]*database.Bike{&bike}, db, change)

	db.Model(&decrease).Update("amount_faked", decrease.AmountFaked)

	// chance bike abandoned
	if rand.Float32() < chanceAbandoned {
		change.EventTime = endTime
		bike.IsAbandoned = sql.NullBool{Bool: true, Valid: true}
		bike.IsReserved = sql.NullBool{Bool: false, Valid: true}
		database.UpdateBike([]*database.Bike{&bike}, db, change)
		return
	} else { // bike defect
		change.EventTime = defectTime
		change.Defect = defects[rand.IntN(len(defects))]
		bike.IsDefect = sql.NullBool{Bool: true, Valid: true}
		bike.IsReserved = sql.NullBool{Bool: false, Valid: true}
		database.UpdateBike([]*database.Bike{&bike}, db, change)

		// chance bike immobilized
		if rand.Float32() < chanceImmobilized {
			change.EventTime = immobilizedTime
			bike.IsImmobilized = sql.NullBool{Bool: true, Valid: true}
			database.UpdateBike([]*database.Bike{&bike}, db, change)
		}
		// chance bike in storage
		if rand.Float32() < chanceInStorage {
			change.EventTime = inStorageTime
			bike.IsInStorage = sql.NullBool{Bool: true, Valid: true}
			database.UpdateBike([]*database.Bike{&bike}, db, change)
		}
	}
}

func getBikeAndUser(db *gorm.DB, startTimeString string) (database.Bike, database.User) {
	//get available bike not in use
	bike := database.Bike{}
	result := db.Where("is_reserved = false AND is_defect = false AND is_immobilized = false AND is_abandoned = false AND is_in_storage = false AND is_returned = true AND in_use_timestamp is null OR extract(epoch from ? - in_use_timestamp)/60 > 0", startTimeString).Order("random()").Limit(1).Find(&bike)
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
	return bike, user
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

var defects = []string{
	"Flat tire",
	"Broken chain",
	"Worn brake pads",
	"Loose spokes",
	"Faulty gear shifting",
	"Bent wheel rim",
	"Damaged pedals",
	"Cracked frame",
	"Stuck brakes",
	"Rusted components",
	"Misaligned wheels",
	"Broken saddle",
	"Malfunctioning gears",
	"Wobbly handlebars",
	"Loose headset",
	"Torn seat cover",
	"Defective bearings",
	"Faulty brakes",
	"Cracked fork",
	"Damaged crankset",
}
