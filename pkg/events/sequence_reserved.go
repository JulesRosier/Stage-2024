package events

import (
	"database/sql"
	"log/slog"
	"math/rand/v2"
	"stage2024/pkg/database"
	h "stage2024/pkg/helper"
)

const chanceAbandoned = 10.1
const chanceDefect = 10.5
const chanceImmobilized = 15.0
const chanceInStorage = 0.5

func (ec EventClient) startReservedsequence(bike database.Bike, change h.Change) h.Change {
	db := database.GetDb()

	user := database.User{}
	db.Order("random()").First(&user)

	//get station that is not empty and active
	station := database.Station{}
	db.Where("occupation > ? AND is_active = ?", 0, true).Order("random()").First(&station)

	change.User_id = user.Id
	change.Station_id = station.Id

	go ec.dorestofsequence(bike, change)

	return change

}

func (ec EventClient) dorestofsequence(bike database.Bike, change h.Change) {
	slog.Info("Starting sequence RESERVED", "bike", bike.OpenDataId)

	db := database.GetDb()

	h.RandSleep(60*5, 60)

	change.Column = "PickedUp"
	ec.Channel <- change

	// Set bike to unavailable, set is_reserved to false
	bike.IsAvailable = sql.NullBool{Bool: false, Valid: true}
	bike.IsReserved = sql.NullBool{Bool: false, Valid: true}
	database.UpdateBike(ec.Channel, []*database.Bike{&bike})

	// Update station occupation
	pickupStation, err := database.GetStationById(change.Station_id)
	h.MaybeDieErr(err)
	pickupStation.Occupation--
	database.UpdateStation(ec.Channel, []*database.Station{&pickupStation})

	h.RandSleep(60*5, 60)

	if rand.Float32() < chanceAbandoned {
		bike.IsAbandoned = sql.NullBool{Bool: true, Valid: true}
		database.UpdateBikeNoNotify(&bike)
		change.Column = "IsAbandoned"
		ec.Channel <- change
		return
	}

	if rand.Float32() < chanceDefect {
		change.Column = "IsDefect"
		bike.IsDefect = sql.NullBool{Bool: true, Valid: true}
		change.Defect = defects[rand.IntN(len(defects))] // get random defect
		ec.Channel <- change

		if rand.Float32() < chanceImmobilized {
			bike.IsImmobilized = sql.NullBool{Bool: true, Valid: true}
			database.UpdateBikeNoNotify(&bike)
			change.Column = "IsImmobilized"
			ec.Channel <- change

			if rand.Float32() < chanceInStorage {
				bike.IsInStorage = sql.NullBool{Bool: true, Valid: true}
				database.UpdateBikeNoNotify(&bike)
				change.Column = "IsInStorage"
				change.NewValue = "true"
				ec.Channel <- change
			}
		}
	}

	h.RandSleep(60*5, 60)

	// Get station that is not full
	change.Column = "Returned"
	returnstation := database.Station{}
	db.Where("occupation < max_capacity AND is_active = ?", true).Order("random()").First(&returnstation)
	change.Station_id = returnstation.Id
	ec.Channel <- change

	// Update station occupation
	returnstation.Occupation++
	database.UpdateStation(ec.Channel, []*database.Station{&returnstation})

	// Set bike to available
	bike.IsAvailable = sql.NullBool{Bool: true, Valid: true}
	database.UpdateBike(ec.Channel, []*database.Bike{&bike})

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
