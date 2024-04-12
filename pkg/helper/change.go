package helper

import "time"

type Change struct {
	Table      string
	Column     string
	Id         string
	OpenDataId string
	OldValue   string
	NewValue   string

	StationId string
	UserId    string
	Defect    string
	EventTime time.Time
}
