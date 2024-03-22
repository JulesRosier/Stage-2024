package writebolt

import (
	"context"
	"encoding/json"
	"log/slog"
	"stage2024/pkg/gentopendata"
	h "stage2024/pkg/helper"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/common"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

// updates every 5 minutes
const fetchdelay = time.Minute * 5
const url = "https://data.stad.gent/api/explore/v2.1/catalog/datasets/bolt-deelfietsen-gent/records"

type ApiData struct {
	BikeId             string `json:"bike_id"`
	CurrentRangeMeters int32  `json:"current_range_meters"`
	PricingPlanId      string `json:"pricing_plan_id"`
	VehicleTypeId      string `json:"vehicle_type_id"`
	IsReserved         int32  `json:"is_reserved"`
	IsDisabled         int32  `json:"is_disabled"`
	RentalUris         string `json:"rental_uris"`
	Loc                struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}
}

const Topic = "bolt-locations"

func WriteBolt(cl *kgo.Client, serde *sr.Serde) {
	slog.Default()

	var wg sync.WaitGroup

	for {
		allItems := gentopendata.Fetch(url,
			func(b []byte) *bikes.BoltLocation {
				var in ApiData
				err := json.Unmarshal(b, &in)
				h.MaybeDieErr(err)

				out := bikes.BoltLocation{}
				out.BikeId = in.BikeId
				out.CurrentRangeMeters = in.CurrentRangeMeters
				out.PricingPlanId = in.PricingPlanId
				out.VehicleTypeId = in.VehicleTypeId
				out.IsReserved = in.IsReserved
				out.IsDisabled = in.IsDisabled
				out.RentalUris = in.RentalUris
				out.Location = &common.Location{
					Lon: in.Loc.Lon,
					Lat: in.Loc.Lat,
				}
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
