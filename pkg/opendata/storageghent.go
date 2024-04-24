package opendata

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/gentopendata"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Fetches data from the OpenData API, sends changes to oltp database and sends events to outbox
func StorageGhent(db *gorm.DB) {
	url := "https://data.stad.gent/api/explore/v2.1/catalog/datasets/real-time-bezettingen-fietsenstallingen-gent/records"
	model := "StorageGhent"

	slog.Info("Fetching data", "model", model)

	records := gentopendata.Fetch(url,
		func(b []byte) *database.Station {

			in := struct {
				Time           string `json:"time"`
				Facilityname   string `json:"facilityname"`
				Id             string `json:"id"`
				Totalplaces    int32  `json:"totalplaces"`
				Freeplaces     int32  `json:"freeplaces"`
				Occupiedplaces int32  `json:"occupiedplaces"`
				Bezetting      int32  `json:"bezetting"`
				Geo_point_2d   struct {
					Lon float64 `json:"lon"`
					Lat float64 `json:"lat"`
				}
			}{}

			err := json.Unmarshal(b, &in)
			if err != nil {
				slog.Warn("Error unmarshalling data", "error", err)
				return &database.Station{}
			}

			out := &database.Station{}
			out.Id = uuid.New().String()
			out.OpenDataId = fmt.Sprint(model+"-", in.Id)
			out.Lat = in.Geo_point_2d.Lat
			out.Lon = in.Geo_point_2d.Lon
			out.Name = in.Facilityname
			out.MaxCapacity = in.Totalplaces
			out.Occupation = in.Occupiedplaces
			out.IsActive = sql.NullBool{Bool: true, Valid: true}
			if in.Facilityname == "Graslei" {
				out.Lat = 51.053970
				out.Lon = 3.721184
			}

			return out
		},
	)
	database.UpdateStation(records, db)

	slog.Debug("Data fetched and processed, waiting...", "model", model)
}
