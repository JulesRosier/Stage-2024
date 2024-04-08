package database

import (
	"fmt"
	"log/slog"
	"stage2024/pkg/helper"
	"stage2024/pkg/protogen/stations"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func occupationChange(station Station, change helper.Change, db *gorm.DB) error {
	//station full
	if station.Occupation == station.MaxCapacity {
		slog.Info("Station is full, sending event...", "station", station.OpenDataId)

		now := timestamppb.Now()
		protostruct := &stations.StationFull{
			TimeStamp:   now,
			Station:     station.IntoId(),
			MaxCapacity: station.MaxCapacity,
		}
		payload, err := proto.Marshal(protostruct)
		if err != nil {
			return err
		}
		if err := createOutboxRecord(now, protostruct, payload, db); err != nil {
			return err
		}

		return err
	}

	//station occupation increased
	if change.NewValue > change.OldValue {
		slog.Info("Station occupation increased, sending event...", "station", station.OpenDataId)

		now := timestamppb.Now()
		protostruct := &stations.StationOccupationIncreased{
			TimeStamp:                now,
			Station:                  station.IntoId(),
			AmountIncreased:          helper.StringToInt(change.NewValue) - helper.StringToInt(change.OldValue),
			CurrentAvailableCapacity: station.Occupation,
			MaxCapacity:              station.MaxCapacity,
		}
		payload, err := proto.Marshal(protostruct)

		if err != nil {
			return err
		}
		if err := createOutboxRecord(now, protostruct, payload, db); err != nil {
			return err
		}
	}

	//station occupation decreased
	if change.NewValue < change.OldValue {
		slog.Info("Station occupation decreased, sending event...", "station", station.OpenDataId)

		now := timestamppb.Now()
		protostruct := &stations.StationOccupationDecreased{
			TimeStamp:                now,
			Station:                  station.IntoId(),
			AmountDecreased:          (helper.StringToInt(change.OldValue) - helper.StringToInt(change.NewValue)),
			CurrentAvailableCapacity: station.Occupation,
			MaxCapacity:              station.MaxCapacity,
		}
		payload, err := proto.Marshal(protostruct)
		if err != nil {
			return err
		}
		if err := createOutboxRecord(now, protostruct, payload, db); err != nil {
			return err
		}
	}
	return nil
}

func activeChange(station Station, change helper.Change, db *gorm.DB) error {
	//station deprecated
	fmt.Println()
	if change.NewValue == "false" {
		slog.Info("Station deprecated, sending event...", "station", station.OpenDataId)
		now := timestamppb.Now()
		protostruct := &stations.StationDeprecated{
			TimeStamp: timestamppb.Now(),
			Station: &stations.DeprecatedStation{
				Station:  station.IntoId(),
				IsActive: station.IsActive.Bool,
			},
		}
		payload, err := proto.Marshal(protostruct)
		if err != nil {
			return err
		}
		if err := createOutboxRecord(now, protostruct, payload, db); err != nil {
			return err
		}
	}

	//station created/activated
	if change.NewValue == "true" {
		slog.Info("Station created, sending event...", "station", station.OpenDataId)
		now := timestamppb.Now()
		protostruct := &stations.StationCreated{
			TimeStamp: timestamppb.Now(),
			Station: &stations.CreatedStation{
				Station:     station.IntoId(),
				IsActive:    station.IsActive.Bool,
				MaxCapacity: station.MaxCapacity},
		}

		payload, err := proto.Marshal(protostruct)
		if err != nil {
			return err
		}
		if err := createOutboxRecord(now, protostruct, payload, db); err != nil {
			return err
		}
	}
	return nil
}

func created(station Station, db *gorm.DB) error {
	slog.Info("Station created, sending event...", "station", station.OpenDataId)
	now := timestamppb.Now()
	protostruct := &stations.StationCreated{
		TimeStamp: timestamppb.Now(),
		Station: &stations.CreatedStation{
			Station:     station.IntoId(),
			IsActive:    station.IsActive.Bool,
			MaxCapacity: station.MaxCapacity},
	}
	payload, err := proto.Marshal(protostruct)
	if err != nil {
		return err
	}
	if err := createOutboxRecord(now, protostruct, payload, db); err != nil {
		return err
	}
	return err
}
