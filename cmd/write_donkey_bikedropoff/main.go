package writedonkeybikedropoff

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"stage2024/pkg/gentopendata"
	h "stage2024/pkg/helper"
	"stage2024/pkg/protogen/common"
	"stage2024/pkg/protogen/events"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/types/known/timestamppb"
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

const Topic = "bike.droppedoff"

func DonkeyBikeDropOff(cl *kgo.Client, serde *sr.Serde) {
	slog.Default()

	var wg sync.WaitGroup

	var prevItems []ApiData

	for {
		allItems := gentopendata.Fetch(url,
			func(b []byte) ApiData {
				var in ApiData
				err := json.Unmarshal(b, &in)
				h.MaybeDie(err, "Failed to unmarshal")
				return in
			},
		)

		if prevItems != nil {
			for _, item := range allItems {
				p := findDataByID(prevItems, item.Station_id)
				bikesDiv := item.Num_bikes_available - p.Num_bikes_available
				if bikesDiv == 0 {
					continue
				}
				// docksDiv := item.Num_docks_available - p.Num_docks_available
				unixTime, err := strconv.ParseInt(item.Last_reported, 10, 64)
				if err != nil {
					unixTime = 0
				}

				item := events.BikeDropOff{
					Reported:  timestamppb.New(time.Unix(unixTime, 0)),
					StationId: item.Station_id,
					Location: &common.Location{
						Lon: item.Geopunt.Lon,
						Lat: item.Geopunt.Lat,
					},
					Company:        "donkey-republic",
					Diff:           bikesDiv,
					DocksAvailable: item.Num_docks_available,
					BikesAvailable: item.Num_bikes_available,
				}
				ctx := context.TODO()

				h.Produce(serde, cl, &wg, &item, ctx, Topic)

			}
		}

		prevItems = allItems

		wg.Wait()
		slog.Info("Sleeping for", "duration", fetchdelay, "topic", Topic)
		time.Sleep(fetchdelay)
	}
}

func findDataByID(data []ApiData, id string) *ApiData {
	for _, d := range data {
		if d.Station_id == id {
			return &d
		}
	}
	return nil
}
