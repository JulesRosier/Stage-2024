package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"stage2024/pkg/helper"
	"time"

	"gorm.io/gorm"
)

// Cleans up bikes in a transaction
func BikeCleanUpTransaction(db *gorm.DB) {
	slog.Debug("Starting bike cleanup Transaction")
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := BikeCleanUp(tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		slog.Warn("error cleaning up bikes", "error", err)
	}
}

// Makes unavailable bikes available again
func BikeCleanUp(db *gorm.DB) error {
	slog.Info("Cleaning up bikes")
	bikes := []Bike{}
	db.Where("(is_defect = true OR is_immobilized = true OR is_abandoned = true OR is_in_storage = true OR is_returned = false) AND in_use_timestamp IS NOT NULL AND extract(epoch from ? - in_use_timestamp)/60 > 10", time.Now().UTC().Format("2006-01-02 15:04:05.999999-07")).Find(&bikes)
	for _, bike := range bikes {
		slog.Info("Cleaning up bike", "bike", bike.Id)
		times := getEventTimes(bike, 3)

		bike.IsAbandoned = sql.NullBool{Bool: false, Valid: true}

		if !bike.IsInStorage.Bool {
			if err := BikeStoredEvent(bike, times[0], db); err != nil {
				return fmt.Errorf("error sending event, error : %v", err)
			}
		}
		bike.IsInStorage = sql.NullBool{Bool: false, Valid: true}

		if bike.IsDefect.Bool || bike.IsImmobilized.Bool {
			if err := BikeRepairedEvent(bike, times[1], db); err != nil {
				return fmt.Errorf("error sending event, error : %v", err)
			}
			bike.IsDefect = sql.NullBool{Bool: false, Valid: true}
			bike.IsImmobilized = sql.NullBool{Bool: false, Valid: true}
		}
		if err := BikeDeployedEvent(bike, times[2], db); err != nil {
			return fmt.Errorf("error sending event, error : %v", err)
		}

		bike.IsReturned = sql.NullBool{Bool: true, Valid: true}
		bike.InUseTimestamp = sql.NullTime{Time: times[2], Valid: true}
		if err := db.Save(&bike).Error; err != nil {
			return err
		}

	}
	return nil
}

// Get event times for bike cleanup function
func getEventTimes(bike Bike, amount int32) []time.Time {
	times := []time.Time{}
	previous := bike.InUseTimestamp.Time
	now := time.Now().UTC()
	for range amount {
		delta := now.Sub(previous)
		r := helper.GetRandomNumber()
		offset := time.Second * time.Duration(r*delta.Seconds())
		eventTime := previous.Add(offset)
		times = append(times, eventTime)
		previous = eventTime
	}

	return times
}
