package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"stage2024/pkg/gentopendata"
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

func main() {
	slog.Default()

	seed := flag.String("seedbroker", "localhost:19092", "brokers port to talk to")
	registry := flag.String("registry", "localhost:18081", "schema registry port to talk to")
	topic := flag.String("topic", "bolt-test", "topic to produce to and consume from")

	slog.Info("Starting kafka client", "seedbrokers", *seed)
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(*seed),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		panic(err)
	}
	defer cl.Close()

	slog.Info("Starting schema registry client", "host", *registry)
	rcl, err := sr.NewClient(sr.URLs(*registry))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	p := common.File_common_location_proto.Path()
	file, err := os.ReadFile(filepath.Join("./proto", p))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	ssLocation, err := rcl.CreateSchema(context.Background(), p,
		sr.Schema{
			Schema: string(file),
			Type:   sr.TypeProtobuf,
		},
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("created or reusing schema",
		"subject", ssLocation.Subject,
		"version", ssLocation.Version,
		"id", ssLocation.ID,
	)

	sub := *topic + "-value"
	file, err = os.ReadFile(filepath.Join("./proto", bikes.File_bikes_bolt_proto.Path()))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	ss, err := rcl.CreateSchema(context.Background(), sub,
		sr.Schema{
			Schema: string(file),
			Type:   sr.TypeProtobuf,
			References: []sr.SchemaReference{
				{
					Name:    p,
					Subject: ssLocation.Subject,
					Version: ssLocation.Version,
				},
			},
		})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("created or reusing schema",
		"subject", ss.Subject,
		"version", ss.Version,
		"id", ss.ID,
	)

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

	slog.Info("Producing records")
	for {
		allBikes := gentopendata.Fetch[*bikes.BoltLocation](url,
			func(b []byte) *bikes.BoltLocation {
				var in ApiData
				json.Unmarshal(b, &in)

				out := bikes.BoltLocation{}
				out.BikeId = in.BikeId
				out.CurrentRangeMeters = in.CurrentRangeMeters
				out.PricingPlanId = in.PricingPlanId
				out.VehicleTypeId = in.VehicleTypeId
				out.IsReserved = in.IsReserved
				out.IsDisabled = in.IsDisabled
				out.RentalUris = in.RentalUris
				out.Loc = &common.Location{
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
			if err != nil {
				slog.Error("Bike encoding", "error", err)
				return
			}
			record := &kgo.Record{Topic: "bolt-test", Value: bikeByte}
			cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
				defer wg.Done()
				if err != nil {
					slog.Error("record had a produce error",
						"error", err,
					)
				}

			})
		}

		wg.Wait()
		slog.Info("uploaded all records")
		slog.Info("Sleeping", "time", fetchdelay)
		time.Sleep(fetchdelay)
	}
}
