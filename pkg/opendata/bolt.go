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
)

func Bolt(channel chan helper.Change) {
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
			err := json.Unmarshal(b, &in)
			helper.MaybeDieErr(err)

			out := &database.Bike{}
			out.OpenDataId = fmt.Sprint(model, "-", in.BikeId)
			out.Id = uuid.New().String()
			out.BikeModel = model
			out.IsElectric = sql.NullBool{Bool: gofakeit.Bool(), Valid: true} //random bool value
			out.Lat = in.Loc.Lat
			out.Lon = in.Loc.Lon
			out.IsImmobilized = sql.NullBool{Bool: false, Valid: true} //fake
			out.IsAbandoned = sql.NullBool{Bool: false, Valid: true}   //fake
			out.IsAvailable = sql.NullBool{Bool: in.IsDisabled == 0, Valid: true}
			out.IsInStorage = sql.NullBool{Bool: false, Valid: true} //fake
			out.IsReserved = sql.NullBool{Bool: in.IsReserved != 0, Valid: true}
			out.IsDefect = sql.NullBool{Bool: false, Valid: true} //fake
			return out
		},
	)
	database.UpdateBike(channel, records)

	biketest := &database.Bike{
		Id:            "00000000-test-test-test-000000000000",
		OpenDataId:    "bike-test",
		BikeModel:     "test",
		Lat:           0,
		Lon:           0,
		IsElectric:    sql.NullBool{Bool: false, Valid: true},
		IsImmobilized: sql.NullBool{Bool: false, Valid: true},
		IsAbandoned:   sql.NullBool{Bool: false, Valid: true},
		IsAvailable:   sql.NullBool{Bool: false, Valid: true},
		IsInStorage:   sql.NullBool{Bool: false, Valid: true},
		IsReserved:    sql.NullBool{Bool: true, Valid: true},
		IsDefect:      sql.NullBool{Bool: false, Valid: true},
	}
	database.UpdateBike(channel, []*database.Bike{biketest})

	slog.Debug("Data fetched and processed, waiting...", "model", model)
}
