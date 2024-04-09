package helper

import "time"

type Change struct {
	Table      string
	Column     string
	Id         string
	OpenDataId string
	OldValue   string
	NewValue   string
	//onderstaande drie zijn standaard leeg, enkel bij gegenereerde data ingevuld
	StationId string
	UserId    string
	Defect    string
	EventTime time.Time
}
