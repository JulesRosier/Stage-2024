package opendata

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/gentopendata"
	"stage2024/pkg/helper"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Bolt(db *gorm.DB, channelCh chan []string) {
	url := "https://data.stad.gent/api/explore/v2.1/catalog/datasets/bolt-deelfietsen-gent/records"
	model := "Bolt"

	slog.Info("Fetching data", "model", model)

	in := struct {
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
	}{}

	records := gentopendata.Fetch(url,
		func(b []byte) *database.Bike {
			faker := gofakeit.New(42) // seed to get same random bool values each time
			err := json.Unmarshal(b, &in)
			helper.MaybeDieErr(err)

			out := &database.Bike{}
			out.OpenDataId = fmt.Sprint(model, "-", in.BikeId)
			out.Id = uuid.New().String()
			out.BikeModel = model
			out.IsElectric = sql.NullBool{Bool: faker.Bool(), Valid: true} //random bool value
			out.Lat = in.Loc.Lat
			out.Lon = in.Loc.Lon
			out.IsImmobilized = sql.NullBool{Valid: false} //fake
			out.IsAbandoned = sql.NullBool{Valid: false}   //fake
			out.IsAvailable = sql.NullBool{Bool: in.IsDisabled == 0, Valid: true}
			out.IsInStorage = sql.NullBool{Valid: false} //fake
			out.IsReserved = sql.NullBool{Bool: in.IsReserved != 0, Valid: true}
			out.IsDefect = sql.NullBool{Valid: false} //fake
			return out
		},
	)
	database.UpdateBike(db, channelCh, records)

	slog.Info("Data fetched and processed, waiting...", "model", model)
}
