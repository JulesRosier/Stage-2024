package events

import (
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
	"stage2024/pkg/protogen/bikes"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func availableChange(bike database.Bike, change helper.Change) {
	// TODO
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
