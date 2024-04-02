package events

import (
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
	"stage2024/pkg/protogen/bikes"

	"google.golang.org/protobuf/types/known/timestamppb"
)

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
		slog.Info("Bike stored, sending event...", "bike", bike.OpenDataId)
		err := kc.Produce(&bikes.BikeStored{
			TimeStamp: timestamppb.Now(),
			Bike:      bike.IntoId(),
		})
		helper.MaybeDieErr(err)
	}
}

func locationchange(bike database.Bike) {
	// TODO

}

func reservedChange(bike database.Bike, change helper.Change) {
	//bike reserved
	if change.NewValue == "true" {
		err := kc.Produce(&bikes.BikeReserved{
			TimeStamp: timestamppb.Now(),
			Bike: &bikes.Bike{
				Bike:       bike.IntoId(),
				IsElectric: bike.IsElectric.Bool,
			},
			//Station: getrandomstation(), // TODO Random station
			//User:    getrandomuser(),    // TODO random user
		})
		helper.MaybeDieErr(err)
	}
}
