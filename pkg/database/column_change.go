package database

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"stage2024/pkg/helper"

	"gorm.io/gorm"
)

// Compares two records and sends changes to a channel
func ColumnChange(oldrecord any, record any, db *gorm.DB) error {

	change := helper.Change{}

	if reflect.TypeOf(oldrecord) != reflect.TypeOf(record) {
		return errors.New("record struct types do not match")
	}

	oldv := reflect.Indirect(reflect.ValueOf(oldrecord))
	newv := reflect.Indirect(reflect.ValueOf(record))
	change.Table = newv.Type().Name()
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

			slog.Debug("change", "", change)

			if err := ChangeDetected(change, db); err != nil {
				return err
			}

		}
	}
	return nil
}
