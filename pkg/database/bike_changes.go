package database

import (
	"fmt"
	"log/slog"
	"stage2024/pkg/helper"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/common"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// Sends BikeImmobilized event. Needs bike and change.
func BikeImmobilizedEvent(bike Bike, change helper.Change, db *gorm.DB) error {
	//TODO: add event for bike mobilized
	//bike immobilized
	slog.Debug("Bike immobilized, sending event...", "bike", bike.Id)
	protostruct := &bikes.BikeImmobilized{
		TimeStamp: timestamppb.New(change.EventTime),
		Bike: &bikes.LocationBike{
			Bike: bike.IntoId(),
			Location: &common.Location{
				Latitude:  bike.Lat,
				Longitude: bike.Lon,
			},
		},
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}
	return nil
}

// Sends BikeAbandoned event. Needs user
func BikeAbandonedEvent(bike Bike, change helper.Change, db *gorm.DB) error {
	//bike abandoned
	user, err := user(change, db)
	if err != nil {
		return err
	}
	slog.Debug("Bike abandoned, sending event...", "bike", bike.Id)
	protostruct := &bikes.BikeAbandoned{
		TimeStamp: timestamppb.New(change.EventTime),
		Bike: &bikes.AbandonedBike{
			Bike: bike.IntoId(),
			Location: &common.Location{
				Latitude:  bike.Lat,
				Longitude: bike.Lon,
			},
		},
		User: user.IntoId(),
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}
	return nil
}

// sends BikeBroughtOut and BikeStored events. Needs bike and change.
func BikeInStorageEvent(bike Bike, change helper.Change, db *gorm.DB) error {
	if change.NewValue == "false" {
		slog.Info("Bike brought out, sending event...", "bike", bike.Id)
		protostruct := &bikes.BikeDeployed{
			TimeStamp: timestamppb.New(change.EventTime),
			Bike:      bike.IntoId(),
		}
		if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
			return err
		}
		return nil
	}
	if change.NewValue == "true" {
		slog.Debug("Bike stored, sending event...", "bike", bike.Id)
		protostruct := &bikes.BikeStored{
			TimeStamp: timestamppb.New(change.EventTime),
			Bike:      bike.IntoId(),
		}
		if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
			return err
		}
	}
	return nil
}

// Sends BikeReserved event. Needs station and user.
func BikeReserved(bike Bike, change helper.Change, db *gorm.DB) error {
	slog.Debug("Bike reserved, sending event...", "bike", bike.Id)
	station, user, err := stationAndUser(change, db)
	if err != nil {
		return err
	}
	protostruct := &bikes.BikeReserved{
		TimeStamp: timestamppb.New(change.EventTime),
		Bike: &bikes.Bike{
			Bike:       bike.IntoId(),
			IsElectric: bike.IsElectric.Bool,
		},
		Station: station.IntoId(),
		User:    user.IntoId(),
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}

	return nil
}

// Sends BikeDefectReported event. Needs bike, user and defect.
func BikeDefectEvent(bike Bike, change helper.Change, db *gorm.DB) error {
	user, defect, err := userDefect(change, db)
	if err != nil {
		return err
	}
	slog.Debug("Bike defect, sending event...", "bike", bike.Id)
	protostruct := &bikes.BikeDefectReported{
		TimeStamp: timestamppb.New(change.EventTime),
		Bike: &bikes.DefectBike{
			Bike: bike.IntoId(),
			Location: &common.Location{
				Latitude:  bike.Lat,
				Longitude: bike.Lon,
			},
			IsElectric: bike.IsElectric.Bool,
		},
		User:           user.IntoId(),
		ReportedDefect: defect,
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}

	return nil
}

// Sends BikePickedUp event. Needs bike, station, user.
func BikePickedUpEvent(bike Bike, change helper.Change, db *gorm.DB) error {
	slog.Debug("Bike picked up, sending event...", "bike", bike.Id)
	station, user, err := stationAndUser(change, db)
	if err != nil {
		return err
	}
	protostruct := &bikes.BikePickedUp{
		TimeStamp: timestamppb.New(change.EventTime),
		Bike:      bike.IntoId(),
		Station:   station.IntoId(),
		User:      user.IntoId(),
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}

	return nil
}

func BikeReturnedEvent(bike Bike, change helper.Change, db *gorm.DB) error {
	slog.Debug("Bike returned, sending event...", "bike", bike.Id)
	station, user, err := stationAndUser(change, db)
	if err != nil {
		return err
	}
	protostruct := &bikes.BikeReturned{
		TimeStamp: timestamppb.New(change.EventTime),
		Bike:      bike.IntoId(),
		Station:   station.IntoId(),
		User:      user.IntoId(),
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}

	return nil
}

func BikeCreatedEvent(bike *Bike, db *gorm.DB) error {
	slog.Debug("Bike created, sending event...", "bike", bike.ID)
	now := timestamppb.Now()
	protostruct := &bikes.BikeDeployed{
		TimeStamp: now,
		Bike:      bike.IntoId(),
	}
	if err := createOutboxRecord(now, protostruct, db); err != nil {
		return err
	}
	return nil
}

// checks if station and user are present and returns them
func stationAndUser(change helper.Change, db *gorm.DB) (Station, User, error) {
	station, err := GetStationById(change.StationId, db)
	if err != nil {
		return Station{}, User{}, err
	}
	user, err := GetUserById(change.UserId, db)
	if err != nil {
		return Station{}, User{}, err
	}
	return station, user, nil
}

// checks if user is present and returns it
func user(change helper.Change, db *gorm.DB) (User, error) {
	user, err := GetUserById(change.UserId, db)
	return user, err
}

// checks if user and defect are present and returns them
func userDefect(change helper.Change, db *gorm.DB) (User, string, error) {
	user, err := GetUserById(change.UserId, db)
	if change.Defect == "" {
		err := fmt.Errorf("defect change detected, but no defect was given. defect=%v", change.Defect)
		return user, change.Defect, err
	}
	return user, change.Defect, err
}
