package database

import (
	"fmt"
	"log/slog"
	"stage2024/pkg/helper"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/common"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// TODO add events for bike location changes
func locationchange(bike Bike) error {
	// slog.Info("Bike location changed, sending event...", "bike", bike.OpenDataId)
	return nil
}

// Sends BikeImmobilized event. Needs bike and change.
func immobilizedChange(bike Bike, change helper.Change) error {
	//TODO add event for bike mobilized
	//bike immobilized
	if change.NewValue == "true" {
		slog.Debug("Bike immobilized, sending event...", "bike", bike.OpenDataId)
		protostruct := &bikes.BikeImmobilized{
			TimeStamp: timestamppb.Now(),
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
	}
	return nil
}

// Sends BikeAbandoned event. Needs user
func abandonedChange(bike Bike, change helper.Change) error {
	//bike abandoned
	if change.NewValue == "true" {
		user := user(change)
		slog.Debug("Bike abandoned, sending event...", "bike", bike.OpenDataId)
		protostruct := &bikes.BikeAbandoned{
			TimeStamp: timestamppb.Now(),
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
	}
	return nil
}

// Sends No event???. Needs station and user
func availableChange(bike Bike, change helper.Change) error {
	//bike available
	// if change.NewValue == "true" {
	// 	slog.Info("Bike available, sending event...", "bike", bike.OpenDataId)
	//	station, user := stationAndUser(change)

	// 	err := kc.Produce(&bikes.BikeReserved{
	// 		TimeStamp: timestamppb.Now(),
	// 		Bike: &bikes.Bike{
	// 			Bike:       bike.IntoId(),
	// 			IsElectric: bike.IsElectric.Bool,
	// 		},
	// 		Station: station.IntoId(),
	// 		User:    user.IntoId(),
	// 	})
	// 	helper.MaybeDieErr(err)
	// }
	return nil
}

// sends BikeBroughtOut and BikeStored events. Needs bike and change.
func isInStorageChange(bike Bike, change helper.Change) error {
	if change.NewValue == "false" {
		slog.Info("Bike brought out, sending event...", "bike", bike.OpenDataId)
		protostruct := &bikes.BikeDeployed{
			TimeStamp: timestamppb.Now(),
			Bike:      bike.IntoId(),
		}
		if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
			return err
		}
		return nil
	}
	if change.NewValue == "true" {
		slog.Debug("Bike stored, sending event...", "bike", bike.OpenDataId)
		protostruct := &bikes.BikeStored{
			TimeStamp: timestamppb.Now(),
			Bike:      bike.IntoId(),
		}
		if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
			return err
		}
	}
	return nil
}

// Sends BikeReserved event. Needs station and user.
func reservedChange(bike Bike, change helper.Change) error {
	//bike reserved
	if change.NewValue == "true" {
		slog.Debug("Bike reserved, sending event...", "bike", bike.OpenDataId)
		station, user := stationAndUser(change)
		protostruct := &bikes.BikeReserved{
			TimeStamp: timestamppb.Now(),
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
	}
	return nil
}

// Sends BikeDefectReported event. Needs bike, user and defect.
func defectChange(bike Bike, change helper.Change) error {
	//bike defect
	if change.NewValue == "true" {
		user, defect := userdefect(change)
		slog.Debug("Bike defect, sending event...", "bike", bike.OpenDataId)
		protostruct := &bikes.BikeDefectReported{
			TimeStamp: timestamppb.Now(),
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
	}
	return nil
}

func pickedUpChange(bike Bike, change helper.Change) error {
	//bike picked up
	slog.Debug("Bike picked up, sending event...", "bike", bike.OpenDataId)
	station, user := stationAndUser(change)
	protostruct := &bikes.BikePickedUp{
		TimeStamp: timestamppb.Now(),
		Bike:      bike.IntoId(),
		Station:   station.IntoId(),
		User:      user.IntoId(),
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}

	return nil
}

func returnedChange(bike Bike, change helper.Change) error {
	//bike returned
	slog.Debug("Bike returned, sending event...", "bike", bike.OpenDataId)
	station, user := stationAndUser(change)
	protostruct := &bikes.BikeReturned{
		TimeStamp: timestamppb.Now(),
		Bike:      bike.IntoId(),
		Station:   station.IntoId(),
		User:      user.IntoId(),
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}
	return nil
}

// checks if station and user are present and returns them
func stationAndUser(change helper.Change) (Station, User) {
	station, err := GetStationById(change.Station_id)
	helper.MaybeDieErr(err)
	user, err := GetUserById(change.User_id)
	helper.MaybeDieErr(err)
	return station, user
}

// checks if user is present and returns it
func user(change helper.Change) User {
	user, err := GetUserById(change.User_id)
	helper.MaybeDieErr(err)
	return user
}

// checks if user and defect are present and returns them
func userdefect(change helper.Change) (User, string) {
	user, err := GetUserById(change.User_id)
	helper.MaybeDieErr(err)
	if change.Defect == "" {
		helper.Die(fmt.Errorf("defect change detected, but no defect was given. defect=%v", change.Defect))
	}
	return user, change.Defect
}