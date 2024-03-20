package writebaqme

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"stage2024/pkg/gentopendata"
	h "stage2024/pkg/helper"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/common"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

// updates every 10 minutes
const fetchdelay = time.Minute * 10
const url = "https://data.stad.gent/api/explore/v2.1/catalog/datasets/baqme-locaties-vrije-deelfietsen-gent/records"

type ApiData struct {
	Bike_id         string `json:"bike_id"`
	Is_reserved     int32  `json:"is_reserved"`
	Is_disabled     int32  `json:"is_disabled"`
	Vehicle_type_id string `json:"vehicle_type"`
	Rental_uris     string `json:"rental_uris"`
	Geopoint        struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}
}

const Topic = "baqme-locations"

// Writes data to the topic
func WriteBaqme(cl *kgo.Client, serde *sr.Serde) {
	slog.Default()

	var wg sync.WaitGroup

	for {
		allItems := gentopendata.Fetch(url,
			func(b []byte) *bikes.BaqmeLocation {
				var in ApiData
				err := json.Unmarshal(b, &in)
				h.MaybeDieErr(err)

				out := bikes.BaqmeLocation{}
				out.BikeId = in.Bike_id
				out.IsReserved = in.Is_reserved
				out.IsDisabled = in.Is_disabled
				out.VehicleTypeId = in.Vehicle_type_id
				out.RentalUris = in.Rental_uris
				out.Location = &common.Location{
					Lon: in.Geopoint.Lon,
					Lat: in.Geopoint.Lat,
				}
				return &out
			},
		)
		ctx := context.Background()
		for _, item := range allItems {
			wg.Add(1)
			itemByte, err := serde.Encode(item)
			h.MaybeDie(err, "Encoding")

			record := &kgo.Record{Topic: Topic, Value: itemByte}
			cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
				defer wg.Done()
				h.MaybeDie(err, "Producing")
			})
		}
		wg.Wait()
		slog.Info("Sleeping for", "duration", fetchdelay, "topic", Topic)
		time.Sleep(fetchdelay)
	}
}
