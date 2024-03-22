package writestallinggentdropoff

import (
	"context"
	"encoding/json"
	"log/slog"
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

// updates every minute
const fetchdelay = time.Minute * 1
const url = "https://data.stad.gent/api/explore/v2.1/catalog/datasets/real-time-bezettingen-fietsenstallingen-gent/records"

type ApiData struct {
	Time           string `json:"time"`
	Facilityname   string `json:"facilityname"`
	Id             string `json:"id"`
	Totalplaces    int32  `json:"totalplaces"`
	Freeplaces     int32  `json:"freeplaces"`
	Occupiedplaces int32  `json:"occupiedplaces"`
	Bezetting      int32  `json:"bezetting"`
	Geopoint       struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}
}

const Topic = "bike.droppedoff"

func StallingGentDropOff(cl *kgo.Client, serde *sr.Serde) {
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
				p := findDataByID(prevItems, item.Id)
				bikesdiff := p.Freeplaces - item.Freeplaces
				if bikesdiff == 0 {
					continue
				}
				timestamp, err := time.Parse(time.RFC3339, item.Time)
				h.MaybeDie(err, "Failed to parse time")

				item := events.BikeDropOff{
					Reported:  timestamppb.New(timestamp),
					StationId: item.Id,
					Location: &common.Location{
						Lon: item.Geopoint.Lon,
						Lat: item.Geopoint.Lat,
					},
					Company:        "Fietsenstalling",
					Diff:           bikesdiff,
					DocksAvailable: item.Freeplaces,
					BikesAvailable: -1,
				}
				ctx := context.TODO()

				wg.Add(1)

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
		if d.Id == id {
			return &d
		}
	}
	return nil
}
