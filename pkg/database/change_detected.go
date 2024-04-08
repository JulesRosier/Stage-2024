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
		// bikechange(change)
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

// // Selects right events to send based on change for bikes
// func bikechange(change helper.Change) {
// 	bike, err := GetBikeById(change.Id)
// 	helper.MaybeDieErr(err)

// 	switch change.Column {
// 	case "Lat", "Lon":
// 		locationchange(bike) // TODO, no event for this
// 	case "IsImmobilized":
// 		immobilizedChange(bike, change) // generated
// 	case "IsAbandoned":
// 		abandonedChange(bike, change) // generated
// 	case "IsAvailable":
// 		//TODO Real event, needs generated station and user
// 		// what event is this?
// 		availableChange(bike, change) //real
// 	case "IsInStorage":
// 		isInStorageChange(bike, change) // generated
// 	case "IsReserved":
// 		//start reserved sequence
// 		reservedChange(bike, change) //real
// 	case "IsDefect":
// 		defectChange(bike, change) //generated
// 	case "PickedUp":
// 		pickedUpChange(bike, change)
// 	case "Returned":
// 		returnedChange(bike, change)
// 	}

// }

//generated values for IsImmobilized, IsAbandoned, IsInStorage, IsDefect
