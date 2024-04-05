package events

import (
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
		ec.abandonedChange(bike, change) // generated
	case "IsAvailable":
		//TODO Real event, needs generated station and user
		// what event is this?
		availableChange(bike, change) //real
	case "IsInStorage":
		ec.isInStorageChange(bike, change) // generated
	case "IsReserved":
		//start reserved sequence
		ec.reservedChange(bike, change) //real
	case "IsDefect":
		ec.defectChange(bike, change) //generated
	case "PickedUp":
		ec.pickedUpChange(bike, change)
	case "Returned":

		ec.returnedChange(bike, change)
	}

}

//generated values for IsImmobilized, IsAbandoned, IsInStorage, IsDefect
