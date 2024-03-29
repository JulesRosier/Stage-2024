package opendata

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/gentopendata"
	"stage2024/pkg/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func StorageTownHall(db *gorm.DB, channelCh chan []string) {
	url := "https://data.stad.gent/api/explore/v2.1/catalog/datasets/real-time-bezetting-fietsenstalling-stadskantoor-gent/records"
	model := "StorageGhent"

	slog.Info("Fetching data", "model", model)

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

	records := gentopendata.Fetch(url,
		func(b []byte) *database.Station {
			err := json.Unmarshal(b, &in)
			helper.MaybeDieErr(err)

			out := &database.Station{}
			out.Id = uuid.New().String()
			out.OpenDataId = fmt.Sprint(model+"-", in.Id)
			out.Lat = in.Geo_point_2d.Lat
			out.Lon = in.Geo_point_2d.Lon
			out.Name = in.Facilityname
			out.MaxCapacity = in.Totalplaces
			out.Occupation = in.Occupiedplaces

			return out
		},
	)
	database.UpdateStation(db, channelCh, records)

	slog.Info("Data fetched and processed, waiting...", "model", model)
}
