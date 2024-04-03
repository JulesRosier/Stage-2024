package events

import (
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
func immobilizedChange(bike database.Bike, change helper.Change) {
	//TODO add event for bike mobilized
	//bike immobilized
	if change.NewValue == "true" {
		slog.Info("Bike immobilized, sending event...", "bike", bike.OpenDataId)
		err := kc.Produce(&bikes.BikeImmobilized{
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
func abandonedChange(bike database.Bike, user database.User, change helper.Change) {
	//bike abandoned
	if change.NewValue == "true" {
		slog.Info("Bike abandoned, sending event...", "bike", bike.OpenDataId)
		err := kc.Produce(&bikes.BikeAbandoned{
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
func availableChange(bike database.Bike, station database.Station, user database.User, change helper.Change) {
	//bike available
	// if change.NewValue == "true" {
	// 	slog.Info("Bike available, sending event...", "bike", bike.OpenDataId)
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
func isInStorageChange(bike database.Bike, change helper.Change) {
	if change.NewValue == "false" {
		slog.Info("Bike brought out, sending event...", "bike", bike.OpenDataId)
		err := kc.Produce(&bikes.BikeBroughtOut{
			TimeStamp: timestamppb.Now(),
			Bike:      bike.IntoId(),
		})
		helper.MaybeDieErr(err)
	}
	if change.NewValue == "true" {
		slog.Debug("Bike stored, sending event...", "bike", bike.OpenDataId)
		err := kc.Produce(&bikes.BikeStored{
			TimeStamp: timestamppb.Now(),
			Bike:      bike.IntoId(),
		})
		helper.MaybeDieErr(err)
	}
}

// Sends BikeReserved event. Needs station and user.
func reservedChange(bike database.Bike, station database.Station, user database.User, change helper.Change) {
	//bike reserved
	if change.NewValue == "true" {
		slog.Info("Bike reserved, sending event...", "bike", bike.OpenDataId)
		err := kc.Produce(&bikes.BikeReserved{
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
func defectChange(bike database.Bike, user database.User, change helper.Change, defect string) {
	//bike defect
	if change.NewValue == "true" {
		slog.Info("Bike defect, sending event...", "bike", bike.OpenDataId)
		err := kc.Produce(&bikes.BikeDefectReported{
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