package database

import (
	"log/slog"
	"stage2024/pkg/helper"

	"gorm.io/gorm"
)

// Reacts to change in database, sends appropiate events
func ChangeDetected(change helper.Change, db *gorm.DB) error {
	//set kafka client to global var

	switch change.Table {
	case "Bike":
		slog.Warn("DEPRECATED, don't send events from db update")
	case "Station":
		err := stationchange(change, db)
		return err
	}
	return nil
}

// Selects right events to send based on change for stations
func stationchange(change helper.Change, db *gorm.DB) error {
	station, err := GetStationById(change.Id, db)
	if err != nil {
		return err
	}

	switch change.Column {
	case "lat", "lon":
		// kan denkik niet?
	case "MaxCapacity":
		// TODO hier geen event van?, kan wel
	case "Occupation":
		err := occupationChange(station, change, db)
		return err
	case "IsActive":
		err := activeChange(station, change, db)
		return err
	}

	return err
}
