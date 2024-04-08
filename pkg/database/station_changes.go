package database

import (
	"log/slog"
	"stage2024/pkg/helper"
	"stage2024/pkg/protogen/stations"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func occupationChange(station Station, change helper.Change, db *gorm.DB) error {
	//station full
	if station.Occupation == station.MaxCapacity {
		slog.Debug("Station is full, sending event...", "station", station.OpenDataId)

		now := timestamppb.Now()
		protostruct := &stations.StationFull{
			TimeStamp:   now,
			Station:     station.IntoId(),
			MaxCapacity: station.MaxCapacity,
		}
		if err := createOutboxRecord(now, protostruct, db); err != nil {
			return err
		}

		return nil
	}

	//station occupation increased
	if change.NewValue > change.OldValue {
		slog.Debug("Station occupation increased, sending event...", "station", station.OpenDataId)

		now := timestamppb.Now()
		protostruct := &stations.StationOccupationIncreased{
			TimeStamp:                now,
			Station:                  station.IntoId(),
			AmountIncreased:          helper.StringToInt(change.NewValue) - helper.StringToInt(change.OldValue),
			CurrentAvailableCapacity: station.Occupation,
			MaxCapacity:              station.MaxCapacity,
		}

		if err := createOutboxRecord(now, protostruct, db); err != nil {
			return err
		}
	}

	//station occupation decreased
	if change.NewValue < change.OldValue {
		slog.Debug("Station occupation decreased, sending event...", "station", station.OpenDataId)

		now := timestamppb.Now()
		protostruct := &stations.StationOccupationDecreased{
			TimeStamp:                now,
			Station:                  station.IntoId(),
			AmountDecreased:          (helper.StringToInt(change.OldValue) - helper.StringToInt(change.NewValue)),
			CurrentAvailableCapacity: station.Occupation,
			MaxCapacity:              station.MaxCapacity,
		}

		if err := createOutboxRecord(now, protostruct, db); err != nil {
			return err
		}
	}
	return nil
}

func activeChange(station Station, change helper.Change, db *gorm.DB) error {
	//station deprecated
	if change.NewValue == "false" {
		slog.Debug("Station deprecated, sending event...", "station", station.OpenDataId)
		now := timestamppb.Now()
		protostruct := &stations.StationDeprecated{
			TimeStamp: timestamppb.Now(),
			Station: &stations.DeprecatedStation{
				Station:  station.IntoId(),
				IsActive: station.IsActive.Bool,
			},
		}

		if err := createOutboxRecord(now, protostruct, db); err != nil {
			return err
		}
	}

	//station created/activated
	if change.NewValue == "true" {
		slog.Debug("Station created, sending event...", "station", station.OpenDataId)
		now := timestamppb.Now()
		protostruct := &stations.StationCreated{
			TimeStamp: timestamppb.Now(),
			Station: &stations.CreatedStation{
				Station:     station.IntoId(),
				IsActive:    station.IsActive.Bool,
				MaxCapacity: station.MaxCapacity},
		}

		if err := createOutboxRecord(now, protostruct, db); err != nil {
			return err
		}
	}
	return nil
}

func created(station Station, db *gorm.DB) error {
	slog.Debug("Station created, sending event...", "station", station.OpenDataId)
	now := timestamppb.Now()
	protostruct := &stations.StationCreated{
		TimeStamp: timestamppb.Now(),
		Station: &stations.CreatedStation{
			Station:     station.IntoId(),
			IsActive:    station.IsActive.Bool,
			MaxCapacity: station.MaxCapacity},
	}

	if err := createOutboxRecord(now, protostruct, db); err != nil {
		return err
	}
	return nil
}
