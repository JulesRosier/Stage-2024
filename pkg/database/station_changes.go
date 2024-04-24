package database

import (
	"log/slog"
	"reflect"
	"stage2024/pkg/helper"
	"stage2024/pkg/protogen/stations"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func OccupationChange(station Station, change helper.Change, eventTimeStamp time.Time, db *gorm.DB) error {
	newValue, err := helper.StringToInt(change.NewValue)
	if err != nil {
		slog.Warn("Error converting string to int", "error", err)
	}
	oldValue, err := helper.StringToInt(change.OldValue)
	if err != nil {
		slog.Warn("Error converting string to int", "error", err)
	}

	timestamp := timestamppb.New(eventTimeStamp)

	//station full
	if station.Occupation == station.MaxCapacity {
		slog.Debug("Station is full, sending event...", "station", station.OpenDataId)

		protostruct := &stations.StationFull{
			TimeStamp:   timestamp,
			Station:     station.IntoId(),
			MaxCapacity: station.MaxCapacity,
		}
		if err := createOutboxRecord(timestamp, protostruct, db); err != nil {
			return err
		}

		//add historical data and add topic name
		record, err := GetStationById(change.Id, db)
		if err != nil {
			return err
		}

		topic := helper.ToSnakeCase(reflect.TypeOf(protostruct).Elem().Name())
		if err := addHistoricalData(&record, topic, db, 0, protostruct.TimeStamp); err != nil {
			return err
		}
	}

	//station occupation increased
	if newValue > oldValue {
		slog.Debug("Station occupation increased, sending event...", "station", station.OpenDataId)

		protostruct := &stations.StationOccupationIncreased{
			TimeStamp:                timestamp,
			Station:                  station.IntoId(),
			AmountIncreased:          newValue - oldValue,
			CurrentAvailableCapacity: station.Occupation,
			MaxCapacity:              station.MaxCapacity,
		}

		if err := createOutboxRecord(timestamp, protostruct, db); err != nil {
			return err
		}

		//add historical data, get record and add topic name

		record, err := GetStationById(change.Id, db)
		if err != nil {
			return err
		}
		topic := helper.ToSnakeCase(reflect.TypeOf(protostruct).Elem().Name())
		if err := addHistoricalData(&record, topic, db, protostruct.AmountIncreased, protostruct.TimeStamp); err != nil {
			return err
		}
	}

	//station occupation decreased
	if newValue < oldValue {
		slog.Debug("Station occupation decreased, sending event...", "station", station.OpenDataId)

		protostruct := &stations.StationOccupationDecreased{
			TimeStamp:                timestamp,
			Station:                  station.IntoId(),
			AmountDecreased:          (oldValue - newValue),
			CurrentAvailableCapacity: station.Occupation,
			MaxCapacity:              station.MaxCapacity,
		}

		if err := createOutboxRecord(timestamp, protostruct, db); err != nil {
			return err
		}

		//add historical data, get record and add topic name
		record, err := GetStationById(change.Id, db)
		if err != nil {
			return err
		}
		topic := helper.ToSnakeCase(reflect.TypeOf(protostruct).Elem().Name())
		if err := addHistoricalData(&record, topic, db, protostruct.AmountDecreased, protostruct.TimeStamp); err != nil {
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

func createStationEvent(station *Station, db *gorm.DB) error {
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

	//add historical data, get record and add topic name
	topic := helper.ToSnakeCase(reflect.TypeOf(protostruct).Elem().Name())
	if err := addHistoricalData(station, topic, db, 0, protostruct.TimeStamp); err != nil {
		return err
	}

	return nil
}
