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
		err := bikechange(change, db)
		return err
	case "Station":
		err := stationchange(change, db)
		return err
	case "User":
		err := userchange(change, db)
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
	case "Created":
		err := createdStation(station, change, db)
		return err
	}

	return err
}

// Selects right events to send based on change for bikes
func bikechange(change helper.Change, db *gorm.DB) error {
	bike, err := GetBikeById(change.Id, db)
	helper.MaybeDieErr(err)

	switch change.Column {
	case "IsImmobilized":
		err := immobilizedChange(bike, change, db) // generated
		return err
	case "IsAbandoned":
		err := abandonedChange(bike, change, db) // generated
		return err
	case "IsInStorage":
		err := isInStorageChange(bike, change, db) // generated
		return err
	case "IsReserved":
		//start reserved sequence
		err := reservedChange(bike, change, db) //real
		return err
	case "IsDefect":
		err := defectChange(bike, change, db) //generated
		return err
	case "PickedUp":
		err := pickedUpChange(bike, change, db)
		return err
	case "IsReturned":
		err := returnedChange(bike, change, db)
		return err
	}
	return nil
}

//generated values for IsImmobilized, IsAbandoned, IsInStorage, IsDefect

func userchange(change helper.Change, db *gorm.DB) error {
	user, err := GetUserById(change.Id, db)
	helper.MaybeDieErr(err)

	switch change.Column {
	case "Created":
		err := createdUser(user, change, db)
		return err
	}
	return nil
}
