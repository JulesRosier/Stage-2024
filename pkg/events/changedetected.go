package events

import (
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
	"stage2024/pkg/kafka"
)

var kc *kafka.Client

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

func stationchange(change helper.Change) {
	slog.Debug("Station change detected")
	station := database.GetStationById(change.Id)
	switch change.Column {
	case "lat", "lon":
		// kan denkik niet?
	case "MaxCapacity":
		// TODO
	case "Occupation":
		occupationChange(station, change)
	case "IsActive":
		activeChange(station, change)
	}
}

func bikechange(change helper.Change) {
	slog.Debug("Bike change detected")
	bike := database.GetBikeById(change.Id)
	switch change.Column {
	case "Lat", "Lon":
		locationchange(bike) // TODO
	case "IsElectric":
		// TODO waarschijnlijk niets, want is geen event, zal niet veranderen
	case "IsImmobilized":
		// TODO gen
	case "IsAbandoned":
		// TODO gen
	case "IsAvailable":
		availableChange(bike, change) //TODO
	case "IsInStorage":
		// TODO gen
	case "IsReserved":
		reservedChange(bike, change) //TDOD
	case "IsDefect":
		// TODO gen
	}
}
