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

func Baqme(channel chan helper.Change) {
	url := "https://data.stad.gent/api/explore/v2.1/catalog/datasets/baqme-locaties-vrije-deelfietsen-gent/records"
	model := "Baqme"

	slog.Info("Fetching data", "model", model)

	in := struct {
		Bike_id         string `json:"bike_id"`
		Is_reserved     int32  `json:"is_reserved"`
		Is_disabled     int32  `json:"is_disabled"`
		Vehicle_type_id string `json:"vehicle_type_id"`
		Rental_uris     string `json:"rental_uris"`
		Geopoint        struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		}
	}{}

	records := gentopendata.Fetch(url,
		func(b []byte) *database.Bike {
			err := json.Unmarshal(b, &in)
			helper.MaybeDieErr(err)

			out := &database.Bike{}
			out.OpenDataId = fmt.Sprint(model, "-", in.Bike_id)
			out.Id = uuid.New().String()
			out.BikeModel = model
			out.IsElectric = sql.NullBool{Bool: gofakeit.Bool(), Valid: true} //random bool value
			out.Lat = in.Geopoint.Lat
			out.Lon = in.Geopoint.Lon
			out.IsImmobilized = sql.NullBool{Bool: false, Valid: true} //fake
			out.IsAbandoned = sql.NullBool{Bool: false, Valid: true}   //fake
			out.IsAvailable = sql.NullBool{Bool: in.Is_disabled == 0, Valid: true}
			out.IsInStorage = sql.NullBool{Bool: false, Valid: true} //fake
			out.IsReserved = sql.NullBool{Bool: in.Is_reserved != 0, Valid: true}
			out.IsDefect = sql.NullBool{Bool: false, Valid: true} //fake
			return out
		},
	)
	database.UpdateBike(channel, records)

	slog.Debug("Data fetched and processed, waiting...", "model", model)
}
