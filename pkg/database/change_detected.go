package database

import (
	"stage2024/pkg/helper"

	"gorm.io/gorm"
)

// Reacts to change in database, sends appropiate events
func ChangeDetected(change helper.Change, db *gorm.DB) error {
	//set kafka client to global var

	switch change.Table {
	case "Bike":
		err := bikechange(change)
		return err
	case "Station":
		err := stationchange(change, db)
		return err
	}
	return nil
}

// Selects right events to send based on change for stations
func stationchange(change helper.Change, db *gorm.DB) error {
	station, err := GetStationById(change.Id)
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
	case "Created":
		err := created(station, db)
		return err
	}

	return err
}

// Selects right events to send based on change for bikes
func bikechange(change helper.Change) error {
	bike, err := GetBikeById(change.Id)
	helper.MaybeDieErr(err)

	switch change.Column {
	case "Lat", "Lon":
		err := locationchange(bike) // TODO, no event for this
		return err
	case "IsImmobilized":
		err := immobilizedChange(bike, change) // generated
		return err
	case "IsAbandoned":
		err := abandonedChange(bike, change) // generated
		return err
	case "IsAvailable":
		//TODO Real event, needs generated station and user
		// what event is this?
		err := availableChange(bike, change) //real
		return err
	case "IsInStorage":
		err := isInStorageChange(bike, change) // generated
		return err
	case "IsReserved":
		//start reserved sequence
		err := reservedChange(bike, change) //real
		return err
	case "IsDefect":
		err := defectChange(bike, change) //generated
		return err
	case "PickedUp":
		err := pickedUpChange(bike, change)
		return err
	case "Returned":
		err := returnedChange(bike, change)
		return err
	}
	return nil
}

//generated values for IsImmobilized, IsAbandoned, IsInStorage, IsDefect
