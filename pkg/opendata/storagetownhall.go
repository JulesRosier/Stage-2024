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
func StorageTownHall(db *gorm.DB) {
	url := "https://data.stad.gent/api/explore/v2.1/catalog/datasets/real-time-bezetting-fietsenstalling-stadskantoor-gent/records"
	model := "StorageTownHall"

	slog.Info("Fetching data", "model", model)

	inStruct := struct {
		Name            string  `json:"name"`
		Parkingcapacity float32 `json:"parkingCapacity"`
		Vacantspaces    float32 `json:"vacantSpaces"`
		Naam            string  `json:"naam"`
		Parking         string  `json:"parking"`
		Occupation      int32   `json:"occupation"`
		Infotekst       string  `json:"infotekst"`
		Enginfotekst    string  `json:"enginfotekst"`
		Frinfotekst     string  `json:"frinfotekst"`
		Locatie         struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		} `json:"locatie"`
	}{}
	in := inStruct

	records := gentopendata.Fetch(url,
		func(b []byte) *database.Station {
			err := json.Unmarshal(b, &in)
			if err != nil {
				slog.Warn("Error unmarshalling data", "error", err)
				return &database.Station{}
			}

			out := &database.Station{}
			out.Id = uuid.New().String()
			out.OpenDataId = fmt.Sprint(model+"-", in.Name)
			out.Lat = in.Locatie.Lat
			out.Lon = in.Locatie.Lon
			out.Name = in.Name
			out.MaxCapacity = int32(in.Parkingcapacity)
			out.Occupation = int32(in.Parkingcapacity - in.Vacantspaces)
			out.IsActive = sql.NullBool{Bool: true, Valid: true}
			in = inStruct
			return out
		},
	)
	database.UpdateStation(records, db)

	slog.Debug("Data fetched and processed, waiting...", "model", model)
}
