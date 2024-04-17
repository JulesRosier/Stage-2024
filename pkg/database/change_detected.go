package database

import (
	"stage2024/pkg/helper"

	"gorm.io/gorm"
)

// Selects right events to send based on change for stations
func ChangeDetected(change helper.Change, db *gorm.DB) error {
	station, err := GetStationById(change.Id, db)
	if err != nil {
		return err
	}

	switch change.Column {
	case "lat", "lon":
		// kan denkik niet?
	case "Occupation":
		err := occupationChange(station, change, db)
		return err
	case "IsActive":
		err := activeChange(station, change, db)
		return err
	}

	return err
}
