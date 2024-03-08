package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"stage2024/pkg/gentopendata"
	"stage2024/pkg/helper"
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
	helper.MaybeDie(err, "Failed to start kafka client")
	defer cl.Close()

	slog.Info("Starting schema registry client", "host", *registry)

	rcl, err := sr.NewClient(sr.URLs(*registry))
	helper.MaybeDieErr(err)

	sub := *topic + "-value"
	file, err := os.ReadFile(filepath.Join("./proto", bikes.File_bikes_bolt_proto.Path()))
	helper.MaybeDieErr(err)

	ss, err := rcl.CreateSchema(context.Background(), sub,
		sr.Schema{
			Schema: string(file),
			Type:   sr.TypeProtobuf,
			References: []sr.SchemaReference{
				helper.ReferenceLocation(rcl),
			},
		})
	helper.MaybeDieErr(err)
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
		allBikes := gentopendata.Fetch(url,
			func(b []byte) *bikes.BoltLocation {
				var in ApiData
				err := json.Unmarshal(b, &in)
				if err != nil {
					slog.Warn(err.Error())
					return nil
				}

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
			if err != nil {
				slog.Warn("Bike encoding", "error", err)
				continue
			}
			record := &kgo.Record{Topic: "bolt-test", Value: bikeByte}
			cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
				defer wg.Done()
				if err != nil {
					slog.Warn("record had a produce error",
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
