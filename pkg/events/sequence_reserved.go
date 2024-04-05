package events

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"stage2024/pkg/database"
	h "stage2024/pkg/helper"
)

const chanceAbandoned = 0.01
const chanceDefect = 0.01
const chanceImmobilized = 0.5

func (ec EventClient) startReservedsequence(bike database.Bike, change h.Change) h.Change {
	db := database.GetDb()

	user := database.User{}
	db.Order("random()").Find(&user)

	station := database.Station{}
	db.Order("random()").Find(&station)

	change.User_id = user.Id
	change.Station_id = station.Id

	go ec.dorestofsequence(bike, change)

	return change

}

func (ec EventClient) dorestofsequence(bike database.Bike, change h.Change) {
	slog.Info("Starting sequence RESERVED", "bike", bike.OpenDataId)
	h.RandSleep(60*5, 60)

	change.Column = "PickedUp"
	ec.Channel <- change

	changestation := ec.stationOccupationDecrease(change)

	ec.Channel <- changestation

	h.RandSleep(60*5, 60)

	if rand.Float32() < chanceAbandoned {
		change.Column = "IsAbandoned"
		ec.Channel <- change
		return
	}

	if rand.Float32() < chanceDefect {
		change.Column = "IsDefect"
		change.Defect = defects[rand.IntN(len(defects))] // get random defect
		ec.Channel <- change
		if rand.Float32() < chanceImmobilized {
			change.Column = "IsImmobilized"
			ec.Channel <- change
		}
	}

	h.RandSleep(60*5, 60)

	change.Column = "Returned"
	// change.Station_id = //TODO generate station where bike is returned
	ec.Channel <- change
}

func (ec EventClient) stationOccupationDecrease(change h.Change) h.Change {
	db := database.GetDb()
	station := database.Station{}
	//TODO  get old value from db to add into event
	db.Find(&station, "id = ?", change.Station_id)

	return h.Change{
		Table:    "Station",
		Column:   "occupation",
		Id:       change.Station_id,
		OldValue: fmt.Sprint(station.Occupation),
		NewValue: fmt.Sprint(station.Occupation - 1),
	}
}

var defects = []string{
	"Flat tire",
	"Broken chain",
	"Worn brake pads",
	"Loose spokes",
	"Faulty gear shifting",
	"Bent wheel rim",
	"Damaged pedals",
	"Cracked frame",
	"Stuck brakes",
	"Rusted components",
	"Misaligned wheels",
	"Broken saddle",
	"Malfunctioning gears",
	"Wobbly handlebars",
	"Loose headset",
	"Torn seat cover",
	"Defective bearings",
	"Faulty brakes",
	"Cracked fork",
	"Damaged crankset",
}
