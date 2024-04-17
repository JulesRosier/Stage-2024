package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"stage2024/pkg/helper"
	"time"

	"gorm.io/gorm"
)

// makes unavailable bikes available again
func BikeCleanUp(db *gorm.DB) error {
	slog.Debug("Cleaning up bikes")
	bikes := []Bike{}
	db.Where("(is_defect = true OR is_immobilized = true OR is_abandoned = true OR is_in_storage = true OR is_returned = false) AND in_use_timestamp IS NOT NULL AND extract(epoch from ? - in_use_timestamp)/60 > 0", time.Now().UTC().Format("2006-01-02 15:04:05.999999-07")).Find(&bikes)
	for _, bike := range bikes {

		eventTime1, eventTime2, eventTime3 := getEventTimes(bike)

		change := helper.Change{}
		bike.IsDefect = sql.NullBool{Bool: false, Valid: true}
		bike.IsImmobilized = sql.NullBool{Bool: false, Valid: true}
		bike.IsAbandoned = sql.NullBool{Bool: false, Valid: true}
		if !bike.IsInStorage.Bool {
			change.EventTime = eventTime1
			if err := BikeStoredEvent(bike, change, db); err != nil {
				slog.Warn("Error sending event", "error", err)
			}
		}
		bike.IsInStorage = sql.NullBool{Bool: false, Valid: true}
		if err := BikeRepairedEvent(bike, eventTime2, db); err != nil {
			slog.Warn("Error sending event", "error", err)
		}
		if err := BikeDeployedEvent(bike, eventTime3, db); err != nil {
			return fmt.Errorf("error sending event, error : %v", err)
		}

		bike.IsReturned = sql.NullBool{Bool: true, Valid: true}

		db.Save(&bike)
	}
	return nil
}

// Get event times for bike cleanup function
func getEventTimes(bike Bike) (time.Time, time.Time, time.Time) {
	times := []time.Time{}
	previous := bike.InUseTimestamp.Time
	now := time.Now().UTC()
	for range 3 {
		delta := now.Sub(previous)
		r := rand.Float64()
		if r == 0 {
			r = 0.5
		}
		offset := time.Second * time.Duration(r*delta.Seconds())
		eventTime := previous.Add(offset)
		times = append(times, eventTime)
		previous = eventTime
		slog.Debug("Event time", "time", eventTime)
	}

	return times[0], times[1], times[2]
}