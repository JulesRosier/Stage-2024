package database

import (
	"log/slog"
	"reflect"
	"stage2024/pkg/helper"
	"stage2024/pkg/protogen/stations"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func occupationChange(station Station, change helper.Change, db *gorm.DB) error {
	newValue := helper.StringToInt(change.NewValue)
	oldValue := helper.StringToInt(change.OldValue)
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

		//add historical data and add topic name
		record := &Station{}
		db.Limit(1).Find(&record, "id = ?", change.Id)
		topic := helper.ToSnakeCase(reflect.TypeOf(protostruct).Elem().Name())
		if err := addHistoricaldata(record, topic, db, 0, protostruct.TimeStamp); err != nil {
			return err
		}
	}

	//station occupation increased
	if newValue > oldValue {
		slog.Debug("Station occupation increased, sending event...", "station", station.OpenDataId)

		now := timestamppb.Now()
		protostruct := &stations.StationOccupationIncreased{
			TimeStamp:                now,
			Station:                  station.IntoId(),
			AmountIncreased:          newValue - oldValue,
			CurrentAvailableCapacity: station.Occupation,
			MaxCapacity:              station.MaxCapacity,
		}

		if err := createOutboxRecord(now, protostruct, db); err != nil {
			return err
		}

		//add historical data, get record and add topic name
		record := &Station{}
		db.Limit(1).Find(&record, "id = ?", change.Id)
		topic := helper.ToSnakeCase(reflect.TypeOf(protostruct).Elem().Name())
		if err := addHistoricaldata(record, topic, db, protostruct.AmountIncreased, protostruct.TimeStamp); err != nil {
			return err
		}
	}

	//station occupation decreased
	if newValue < oldValue {
		slog.Debug("Station occupation decreased, sending event...", "station", station.OpenDataId)

		now := timestamppb.Now()
		protostruct := &stations.StationOccupationDecreased{
			TimeStamp:                now,
			Station:                  station.IntoId(),
			AmountDecreased:          (oldValue - newValue),
			CurrentAvailableCapacity: station.Occupation,
			MaxCapacity:              station.MaxCapacity,
		}

		if err := createOutboxRecord(now, protostruct, db); err != nil {
			return err
		}

		//add historical data, get record and add topic name
		record := &Station{}
		db.Limit(1).Find(&record, "id = ?", change.Id)
		topic := helper.ToSnakeCase(reflect.TypeOf(protostruct).Elem().Name())
		if err := addHistoricaldata(record, topic, db, protostruct.AmountDecreased, protostruct.TimeStamp); err != nil {
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
	record := &Station{}
	db.Limit(1).Find(&record, "id = ?", station.Id)
	topic := helper.ToSnakeCase(reflect.TypeOf(protostruct).Elem().Name())
	if err := addHistoricaldata(record, topic, db, 0, protostruct.TimeStamp); err != nil {
		return err
	}

	return nil
}
