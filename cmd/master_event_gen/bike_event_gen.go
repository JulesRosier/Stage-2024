package main

import (
	"database/sql"
	"log/slog"
	"math/rand/v2"
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
	"time"

	"gorm.io/gorm"
)

const minDuration = time.Minute * 5
const windowSize = time.Minute * 30

const chanceDefect = 0.03
const chanceImmobilized = 0.5
const chanceInStorage = 0.2
const format = "2006-01-02 15:04:05.999999-07"

// BikeEventGen generates bike events based on station occupation changes in the database for the last 'frequency' +20 minutes
func BikeEventGen(db *gorm.DB) {
	slog.Info("Generating bike events")
	nowUtc := time.Now().UTC().Format(format)

	decreases := []database.HistoricalStationData{}
	//get all unchecked decreases older than minDuration+windowSize
	db.Where("extract(epoch from ? - event_time_stamp)/60 >= ? and topic_name = 'station_occupation_decreased' and amount_changed > amount_faked", nowUtc, (minDuration + windowSize).Minutes()).Order("id asc").Find(&decreases)
	for _, decrease := range decreases {

		// generate amount of sequences for amount decreased/increased
		//TODO: check for increases that have no decrease
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
					err := generateBikeNotReturned(db, decrease)
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
		slog.Info("Station not found", "station", decrease.Uuid, "error", err)
		return 0, err
	}

	delta := decrease.EventTimeStamp.Sub(stationCreated.EventTimeStamp)

	if delta > windowSize {
		delta = windowSize
	}

	r := rand.Float64()
	if r == 0 {
		r = 0.5
	}

	offset := time.Second * time.Duration(r*delta.Seconds())

	slog.Debug("Start offset", "offset", offset)
	return offset, err
}

// Generates events based on increase and decrease in station occupation
func generate(db *gorm.DB, increase database.HistoricalStationData, decrease database.HistoricalStationData) error {
	slog.Info("Generating event sequence", "station", decrease.OpenDataId)
	slog.Debug("Stations", "increase", increase.OpenDataId, "decrease", decrease.OpenDataId)

	// number of minutes bike is reserved before it is picked up
	startoffset, err := getStartOffset(db, decrease)
	if err != nil {
		return err
	}

	startTime := decrease.EventTimeStamp.Add(-startoffset).UTC()
	endTime := increase.EventTimeStamp.UTC()

	bike, user, err := getBikeAndUser(db, startTime)
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
	database.BikeReservedEvent(bike, change, db)

	// same time as capacity decrease
	// bike picked up
	change.EventTime = decrease.EventTimeStamp.UTC()
	bike.IsReturned = sql.NullBool{Bool: false, Valid: true}
	database.UpdateBike(db, &bike)
	database.BikePickedUpEvent(bike, change, db)

	// after capacity decrease
	delta := decrease.EventTimeStamp.Sub(endTime)
	r := rand.Float64()
	if r == 0 {
		r = 0.5
	}
	minutes := r * delta.Minutes()

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
	bike.IsReturned = sql.NullBool{Bool: true, Valid: true}
	change.StationId = increase.Uuid
	change.EventTime = endTime

	database.UpdateBike(db, &bike)
	database.BikeReturnedEvent(bike, change, db)

	db.Model(&decrease).Update("amount_faked", decrease.AmountFaked)
	db.Model(&increase).Update("amount_faked", increase.AmountFaked)
	return nil
}

func generateBikeNotReturned(db *gorm.DB, decrease database.HistoricalStationData) error {
	slog.Info("Generating event sequence, Bike is not returned", "station", decrease.OpenDataId)

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

	bike, user, err := getBikeAndUser(db, startTime)
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
	database.UpdateBike(db, &bike)
	database.BikeReservedEvent(bike, change, db)

	// same time as capacity decrease
	// bike picked up
	change.EventTime = pickedUpTime
	bike.IsReturned = sql.NullBool{Bool: false, Valid: true}
	database.UpdateBike(db, &bike)
	database.BikePickedUpEvent(bike, change, db)

	// chance bikedefect
	if rand.Float32() < chanceDefect+0.5 {
		change.EventTime = defectTime
		change.Defect = defects[rand.IntN(len(defects))]
		bike.IsDefect = sql.NullBool{Bool: true, Valid: true}
		database.UpdateBike(db, &bike)
		database.BikeDefectEvent(bike, change, db)

		// chance bike immobilized
		if rand.Float32() < chanceImmobilized {
			change.EventTime = immobilizedTime
			bike.IsImmobilized = sql.NullBool{Bool: true, Valid: true}
			database.UpdateBike(db, &bike)
			database.BikeImmobilizedEvent(bike, change, db)
		}
	}
	// bike abandoned
	change.EventTime = abandonedTime
	bike.IsAbandoned = sql.NullBool{Bool: true, Valid: true}
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

	db.Model(&decrease).Update("amount_faked", decrease.AmountFaked)

	return nil
}

func getBikeAndUser(db *gorm.DB, startTime time.Time) (database.Bike, database.User, error) {
	startTimeString := startTime.Format(format)
	//get available bike not in use
	bike := database.Bike{}
	if result := db.Where("(is_defect = false AND is_immobilized = false AND is_abandoned = false AND is_in_storage = false AND is_returned = true) AND (in_use_timestamp is null OR (extract(epoch from ? - in_use_timestamp)/60 > 0) AND extract(epoch from ? - created_at)/60 > 0)", startTimeString, startTimeString).Order("random()").Limit(1).Find(&bike); result.RowsAffected == 0 {
		// If no bike available, clean up bikes
		if err := database.BikeCleanUp(db); err != nil {
			return database.Bike{}, database.User{}, err
		}
		if result := db.Where("(AND is_defect = false AND is_immobilized = false AND is_abandoned = false AND is_in_storage = false AND is_returned = true) AND (in_use_timestamp is null OR (extract(epoch from ? - in_use_timestamp)/60 > 0) AND extract(epoch from ? - created_at)/60 > 0)", startTimeString, startTimeString).Order("random()").Limit(1).Find(&bike); result.RowsAffected == 0 {
			// If still no bike available, create bike
			slog.Info("No available bike found, creating bike")
			newBike, err := database.CreateBike(db, startTime.Add(-time.Minute*30))
			if err != nil {
				return database.Bike{}, database.User{}, err
			}
			bike = *newBike
		}

	}
	// get available user
	user := database.User{}
	result := db.Where("is_available_timestamp is null OR (extract(epoch from ? - is_available_timestamp)/60 > 0 AND extract(epoch from ? - created_at)/60 > 0)", startTimeString, startTimeString).Order("random()").Limit(1).Find(&user)
	if result.RowsAffected == 0 {
		slog.Warn("No available user found, creating user")
		newUser, err := database.CreateUser(db)
		if err != nil {
			return database.Bike{}, database.User{}, err
		}
		user = *newUser
	}

	slog.Info("Bike & User selected", "bike", bike.Id, "user", user.Id)
	return bike, user, nil
}

var defects = []string{
	"Handlebars are actually spaghetti strands",
	"Wheels keep trying to escape to join the circus",
	"Saddle is a whoopee cushion",
	"Pedals rotate backward, propelling you into the past",
	"Bell plays 'Jingle Bells' at random intervals",
	"Chain is made of rubber bands",
	"Brakes function as accelerators",
	"Bike frame is held together by duct tape",
	"Lights only work when the moon is full",
	"Tires are square",
	"Bike speaks only in Morse code",
	"Basket is actually a miniature black hole",
	"Reflectors emit disco lights",
	"Kickstand has commitment issues",
	"Water bottle holder dispenses hot sauce instead",
	"GPS always directs you to the nearest ice cream parlor",
	"Bike lock is just a piece of string",
	"Bells and whistles are literal bells and whistles",
	"Horn plays 'La Cucaracha' off-key",
	"Handlebars rotate 360 degrees uncontrollably",
	"Seat cushion is made of cactus needles",
	"Pedals detach mid-ride for impromptu dance parties",
	"Frame is magnetically attracted to garbage cans",
	"Saddle has a 'kick me' sign taped to it",
	"Bike basket has a pet rock as a passenger",
	"Bike chain sings 'The Wheels on the Bus' endlessly",
	"Gears shift randomly to reverse",
	"Reflectors reflect sarcastic remarks",
	"Handlebar grips are actually bananas",
}
