package opendata

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/gentopendata"
	"stage2024/pkg/helper"
	"stage2024/pkg/protogen/common"

	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/gorm"
)

const url = "https://data.stad.gent/api/explore/v2.1/catalog/datasets/bolt-deelfietsen-gent/records"
const Model = "Bolt"

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

func Bolt(db *gorm.DB, channelCh chan []string) {
	slog.Info("Fetching data", "model", Model)

	records := gentopendata.Fetch(url,
		func(b []byte) *database.Bike {
			faker := gofakeit.New(42) // seed to get same random bool values each time
			var in ApiData
			err := json.Unmarshal(b, &in)
			helper.MaybeDieErr(err)

			out := &database.Bike{}
			out.Id = in.BikeId
			out.BikeModel = Model
			out.IsElectric = sql.NullBool{Bool: faker.Bool(), Valid: true} //random bool value
			out.Location = fmt.Sprintf("%v", &common.Location{
				Latitude:  in.Loc.Lat,
				Longitude: in.Loc.Lon,
			})
			out.IsImmobilized = sql.NullBool{Valid: false} //fake
			out.IsAbandoned = sql.NullBool{Valid: false}   //fake
			out.IsAvailable = sql.NullBool{Bool: in.IsDisabled == 0, Valid: true}
			out.IsInStorage = sql.NullBool{Valid: false} //fake
			out.IsReserved = sql.NullBool{Bool: in.IsReserved != 0, Valid: true}
			out.IsDefect = sql.NullBool{Valid: false} //fake
			return out
		},
	)
	database.UpdateRecords(db, channelCh, records)

	slog.Info("Data fetched and processed, waiting...", "model", Model)
}
