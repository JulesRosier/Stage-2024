package stallingstadskantoordropoff

import (
	"context"
	"encoding/json"
	"fmt"
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

// updates every 5 minutes
const fetchdelay = time.Minute * 5
const url = "https://data.stad.gent/api/explore/v2.1/catalog/datasets/real-time-bezetting-fietsenstalling-stadskantoor-gent/records"

type ApiData struct {
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
}

const Topic = "bike.droppedoff"

func StadskantoorDropoff(cl *kgo.Client, serde *sr.Serde) {
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
				p := findDataByID(prevItems, item.Name)
				// diff in bikes positive number means bikes dropped off, negative means bikes picked up
				bikesdiff := p.Vacantspaces - item.Vacantspaces
				if bikesdiff == 0 {
					continue
				}
				timestamp, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
				h.MaybeDie(err, "Failed to parse time")

				item := events.BikeDropOff{
					// geen timestamp in data, dus tijd van nu
					Reported:  timestamppb.New(timestamp),
					StationId: fmt.Sprint(item.Parking, " ", item.Naam),
					Location: &common.Location{
						Lon: item.Locatie.Lon,
						Lat: item.Locatie.Lat,
					},
					Company:        "Fietsenstalling",
					Diff:           int32(bikesdiff),
					DocksAvailable: int32(item.Vacantspaces),
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
		if d.Name == id {
			return &d
		}
	}
	return nil
}
