package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
	"time"

	"gorm.io/gorm"
)

const minDuration = time.Minute * 5
const windowSize = time.Minute * 30

const chanceDefect = 0.1
const chanceImmobilized = 0.5
const chanceInStorage = 0.2
const format = "2006-01-02 15:04:05.999999-07"

// BikeEventGen generates bike events based on station occupation changes in the database for the last 'frequency' +20 minutes
func BikeEventGen(db *gorm.DB, frequency int) {
	slog.Info("Generating bike events")
	nowUtc := time.Now().UTC().Format(format)

	decreases := []database.HistoricalStationData{}
	//get all unchecked decreases older than minDuration+windowSize
	db.Where("extract(epoch from ? - event_time_stamp)/60 >= ? and topic_name = 'station_occupation_decreased' and amount_changed > amount_faked", nowUtc, (minDuration + windowSize).Minutes()).Order("updated_at asc").Find(&decreases)
	for _, decrease := range decreases {

		// generate amount of sequences for amount decreased/increased
		// start sequence with decrease
		for range decrease.AmountChanged {
			decrease.AmountFaked++

			// start transaction here
			err := db.Transaction(func(tx *gorm.DB) error {
				//get increase for decrease within window
				increase := database.HistoricalStationData{}
				result := db.Where("event_time_stamp between ? and ? AND topic_name = 'station_occupation_increased' AND amount_changed > amount_faked", decrease.EventTimeStamp.Add(minDuration), decrease.EventTimeStamp.Add(minDuration+windowSize)).Order("random()").Limit(1).Find(&increase)

				if result.RowsAffected == 0 {
					slog.Debug("No increase found for decrease", "increase", decrease.OpenDataId)
					// get bike abandoned/ immobilized events here
					err := generateNotReturned(db, decrease)
					return err
				} else {
					increase.AmountFaked++
					err := generate(db, increase, decrease)
					return err
				}
			})
			if err != nil {
				slog.Warn("Fake event transaction failed", "error", err)
				decrease.AmountFaked--
			}
		}
	}
}

// Gets random start offset for event sequence that is not before station creation
func getStartOffset(db *gorm.DB, decrease database.HistoricalStationData) (time.Duration, error) {
	stationCreated := database.HistoricalStationData{}
	err := db.Where("uuid = ? AND topic_name = 'station_created'", decrease.Uuid).Limit(1).Find(&stationCreated).Error
	if err != nil {
		slog.Warn("Station not found", "station", decrease.Uuid)
		helper.Die(err)
		return 0, err
	}

	delta := decrease.EventTimeStamp.Sub(stationCreated.EventTimeStamp)

	if delta > windowSize {
		delta = windowSize
	}
	offset := time.Minute * time.Duration(rand.Float64()*delta.Minutes())

	slog.Info("Start offset", "offset", offset)
	return offset, err
}

// Generates events based on increase and decrease in station occupation
func generate(db *gorm.DB, increase database.HistoricalStationData, decrease database.HistoricalStationData) error {
	slog.Info("Generating event sequence", "station", decrease.OpenDataId)

	// number of minutes bike is reserved before it is picked up
	startoffset, err := getStartOffset(db, decrease)
	if err != nil {
		return err
	}

	startTime := decrease.EventTimeStamp.Add(-startoffset).UTC()
	startTimeString := startTime.Format(format)
	endTime := increase.EventTimeStamp.UTC()
	delta := startTime.Sub(endTime)

	bike, user, err := getBikeAndUser(db, startTimeString)
	if err != nil {
		return err
	}

	user.IsAvailableTimestamp = sql.NullTime{Time: endTime, Valid: true}
	bike.InUseTimestamp = sql.NullTime{Time: endTime, Valid: true}
	database.UpdateBike(db, &bike)
	database.UpdateUser(db, &user)

	change := helper.Change{
		EventTime: startTime,
		StationId: decrease.Uuid,
		UserId:    user.Id,
	}

	// before station capacity decrease
	// bike reserved
	bike.IsReserved = sql.NullBool{Bool: true, Valid: true}
	database.UpdateBike(db, &bike)
	database.BikeReserved(bike, change, db)

	// same time as capacity decrease
	// bike picked up
	change.EventTime = decrease.EventTimeStamp.UTC()
	bike.PickedUp = sql.NullBool{Bool: true, Valid: true}
	bike.IsReturned = sql.NullBool{Bool: false, Valid: true}
	database.UpdateBike(db, &bike)
	database.BikePickedUpEvent(bike, change, db)

	// after capacity decrease
	minutes := rand.Float64() * delta.Minutes()

	// chance bike defect
	if rand.Float32() < chanceDefect {
		change.EventTime = startTime.Add(time.Minute * time.Duration(minutes))
		change.Defect = defects[rand.IntN(len(defects))]
		bike.IsDefect = sql.NullBool{Bool: true, Valid: true}
		database.UpdateBike(db, &bike)
		database.BikeDefectEvent(bike, change, db)
	}

	// same time as capacity increase
	// bike returned
	bike.PickedUp = sql.NullBool{Bool: false, Valid: true}
	bike.IsReserved = sql.NullBool{Bool: false, Valid: true}
	bike.IsReturned = sql.NullBool{Bool: true, Valid: true}
	change.StationId = increase.Uuid
	change.EventTime = endTime

	database.UpdateBike(db, &bike)
	database.BikeReturnedEvent(bike, change, db)

	db.Model(&decrease).Update("amount_faked", decrease.AmountFaked)
	db.Model(&increase).Update("amount_faked", increase.AmountFaked)
	return nil
}

func generateNotReturned(db *gorm.DB, decrease database.HistoricalStationData) error {
	slog.Info("Generating NOT returned event sequence", "decrease", decrease.OpenDataId)

	// number of minutes bike is reserved before it is picked up
	startoffset, err := getStartOffset(db, decrease)
	if err != nil {
		return err
	}

	startTime := decrease.EventTimeStamp.Add(-startoffset).UTC()
	pickedUpTime := decrease.EventTimeStamp.UTC()
	defectTime := pickedUpTime.Add(helper.RandMinutes(60*5, 60)).UTC()
	immobilizedTime := defectTime.Add(helper.RandMinutes(2*5, 2)).UTC()
	abandonedTime := immobilizedTime.Add(helper.RandMinutes(2*5, 2)).UTC()
	endTime := abandonedTime.Add(helper.RandMinutes(60*5, 60)).UTC()

	bike, user, err := getBikeAndUser(db, startTime.Format(format))
	if err != nil {
		return err
	}

	user.IsAvailableTimestamp = sql.NullTime{Time: endTime, Valid: true}
	bike.InUseTimestamp = sql.NullTime{Time: endTime, Valid: true}
	database.UpdateBike(db, &bike)
	database.UpdateUser(db, &user)

	change := helper.Change{
		EventTime: startTime,
		StationId: decrease.Uuid,
		UserId:    user.Id,
	}
	// before station capacity decrease
	// bike reserved
	bike.IsReserved = sql.NullBool{Bool: true, Valid: true}
	database.UpdateBike(db, &bike)
	database.BikeReserved(bike, change, db)

	// same time as capacity decrease
	// bike picked up
	change.EventTime = pickedUpTime
	bike.PickedUp = sql.NullBool{Bool: true, Valid: true}
	bike.IsReturned = sql.NullBool{Bool: false, Valid: true}
	database.UpdateBike(db, &bike)
	database.BikePickedUpEvent(bike, change, db)

	// chance bikedefect
	if rand.Float32() < chanceDefect+0.5 {
		change.EventTime = defectTime
		change.Defect = defects[rand.IntN(len(defects))]
		bike.IsDefect = sql.NullBool{Bool: true, Valid: true}
		bike.IsReserved = sql.NullBool{Bool: false, Valid: true}
		database.UpdateBike(db, &bike)
		database.BikeDefectEvent(bike, change, db)

		// chance bike immobilized
		if rand.Float32() < chanceImmobilized {
			change.EventTime = immobilizedTime
			bike.IsImmobilized = sql.NullBool{Bool: true, Valid: true}
			database.UpdateBike(db, &bike)
			database.BikeImmobilizedEvent(bike, change, db)
		}

		// bike abandoned
		change.EventTime = abandonedTime
		bike.IsAbandoned = sql.NullBool{Bool: true, Valid: true}
		bike.IsReserved = sql.NullBool{Bool: false, Valid: true}
		database.UpdateBike(db, &bike)
		database.BikeAbandonedEvent(bike, change, db)

		// chance bike in storage
		if rand.Float32() < chanceInStorage {
			change.EventTime = endTime
			change.NewValue = "true"
			bike.IsInStorage = sql.NullBool{Bool: true, Valid: true}
			database.UpdateBike(db, &bike)
			database.BikeInStorageEvent(bike, change, db)
		}
	}

	db.Model(&decrease).Update("amount_faked", decrease.AmountFaked)

	return nil
}

func getBikeAndUser(db *gorm.DB, startTimeString string) (database.Bike, database.User, error) {
	//get available bike not in use
	bike := database.Bike{}
	result := db.Where("is_reserved = false AND is_defect = false AND is_immobilized = false AND is_abandoned = false AND is_in_storage = false AND is_returned = true AND in_use_timestamp is null OR extract(epoch from ? - in_use_timestamp)/60 > 0 and extract(epoch from ? - created_at)/60 > 0", startTimeString, startTimeString).Order("random()").Limit(1).Find(&bike)
	if result.RowsAffected == 0 {
		slog.Warn("No available bike found")
		return database.Bike{}, database.User{}, fmt.Errorf("no available bike found")
	}
	// get available user
	user := database.User{}
	result = db.Where("is_available_timestamp is null OR extract(epoch from ? - is_available_timestamp)/60 > 0 and extract(epoch from ? - created_at)/60 > 0", startTimeString, startTimeString).Order("random()").Limit(1).Find(&user)
	if result.RowsAffected == 0 {
		slog.Warn("No available user found")
		return database.Bike{}, database.User{}, fmt.Errorf("no available user found")
	}

	slog.Info("Bike selected", "bike", bike.Id)
	slog.Info("User selected", "user", user.Id)
	return bike, user, nil
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
