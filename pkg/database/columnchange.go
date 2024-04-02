package database

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"stage2024/pkg/helper"
)

func ColumnChange(oldrecord any, record any, channel chan helper.Change) {

	change := helper.Change{}

	if reflect.TypeOf(oldrecord) != reflect.TypeOf(record) {
		helper.Die(errors.New("record struct types do not match"))
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

			if change.Column != "Lat" && change.Column != "Lon" {
				change.OldValue = fmt.Sprintf("%v", oldv.Field(i))
				change.NewValue = fmt.Sprintf("%v", newv.Field(i))
			}
			// Special case for geopoints
			if change.Column == "Lat" || change.Column == "Lon" && !latchange {
				change.OldValue = fmt.Sprintf("(%v, %v)", oldv.Field(i), oldv.Field(i+1))
				change.NewValue = fmt.Sprintf("(%v, %v)", newv.Field(i), newv.Field(i+1))
				latchange = true
			}

			slog.Debug("change", "", change)

			channel <- change
		}
	}
}
