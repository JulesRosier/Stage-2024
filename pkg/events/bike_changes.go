package events

import (
	"fmt"
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/common"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// TODO add events for bike location changes
func locationchange(bike database.Bike) {
	// slog.Info("Bike location changed, sending event...", "bike", bike.OpenDataId)
}

// Sends BikeImmobilized event. Needs bike and change.
func (ec EventClient) immobilizedChange(bike database.Bike, change helper.Change) {
	//TODO add event for bike mobilized
	//bike immobilized
	if change.NewValue == "true" {
		slog.Info("Bike immobilized, sending event...", "bike", bike.OpenDataId)
		err := ec.Kc.Produce(&bikes.BikeImmobilized{
			TimeStamp: timestamppb.Now(),
			Bike: &bikes.LocationBike{
				Bike: bike.IntoId(),
				Location: &common.Location{
					Latitude:  bike.Lat,
					Longitude: bike.Lon,
				},
			},
		})
		helper.MaybeDieErr(err)
	}
}

// Sends BikeAbandoned event. Needs user
func (ec EventClient) abandonedChange(bike database.Bike, change helper.Change) {
	//bike abandoned
	if change.NewValue == "true" {
		user := user(change)
		slog.Info("Bike abandoned, sending event...", "bike", bike.OpenDataId)
		err := ec.Kc.Produce(&bikes.BikeAbandoned{
			TimeStamp: timestamppb.Now(),
			Bike: &bikes.AbandonedBike{
				Bike: bike.IntoId(),
				Location: &common.Location{
					Latitude:  bike.Lat,
					Longitude: bike.Lon,
				},
			},
			User: user.IntoId(),
		})
		helper.MaybeDieErr(err)
	}
}

// Sends No event???. Needs station and user
func availableChange(bike database.Bike, change helper.Change) {
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
}

// sends BikeBroughtOut and BikeStored events. Needs bike and change.
func (ec EventClient) isInStorageChange(bike database.Bike, change helper.Change) {
	if change.NewValue == "false" {
		slog.Info("Bike brought out, sending event...", "bike", bike.OpenDataId)
		err := ec.Kc.Produce(&bikes.BikeDeployed{
			TimeStamp: timestamppb.Now(),
			Bike:      bike.IntoId(),
		})
		helper.MaybeDieErr(err)
	}
	if change.NewValue == "true" {
		slog.Debug("Bike stored, sending event...", "bike", bike.OpenDataId)
		err := ec.Kc.Produce(&bikes.BikeStored{
			TimeStamp: timestamppb.Now(),
			Bike:      bike.IntoId(),
		})
		helper.MaybeDieErr(err)
	}
}

// Sends BikeReserved event. Needs station and user.
func (ec EventClient) reservedChange(bike database.Bike, change helper.Change) {
	//bike reserved
	if change.NewValue == "true" {
		slog.Info("Bike reserved, sending event...", "bike", bike.OpenDataId)

		change = ec.startReservedsequence(bike, change)
		station, user := stationAndUser(change)

		err := ec.Kc.Produce(&bikes.BikeReserved{
			TimeStamp: timestamppb.Now(),
			Bike: &bikes.Bike{
				Bike:       bike.IntoId(),
				IsElectric: bike.IsElectric.Bool,
			},
			Station: station.IntoId(),
			User:    user.IntoId(),
		})
		helper.MaybeDieErr(err)
	}
}

// Sends BikeDefectReported event. Needs bike, user and defect.
func (ec EventClient) defectChange(bike database.Bike, change helper.Change) {
	//bike defect
	if change.NewValue == "true" {
		user, defect := userdefect(change)
		slog.Info("Bike defect, sending event...", "bike", bike.OpenDataId)
		err := ec.Kc.Produce(&bikes.BikeDefectReported{
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
		})

		helper.MaybeDieErr(err)
	}
}

func (ec EventClient) pickedUpChange(bike database.Bike, change helper.Change) {
	//bike picked up
	slog.Info("Bike picked up, sending event...", "bike", bike.OpenDataId)
	station, user := stationAndUser(change)
	err := ec.Kc.Produce(&bikes.BikePickedUp{
		TimeStamp: timestamppb.Now(),
		Bike:      bike.IntoId(),
		Station:   station.IntoId(),
		User:      user.IntoId(),
	})
	helper.MaybeDieErr(err)
}

func (ec EventClient) returnedChange(bike database.Bike, change helper.Change) {
	//bike returned
	slog.Info("Bike returned, sending event...", "bike", bike.OpenDataId)
	station, user := stationAndUser(change)
	err := ec.Kc.Produce(&bikes.BikeReturned{
		TimeStamp: timestamppb.Now(),
		Bike:      bike.IntoId(),
		Station:   station.IntoId(),
		User:      user.IntoId(),
	})
	helper.MaybeDieErr(err)
}

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
