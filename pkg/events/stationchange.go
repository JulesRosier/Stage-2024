package events

import (
	"fmt"
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
	"stage2024/pkg/protogen/common"
	"stage2024/pkg/protogen/stations"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func activeChange(station database.Station, change helper.Change) {
	//station deprecated
	slog.Info("Active change detected", "station", station.OpenDataId, "newvalue", change.NewValue)
	fmt.Println()
	if change.NewValue == "false" {
		slog.Info("Station deprecated, sending event...", "station", station.OpenDataId)
		err := kc.Produce(&stations.StationDeprecated{
			TimeStamp: timestamppb.Now(),
			Station: &stations.DeprecatedStation{
				Station: &stations.StationIdentification{
					Id: station.Id,
					Location: &common.Location{
						Latitude:  station.Lat,
						Longitude: station.Lon,
					},
					Name: station.Name,
				},
				IsActive: station.IsActive.Bool,
			},
		})
		helper.MaybeDieErr(err)
	}

	//station created/activated
	if change.NewValue == "true" {
		slog.Info("Station created, sending event...", "station", station.OpenDataId)
		err := kc.Produce(&stations.StationCreated{
			TimeStamp: timestamppb.Now(),
			Station: &stations.CreatedStation{
				Station: &stations.StationIdentification{
					Id: station.Id,
					Location: &common.Location{
						Latitude:  station.Lat,
						Longitude: station.Lon,
					},
					Name: station.Name,
				},
				IsActive:    station.IsActive.Bool,
				MaxCapacity: station.MaxCapacity},
		})
		helper.MaybeDieErr(err)

	}
}

func occupationChange(station database.Station, change helper.Change) {
	//station full
	if station.Occupation == station.MaxCapacity {
		slog.Info("Station is full, sending event...", "station", station.OpenDataId)

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

	//station occupation increased
	if change.NewValue > change.OldValue {
		slog.Info("Station occupation increased, sending event...", "station", station.OpenDataId)

		err := kc.Produce(&stations.StationCapacityIncreased{
			TimeStamp: timestamppb.Now(),
			Station: &stations.StationIdentification{
				Id: station.Id,
				Location: &common.Location{
					Latitude:  station.Lat,
					Longitude: station.Lon,
				},
				Name: station.Name,
			},
			AmountIncreased:          helper.StringToInt(change.NewValue) - helper.StringToInt(change.OldValue),
			CurrentAvailableCapacity: station.Occupation,
			MaxCapacity:              station.MaxCapacity,
		})
		helper.MaybeDieErr(err)
	}

	//station occupation decreased
	if change.NewValue < change.OldValue {
		slog.Info("Station occupation decreased, sending event...", "station", station.OpenDataId)

		err := kc.Produce(&stations.StationCapacityDecreased{
			TimeStamp: timestamppb.Now(),
			Station: &stations.StationIdentification{
				Id: station.Id,
				Location: &common.Location{
					Latitude:  station.Lat,
					Longitude: station.Lon,
				},
				Name: station.Name,
			},
			AmountDecreased:          (helper.StringToInt(change.OldValue) - helper.StringToInt(change.NewValue)),
			CurrentAvailableCapacity: station.Occupation,
			MaxCapacity:              station.MaxCapacity,
		})
		helper.MaybeDieErr(err)
	}
}
