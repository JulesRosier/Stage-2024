package main

import (
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
	"stage2024/pkg/kafka"
	"stage2024/pkg/protogen/common"
	"stage2024/pkg/protogen/stations"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var kc *kafka.Client

func ChangeDetected(k *kafka.Client, change helper.Change) {
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
	switch change.Column {
	case "lat", "lon":
		// TODO
	case "MaxCapacity":
		// TODO
	case "Occupation":
		occupationchange(change.Id)
	case "IsActive":
		// TODO
	}
}

func bikechange(change helper.Change) {
	slog.Debug("Bike change detected")
	switch change.Column {
	case "Lat", "Lon":
		// TODO
	case "IsElectric":
		// TODO
	case "IsImmobilized":
		// TODO
	case "IsAbandoned":
		// TODO
	case "IsAvailable":
		// TODO
	case "IsInStorage":
		// TODO
	case "IsReserved":
		// TODO
	case "IsDefect":
		// TODO
	}
}

func occupationchange(id string) {
	db := database.GetDb()
	var station database.Station
	db.Where("id = ?", id).First(&station)

	//station full
	if station.Occupation == station.MaxCapacity {
		slog.Info("Station is full, sending event")

		err := kc.Produce(&stations.StationCapacityExhausted{
			TimeStamp: timestamppb.Now(),
			Station: &stations.StationIdentification{
				Id: station.Id,
				Location: &common.Location{
					Latitude:  station.Lat,
					Longitude: station.Lon,
				},
				Name: station.Name,
			},
			MaxCapacity: station.MaxCapacity,
		})
		helper.MaybeDieErr(err)
	}
}
