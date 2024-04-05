package events

import (
	"fmt"
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
)

// Reacts to change in database, sends appropiate events
func (ec EventClient) ChangeDetected(change helper.Change) {
	//set kafka client to global var

	switch change.Table {
	case "Bike":
		ec.bikechange(change)
	case "Station":
		ec.stationchange(change)
	}
}

// Selects right events to send based on change for stations
func (ec EventClient) stationchange(change helper.Change) {
	station, err := database.GetStationById(change.Id)

	helper.MaybeDieErr(err)

	switch change.Column {
	case "lat", "lon":
		// kan denkik niet?
	case "MaxCapacity":
		// TODO hier geen event van?, kan wel
	case "Occupation":
		ec.occupationChange(station, change)
	case "IsActive":
		ec.activeChange(station, change)
	case "Created":
		ec.created(station)
	}

}

// Selects right events to send based on change for bikes
func (ec EventClient) bikechange(change helper.Change) {
	bike, err := database.GetBikeById(change.Id)
	helper.MaybeDieErr(err)

	switch change.Column {
	case "Lat", "Lon":
		locationchange(bike) // TODO, no event for this
	case "IsImmobilized":
		ec.immobilizedChange(bike, change) // generated
	case "IsAbandoned":
		user := user(change)
		ec.abandonedChange(bike, user, change) // generated
	case "IsAvailable":
		//TODO Real event, needs generated station and user
		// what event is this?
		change.Station_id = "00000000-test-test-test-000000000000"
		change.User_id = "00f9c353-ec33-42c0-81fa-125fdf66c6a3"
		station, user := stationAndUser(change)
		availableChange(bike, station, user, change) //real
	case "IsInStorage":
		ec.isInStorageChange(bike, change) // generated
	case "IsReserved":
		//TODO real event, needs generated station and user
		//start reserved sequence

		ec.reservedChange(bike, change) //real
	case "IsDefect":
		user, defect := userdefect(change)
		ec.defectChange(bike, user, change, defect) //generated
	case "PickedUp":
		station, user := stationAndUser(change)
		ec.pickedUpChange(bike, station, user)
	case "Returned":
		station, user := stationAndUser(change)
		ec.returnedChange(bike, station, user)
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
