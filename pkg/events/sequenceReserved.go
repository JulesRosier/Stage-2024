package events

import (
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
)

func startReservedsequence(bike database.Bike, change helper.Change) helper.Change {
	db := database.GetDb()

	user := database.User{}
	db.Order("random()").Find(&user)

	station := database.Station{}
	db.Order("random()").Find(&station)

	change.User_id = user.Id
	change.Station_id = station.Id

	go dorestofsequence(bike, change)

	return change

}

func dorestofsequence(bike database.Bike, change helper.Change) {
	slog.Info("Starting sequence RESERVED", "bike", bike.OpenDataId)
	//TODO
}
