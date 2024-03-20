package writebolt

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"stage2024/pkg/gentopendata"
	h "stage2024/pkg/helper"
	"stage2024/pkg/kafka"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/common"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
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

func WriteBolt(cl *kgo.Client, rcl *sr.Client) {
	slog.Default()

	topic := "bolt-locations"

	// Proto file wordt gelezen, dit moet dus ieder keer aangepast worden naar de juiste file
	file, err := os.ReadFile(filepath.Join("./proto", bikes.File_bikes_bolt_proto.Path()))
	h.MaybeDieErr(err)

	// schema ophalen
	ss := kafka.GetSchema(topic, rcl, file)

	var serde sr.Serde
	serde.Register(
		ss.ID,
		&bikes.BoltLocation{},
		sr.EncodeFn(func(a any) ([]byte, error) {
			return proto.Marshal(a.(*bikes.BoltLocation))
		}),
		sr.Index(0),
		sr.DecodeFn(func(b []byte, a any) error {
			return proto.Unmarshal(b, a.(*bikes.BoltLocation))
		}),
	)

	var wg sync.WaitGroup

	for {
		slog.Info("Producing to", "topic", topic)
		allBikes := gentopendata.Fetch(url,
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
		for _, bike := range allBikes {
			wg.Add(1)
			bikeByte, err := serde.Encode(bike)
			h.MaybeDie(err, "Encoding error")
			record := &kgo.Record{Topic: topic, Value: bikeByte}
			cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
				defer wg.Done()
				h.MaybeDie(err, "Producing error")
			})
		}
		wg.Wait()
		slog.Info("Uploaded all records", "topic", topic)
		slog.Info("Sleeping", "time", fetchdelay)
		time.Sleep(fetchdelay)
	}
}
