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
		// TODO hier geen event van?, kan wel
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
		slog.Info("Hoe kan een fiets opeens wel of niet elektrisch worden?")
		// TODO waarschijnlijk niets, want is geen event, zal niet veranderen
	case "IsImmobilized":
		// TODO generated
	case "IsAbandoned":
		// TODO generated
	case "IsAvailable":
		// TODO
	case "IsInStorage":
		isInStorageChange(bike, change) // will never change on its own, generated
	case "IsReserved":
		reservedChange(bike, change) //TODO
	case "IsDefect":
		// TODO generated
	}
}

//generated values for IsImmobilized, IsAbandoned, IsInStorage, IsDefect
