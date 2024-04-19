package events

import (
	"database/sql"
	"log/slog"
	"math/rand/v2"
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
	"time"

	"gorm.io/gorm"
)

const (
	minDuration             = time.Minute * 10
	windowSize              = time.Minute * 30
	chanceDefect            = 0.03
	chanceDefectNotReturned = 0.5
	chanceImmobilized       = 0.5
	format                  = "2006-01-02 15:04:05.999999-07"
)

// BikeEventGen generates bike events based on station occupation changes in the database
func BikeEventGen(db *gorm.DB) {
	slog.Info("Generating bike events")
	nowUtc := time.Now().UTC().Format(format)

	decreases := []database.HistoricalStationData{}
	//get all unchecked decreases older than minDuration+windowSize
	err := db.Where("extract(epoch from ? - event_time_stamp)/60 >= ? and topic_name = 'station_occupation_decreased' and amount_changed > amount_faked",
		nowUtc, (minDuration + windowSize).Minutes()).Order("id asc").Find(&decreases).Error
	if err != nil {
		slog.Warn("Failed to get decreases", "error", err)
	}

	for _, decrease := range decreases {
		// generate amount of sequences for amount decreased/increased
		// start sequence with decrease
		for range decrease.AmountChanged {
			decrease.AmountFaked++
			db.Transaction(func(tx *gorm.DB) error {
				//get increase for decrease within window
				increase := database.HistoricalStationData{}
				db.Where("event_time_stamp between ? and ? AND topic_name = 'station_occupation_increased' AND amount_changed > amount_faked",
					decrease.EventTimeStamp.Add(minDuration), decrease.EventTimeStamp.Add(minDuration+windowSize)).Order("random()").Limit(1).Find(&increase)
				err := generate(db, increase, decrease)
				if err != nil {
					slog.Warn("Fake event transaction failed", "error", err)
					decrease.AmountFaked--
					err := tx.Rollback().Error
					if err != nil {
						return err
					}
				}
				return nil
			})
		}
	}
}

// generates bike event sequence
func generate(db *gorm.DB, increase database.HistoricalStationData, decrease database.HistoricalStationData) error {
	slog.Debug("Generating event sequence")
	decreaseStation, err := database.GetStationById(decrease.Uuid, db)
	if err != nil {
		return err
	}

	// number of minutes bike is reserved before it is picked up
	startoffset, err := getStartOffset(db, decrease)
	if err != nil {
		return err
	}

	startTime := decrease.EventTimeStamp.Add(-startoffset).UTC()
	pickedUpTime := decrease.EventTimeStamp.UTC()
	defectTime := pickedUpTime.Add(helper.RandMinutes(60*5, 60)).UTC()
	immobilizedTime := defectTime.Add(helper.RandMinutes(60*5, 2)).UTC()
	abandonedTime := immobilizedTime.Add(helper.RandMinutes(60*5, 2)).UTC()
	fakeEndTime := abandonedTime.Add(helper.RandMinutes(60*5, 60)).UTC()
	endTime := increase.EventTimeStamp.UTC()

	bike, user, err := getAvailableBikeAndUser(db, startTime)
	if err != nil {
		return err
	}

	// bike reserved
	database.BikeReservedEvent(bike, startTime, decreaseStation, user, db)

	// bike picked up
	bike.IsReturned = sql.NullBool{Bool: false, Valid: true}
	database.BikePickedUpEvent(bike, pickedUpTime, decreaseStation, user, db)

	// bike not returned
	if increase == (database.HistoricalStationData{}) {
		slog.Info("Bike not returned", "station", decreaseStation.Name)
		endTime = fakeEndTime
		// chance bikedefect
		if rand.Float64() < chanceDefectNotReturned {
			bike.IsDefect = sql.NullBool{Bool: true, Valid: true}
			database.BikeDefectEvent(bike, defectTime, user, defects[rand.IntN(len(defects))], db)

			// chance bike immobilized
			if rand.Float64() < chanceImmobilized {
				bike.IsImmobilized = sql.NullBool{Bool: true, Valid: true}
				database.BikeImmobilizedEvent(bike, immobilizedTime, db)
			}
		}
		// bike abandoned
		bike.IsAbandoned = sql.NullBool{Bool: true, Valid: true}
		database.BikeAbandonedEvent(bike, abandonedTime, user, db)

		db.Model(&decrease).Update("amount_faked", decrease.AmountFaked)
	} else { // bike returned
		slog.Info("Bike returned", "station", decreaseStation.Name)
		increase.AmountFaked++
		increaseStation, err := database.GetStationById(increase.Uuid, db)
		if err != nil {
			return err
		}
		// after capacity decrease
		delta := endTime.Sub(decrease.EventTimeStamp)
		r := rand.Float64()
		if r == 0 {
			r = 0.5
		}
		offset := time.Second * time.Duration(r*delta.Seconds())
		// chance bike defect
		if rand.Float64() < chanceDefect {
			bike.IsDefect = sql.NullBool{Bool: true, Valid: true}
			database.BikeDefectEvent(bike, startTime.Add(offset), user, defects[rand.IntN(len(defects))], db)
		}

		// same time as capacity increase
		// bike returned
		bike.IsReturned = sql.NullBool{Bool: true, Valid: true}
		database.BikeReturnedEvent(bike, endTime, increaseStation, user, db)

		db.Model(&decrease).Update("amount_faked", decrease.AmountFaked)
		db.Model(&increase).Update("amount_faked", increase.AmountFaked)
	}

	user.IsAvailableTimestamp = sql.NullTime{Time: endTime, Valid: true}
	bike.InUseTimestamp = sql.NullTime{Time: endTime, Valid: true}
	database.UpdateBike(db, &bike)
	database.UpdateUser(db, &user)

	return nil
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

// returns an available bike and user
func getAvailableBikeAndUser(db *gorm.DB, startTime time.Time) (database.Bike, database.User, error) {
	startTimeString := startTime.Format(format)
	//get available bike not in use
	bike, user := database.Bike{}, database.User{}
	if result := db.Where("(is_defect = false AND is_immobilized = false AND is_abandoned = false AND is_in_storage = false AND is_returned = true) AND "+
		"(in_use_timestamp is null OR (extract(epoch from ? - in_use_timestamp)/60 > 0) AND extract(epoch from ? - created_at)/60 > 0)",
		startTimeString, startTimeString).Order("random()").Limit(1).Find(&bike); result.RowsAffected == 0 {
		// If no bike available, clean up bikes
		if err := database.BikeCleanUp(db); err != nil {
			return database.Bike{}, database.User{}, err
		}
		if result := db.Where("(is_defect = false AND is_immobilized = false AND is_abandoned = false AND is_in_storage = false AND is_returned = true) AND "+
			"(in_use_timestamp is null OR (extract(epoch from ? - in_use_timestamp)/60 > 0) AND extract(epoch from ? - created_at)/60 > 0)",
			startTimeString, startTimeString).Order("random()").Limit(1).Find(&bike); result.RowsAffected == 0 {
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
	result := db.Where("is_available_timestamp is null OR (extract(epoch from ? - is_available_timestamp)/60 > 0 AND extract(epoch from ? - created_at)/60 > 0)", startTimeString, startTimeString).Order("random()").Limit(1).Find(&user)
	if result.RowsAffected == 0 {
		slog.Info("No available user found, creating user")
		newUser, err := database.CreateUser(db, startTime.Add(-time.Minute*30))
		if err != nil {
			return database.Bike{}, database.User{}, err
		}
		user = *newUser
	}

	slog.Info("Bike & User selected", "bike", bike.Id, "user", user.Id)
	return bike, user, nil
}

// gpt generated defects
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
