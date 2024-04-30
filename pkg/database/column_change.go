package database

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"stage2024/pkg/helper"
	"time"

	"gorm.io/gorm"
)

// Compares two records and sends changes to the outbox
func ColumnChange(oldrecord any, record any, db *gorm.DB, change helper.Change) error {

	if reflect.TypeOf(oldrecord) != reflect.TypeOf(record) {
		return errors.New("record struct types do not match")
	}

	oldv := reflect.Indirect(reflect.ValueOf(oldrecord))
	newv := reflect.Indirect(reflect.ValueOf(record))
	change.Id = newv.Field(1).String()
	change.OpenDataId = newv.Field(2).String()

	latchange := false

	for i := 1; i < oldv.NumField(); i++ {
		if oldv.Field(i).Interface() != newv.Field(i).Interface() {
			change.Column = newv.Type().Field(i).Name

			// Special case for geopoints
			if change.Column == "Lat" || change.Column == "Lon" && !latchange {
				change.OldValue = fmt.Sprintf("(%v, %v)", oldv.Field(i), oldv.Field(i+1))
				change.NewValue = fmt.Sprintf("(%v, %v)", newv.Field(i), newv.Field(i+1))
				latchange = true
			} else if newv.Field(i).Type().String() == "sql.NullBool" { // Special case for sql.NullBool because string conversion is different
				change.OldValue = fmt.Sprintf("%v", oldv.Field(i).Field(0))
				change.NewValue = fmt.Sprintf("%v", newv.Field(i).Field(0))
			} else {
				change.OldValue = fmt.Sprintf("%v", oldv.Field(i))
				change.NewValue = fmt.Sprintf("%v", newv.Field(i))
			}

			slog.Debug("change", "Column", change.Column, "id", change.Id)

			if err := changeDetected(change, db); err != nil {
				return err
			}

		}
	}
	return nil
}

// Selects right events to send based on change for stations
func changeDetected(change helper.Change, db *gorm.DB) error {
	station, err := GetStationById(change.Id, db)
	if err != nil {
		slog.Warn("Station not found", "station", change.Id, "openDataId", change.OpenDataId, "error", err)
		return err
	}

	switch change.Column {
	case "Occupation":
		err := OccupationChange(station, change, time.Now(), db)
		return err
	case "IsActive":
		err := activeChange(station, change, db)
		return err
	case "MaxCapacity":
		//slog.Info("MaxCapacity change detected", "station", station.OpenDataId, "old", change.OldValue, "new", change.NewValue)
	}

	return err
}
