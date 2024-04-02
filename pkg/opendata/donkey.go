package opendata

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/gentopendata"
	"stage2024/pkg/helper"

	"github.com/google/uuid"
)

func Donkey(channel chan helper.Change) {
	url := "https://data.stad.gent/api/explore/v2.1/catalog/datasets/donkey-republic-beschikbaarheid-deelfietsen-per-station/records"
	model := "Donkey"

	slog.Info("Fetching data", "model", model)

	in := struct {
		Station_id          string `json:"station_id"`
		Num_bikes_available int32  `json:"num_bikes_available"`
		Num_docks_available int32  `json:"num_docks_available"`
		Is_renting          int32  `json:"is_renting"`
		Is_installed        int32  `json:"is_installed"`
		Is_returning        int32  `json:"is_returning"`
		Last_reported       string `json:"last_reported"`
		Geopunt             struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		}
		Name string `json:"name"`
		Type string `json:"type"`
	}{}

	records := gentopendata.Fetch(url,
		func(b []byte) *database.Station {
			err := json.Unmarshal(b, &in)
			helper.MaybeDieErr(err)

			out := &database.Station{}

			out.Id = uuid.New().String()
			out.OpenDataId = fmt.Sprint(model+"-", in.Station_id)
			out.Lat = in.Geopunt.Lat
			out.Lon = in.Geopunt.Lon
			out.Name = in.Name
			out.MaxCapacity = in.Num_bikes_available + in.Num_docks_available
			out.Occupation = in.Num_bikes_available
			out.IsActive = sql.NullBool{Bool: in.Is_renting != 0, Valid: true} //sql.NullBool{Bool: gofakeit.Bool(), Valid: true}

			return out
		},
	)
	database.UpdateStation(channel, records)

	// stationfull := &database.Station{
	// 	Id:          "00000000-0000-0000-0000-000000000000",
	// 	OpenDataId:  "station-full",
	// 	Lat:         0,
	// 	Lon:         0,
	// 	Name:        "Station Full",
	// 	MaxCapacity: 5,
	// 	Occupation:  4,
	// 	IsActive:    sql.NullBool{Bool: true, Valid: true},
	// }
	// database.UpdateStation(channel, []*database.Station{stationfull})

	slog.Info("Data fetched and processed, waiting...", "model", model)
}
