package opendata

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/gentopendata"
	"stage2024/pkg/helper"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func BlueBike(db *gorm.DB) {
	urls := []string{"https://data.stad.gent/api/explore/v2.1/catalog/datasets/blue-bike-deelfietsen-gent-sint-pieters-m-hendrikaplein/records",
		"https://data.stad.gent/api/explore/v2.1/catalog/datasets/blue-bike-deelfietsen-gent-dampoort/records",
		"https://data.stad.gent/api/explore/v2.1/catalog/datasets/blue-bike-deelfietsen-gent-sint-pieters-st-denijslaan/records",
		"https://data.stad.gent/api/explore/v2.1/catalog/datasets/blue-bike-deelfietsen-merelbeke-drongen-wondelgem/records"}
	model := "BlueBike"

	slog.Info("Fetching data", "model", model)

	in := struct {
		LastSeen       time.Time `json:"last_seen"`
		Id             int       `json:"id"`
		Name           string    `json:"name"`
		BikesInUse     int32     `json:"bikes_in_use"`
		BikesAvailable int32     `json:"bikes_available"`
		Longitude      string    `json:"longitude"`
		Latitude       string    `json:"latitude"`
		Geopoint       struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		} `json:"geopoint"`
		Type string `json:"type"`
	}{}

	for _, url := range urls {

		records := gentopendata.Fetch(url,
			func(b []byte) *database.Station {
				err := json.Unmarshal(b, &in)
				helper.MaybeDieErr(err)

				out := &database.Station{}

				out.Id = uuid.New().String()
				out.OpenDataId = fmt.Sprint(model+"-", in.Id)
				out.Lat = in.Geopoint.Lat
				out.Lon = in.Geopoint.Lon
				out.Name = in.Name
				out.MaxCapacity = in.BikesAvailable + in.BikesInUse
				out.Occupation = in.BikesAvailable
				out.IsActive = sql.NullBool{Bool: true, Valid: true}

				return out
			},
		)
		database.UpdateStation(records, db)
	}

	slog.Debug("Data fetched and processed, waiting...", "model", model)
}
