package helper

import (
	"errors"
	"reflect"
)

// Returns a slice of strings with the id and names of the columns that have changed between the old and new record.
func ColumnChange(oldrecord any, record any) []string {
	changedColumns := []string{}

	if reflect.TypeOf(oldrecord) != reflect.TypeOf(record) {
		Die(errors.New("record struct types do not match"))
	}

	oldv := reflect.Indirect(reflect.ValueOf(oldrecord))
	newv := reflect.Indirect(reflect.ValueOf(record))

	changedColumns = append(changedColumns, newv.Field(1).String())

	for i := 1; i < oldv.NumField(); i++ {
		if oldv.Field(i).Interface() != newv.Field(i).Interface() {
			changedColumns = append(changedColumns, newv.Type().Field(i).Name)
		}
	}

	return changedColumns
}
