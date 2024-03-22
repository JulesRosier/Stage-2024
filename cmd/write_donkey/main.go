package writedonkey

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"stage2024/pkg/gentopendata"
	h "stage2024/pkg/helper"
	"stage2024/pkg/protogen/common"
	"stage2024/pkg/protogen/occupations"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

// updates every 10 minutes
const fetchdelay = time.Minute * 10
const url = "https://data.stad.gent/api/explore/v2.1/catalog/datasets/donkey-republic-beschikbaarheid-deelfietsen-per-station/records"

type ApiData struct {
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
}

const Topic = "donkey-locations"

func WriteDonkey(cl *kgo.Client, serde *sr.Serde) {
	slog.Default()

	var wg sync.WaitGroup

	for {
		allItems := gentopendata.Fetch(url,
			func(b []byte) *occupations.DonkeyLocation {
				var in ApiData
				err := json.Unmarshal(b, &in)
				h.MaybeDie(err, "Failed to unmarshal")
				out := occupations.DonkeyLocation{}
				out.StationId = in.Station_id
				out.NumBikesAvailable = in.Num_bikes_available
				out.NumDocksAvailable = in.Num_docks_available
				out.IsRenting = in.Is_renting
				out.IsInstalled = in.Is_installed
				out.IsReturning = in.Is_returning
				out.LastReported = in.Last_reported
				out.Location = &common.Location{
					Lon: in.Geopunt.Lon,
					Lat: in.Geopunt.Lat,
				}
				out.Name = in.Name

				return &out
			},
		)
		ctx := context.Background()
		for _, item := range allItems {
			wg.Add(1)
			h.Produce(serde, cl, &wg, item, ctx, Topic)
		}
		wg.Wait()
		slog.Info("Sleeping for", "duration", fetchdelay, "topic", Topic)
		time.Sleep(fetchdelay)
	}
}
