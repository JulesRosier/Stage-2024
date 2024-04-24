package events

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

const (
	minDuration             = time.Minute * 10
	windowSize              = time.Minute * 30
	chanceDefect            = 0.02
	chanceDefectNotReturned = 0.5
	chanceImmobilized       = 0.5
	format                  = "2006-01-02 15:04:05.999999-07"
	fakeStationId           = "ab448be3-5c90-43ce-8c37-74f929ec016f"
)

// BikeEventGen generates bike events based on historical station occupation changes in the database
func BikeEventGen(db *gorm.DB) {
	slog.Info("Generating bike events")

	nowUtc := time.Now().UTC().Format(format)
	decreases := []database.HistoricalStationData{}
	//get all decreases older than minDuration+windowSize with amount faked < amount changed
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
			err := db.Transaction(func(tx *gorm.DB) error {
				//get increase for decrease within window
				increase := database.HistoricalStationData{}
				db.Where("event_time_stamp between ? and ? AND topic_name = 'station_occupation_increased' AND amount_changed > amount_faked",
					decrease.EventTimeStamp.Add(minDuration), decrease.EventTimeStamp.Add(minDuration+windowSize)).Order("random()").Limit(1).Find(&increase)
				err := generate(db, increase, decrease)
				return err
			})
			if err != nil {
				slog.Warn("Fake event transaction failed", "error", err)
				decrease.AmountFaked--
			}
		}
	}

	// get all unchecked increases older than minDuration+windowSize
	increases := []database.HistoricalStationData{}
	err = db.Where("extract(epoch from ? - event_time_stamp)/60 >= ? and topic_name = 'station_occupation_increased' and amount_changed > amount_faked",
		nowUtc, (minDuration + windowSize).Minutes()).Order("id asc").Find(&increases).Error
	if err != nil {
		slog.Warn("Failed to get increases", "error", err)
	}

	// generate amount of sequences for amount decreased/increased
	slog.Info("Generating fake decreases for increase")
	for _, increase := range increases {
		// start sequence with increase
		db.Transaction(func(tx *gorm.DB) error {
			err := generateIncrease(db, increase)
			return err
		})
		if err != nil {
			slog.Warn("Error generating event sequence for increase")
		}
	}
}

// generates fake decrease events for increase
func generateIncrease(db *gorm.DB, increase database.HistoricalStationData) error {
	// bike needs to get picked up from fake station...
	station := database.Station{}
	result := db.Model(&station).Where("id = ?", fakeStationId).First(&station)
	if result.RowsAffected == 0 {
		slog.Debug("Fake station not found", "station", fakeStationId)
		database.MakeFakeStation(db, fakeStationId)
		db.Model(&station).Where("id = ?", fakeStationId).First(&station)
	}
	if station.Occupation < 50 {
		station.Occupation = station.MaxCapacity - 1
		db.Save(&station)
	}

	for range increase.AmountChanged - increase.AmountFaked {
		// Timestamp random amount of minutes inside windowSize + minDuration
		randomMinutes := rand.IntN(int(windowSize) + int(minDuration))
		if randomMinutes < int(minDuration) {
			randomMinutes = int(minDuration)
		}
		timestamp := increase.EventTimeStamp.Add(-time.Duration(randomMinutes))
		if timestamp.Before(station.CreatedAt) {
			slog.Info("Increase too close to station creation")
			db.Model(&increase).Update("amount_faked", increase.AmountChanged)
			return nil
		}
		database.OccupationChange(station, helper.Change{Id: station.Id, OldValue: fmt.Sprint(station.Occupation), NewValue: fmt.Sprint(station.Occupation - 1)}, timestamp, db)
		station.Occupation = station.Occupation - 1
		slog.Debug("generated fake decrease")
		//get generated decrease
		decrease := database.HistoricalStationData{}
		result = db.Model(&decrease).Where("uuid = ? and amount_changed > amount_faked and topic_name = 'station_occupation_decreased'", station.Id).First(&decrease)
		if result.RowsAffected == 0 {
			slog.Info("No decrease found for station", "station", station.Id)
		}

		err := generate(db, increase, decrease)
		if err != nil {
			increase.AmountFaked--
			return err
		}
	}

	return nil
}

// generates bike event sequence
func generate(db *gorm.DB, increase database.HistoricalStationData, decrease database.HistoricalStationData) error {
	slog.Debug("Generating event sequence")
	decreaseStation, err := database.GetStationById(decrease.Uuid, db)
	if err != nil {
		return err
	}

	// number of minutes bike is reserved before it is picked up
	startOffset, err := getStartOffset(db, decrease)
	if err != nil {
		return err
	}

	startTime := decrease.EventTimeStamp.Add(-startOffset).UTC()
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
		// chance bike defect
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
		delta := endTime.Sub(pickedUpTime)
		r := helper.GetRandomNumber()
		offset := time.Second * time.Duration(r*delta.Seconds())
		// chance bike defect
		if rand.Float64() < chanceDefect {
			bike.IsDefect = sql.NullBool{Bool: true, Valid: true}
			database.BikeDefectEvent(bike, pickedUpTime.Add(offset), user, defects[rand.IntN(len(defects))], db)
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
func getStartOffset(db *gorm.DB, station database.HistoricalStationData) (time.Duration, error) {
	stationCreated := database.HistoricalStationData{}
	err := db.Where("uuid = ? AND topic_name = 'station_created'", station.Uuid).Limit(1).Find(&stationCreated).Error
	if err != nil {
		slog.Info("Station not found", "station", station.Uuid, "error", err)
		return 0, err
	}

	delta := station.EventTimeStamp.Sub(stationCreated.EventTimeStamp)
	if delta > windowSize {
		delta = windowSize
	}

	r := helper.GetRandomNumber()

	offset := time.Second * time.Duration(r*delta.Seconds())

	return offset, err
}

// returns an available bike and user
func getAvailableBikeAndUser(db *gorm.DB, startTime time.Time) (database.Bike, database.User, error) {
	startTimeString := startTime.Format(format)
	//get available bike not in use
	bike, user := database.Bike{}, database.User{}
	if result := db.Where("(is_defect = false AND is_immobilized = false AND is_abandoned = false AND is_in_storage = false AND is_returned = true) AND "+
		"((in_use_timestamp is null OR (extract(epoch from ? - in_use_timestamp)/60 > 0)) AND extract(epoch from ? - created_at)/60 > 0)",
		startTimeString, startTimeString).Order("random()").Limit(1).Find(&bike); result.RowsAffected == 0 {
		// If no bike available, clean up bikes
		if err := database.BikeCleanUp(db); err != nil {
			return database.Bike{}, database.User{}, err
		}
		if result := db.Where("(is_defect = false AND is_immobilized = false AND is_abandoned = false AND is_in_storage = false AND is_returned = true) AND "+
			"(in_use_timestamp is null OR (extract(epoch from ? - in_use_timestamp)/60 > 0) AND extract(epoch from ? - created_at)/60 > 0)",
			startTimeString, startTimeString).Order("random()").Limit(1).Find(&bike); result.RowsAffected == 0 {
			// If still no bike available, create bike
			slog.Debug("No available bike found, creating bike")
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
		slog.Debug("No available user found, creating user")
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
	"Squeaky seat post",
	"Slipping seat clamp",
	"Twisted handlebar grips",
	"Rattling fender",
	"Misaligned derailleur",
	"Sticky shift levers",
	"Worn-out cassette",
	"Scratched frame paint",
	"Broken spoke nipples",
	"Frayed brake cables",
}
