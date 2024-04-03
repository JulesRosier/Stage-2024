package events

import (
	"fmt"
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
	"stage2024/pkg/kafka"
)

var kc *kafka.Client

// Reacts to change in database, sends appropiate events
func ChangeDetected(k *kafka.Client, change helper.Change) {
	//set kafka client to global var
	kc = k
	switch change.Table {
	case "Bike":
		bikechange(change)
	case "Station":
		stationchange(change)
	}
}

// Selects right events to send based on change for stations
func stationchange(change helper.Change) {
	station, err := database.GetStationById(change.Id)
	helper.MaybeDieErr(err)

	switch change.Column {
	case "lat", "lon":
		// kan denkik niet?
	case "MaxCapacity":
		// TODO hier geen event van?, kan wel
	case "Occupation":
		occupationChange(station, change)
	case "IsActive":
		activeChange(station, change)
	case "Created":
		created(station)
	}

}

// Selects right events to send based on change for bikes
func bikechange(change helper.Change) {
	bike, err := database.GetBikeById(change.Id)
	helper.MaybeDieErr(err)

	switch change.Column {
	case "Lat", "Lon":
		locationchange(bike) // TODO, no event for this
	case "IsImmobilized":
		immobilizedChange(bike, change) // generated
	case "IsAbandoned":
		user := user(change)
		abandonedChange(bike, user, change) // generated
	case "IsAvailable":
		//TODO Real event, needs generated station and user
		// what event is this?
		change.Station_id = "61a69174-89f4-412c-9fa7-fe129b1ced98"
		change.User_id = "00f9c353-ec33-42c0-81fa-125fdf66c6a3"
		station, user := stationAndUser(change)
		availableChange(bike, station, user, change) //real
	case "IsInStorage":
		isInStorageChange(bike, change) // generated
	case "IsReserved":
		//TODO real event, needs generated station and user
		//start reserved sequence
		if bike.IsReserved.Bool {
			change = startReservedsequence(bike, change)
		}
		station, user := stationAndUser(change)
		reservedChange(bike, station, user, change) //real
	case "IsDefect":
		change.User_id = "00f9c353-ec33-42c0-81fa-125fdf66c6a3"
		change.Defect = "test defect......."
		user, defect := userdefect(change)
		defectChange(bike, user, change, defect) //generated
	}
}

//generated values for IsImmobilized, IsAbandoned, IsInStorage, IsDefect

// checks if station and user are present and returns them
func stationAndUser(change helper.Change) (database.Station, database.User) {
	station, err := database.GetStationById(change.Station_id)
	helper.MaybeDieErr(err)
	user, err := database.GetUserById(change.User_id)
	helper.MaybeDieErr(err)
	return station, user
}

// checks if user is present and returns it
func user(change helper.Change) database.User {
	user, err := database.GetUserById(change.User_id)
	helper.MaybeDieErr(err)
	return user
}

// checks if user and defect are present and returns them
func userdefect(change helper.Change) (database.User, string) {
	user, err := database.GetUserById(change.User_id)
	helper.MaybeDieErr(err)
	if change.Defect == "" {
		helper.Die(fmt.Errorf("defect change detected, but no defect was given. defect=%v", change.Defect))
	}
	return user, change.Defect
}
