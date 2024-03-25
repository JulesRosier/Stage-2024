package opendata

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/gentopendata"
	h "stage2024/pkg/helper"
	"stage2024/pkg/protogen/common"

	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const url = "https://data.stad.gent/api/explore/v2.1/catalog/datasets/bolt-deelfietsen-gent/records"

type ApiData struct {
	BikeId             string `json:"bike_id"`
	CurrentRangeMeters int32  `json:"current_range_meters"`
	PricingPlanId      string `json:"pricing_plan_id"`
	VehicleTypeId      string `json:"vehicle_type_id"`
	IsReserved         int32  `json:"is_reserved"`
	IsDisabled         int32  `json:"is_disabled"`
	RentalUris         string `json:"rental_uris"`
	Loc                struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}
}

func Bolt(db *gorm.DB, channelCh chan string) {
	slog.Info("Fetching Bolt data")

	records := gentopendata.Fetch(url,
		func(b []byte) *database.Bike {
			var in ApiData
			err := json.Unmarshal(b, &in)
			h.MaybeDieErr(err)

			out := &database.Bike{}
			out.Id = in.BikeId
			out.BikeModel = "Bolt"
			out.IsElectric = sql.NullBool{Bool: gofakeit.Bool(), Valid: true} //random bool value
			out.Location = fmt.Sprintf("%v", &common.Location{
				Latitude:  in.Loc.Lat,
				Longitude: in.Loc.Lon,
			})
			out.IsImmobilized = sql.NullBool{Valid: false} //fake
			out.IsAbandoned = sql.NullBool{Valid: false}   //fake
			out.IsAvailable = sql.NullBool{Bool: in.IsDisabled != 0, Valid: true}
			out.IsInStorage = sql.NullBool{Valid: false} //fake
			out.IsReserved = sql.NullBool{Bool: in.IsReserved != 0, Valid: true}
			out.IsDefect = sql.NullBool{Valid: false} //fake
			return out
		},
	)
	for _, record := range records {

		oldrecord := &database.Bike{}
		result := db.Limit(1).Find(&oldrecord, "id = ?", record.Id)

		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&record)

		if result.RowsAffected == 0 {
			return
		}

		//TODO: check for changes and send to event script
		// works using channel to send to event script
		// maybe make a splice of changed columns and send all at once?
		// or send enum with changed columns?
		if oldrecord.IsAvailable.Bool != record.IsAvailable.Bool {
			channelCh <- fmt.Sprintf("Bike %s IsAvailable changed", record.Id)
		}
		if oldrecord.IsReserved.Bool != record.IsReserved.Bool {
			channelCh <- fmt.Sprintf("Bike %s IsReserved changed", record.Id)
		}
		if oldrecord.Location != record.Location {
			channelCh <- fmt.Sprintf("Bike %s Location changed, Location: %v", record.Id, record.Location)
		}

	}
	slog.Info("Done")
}
