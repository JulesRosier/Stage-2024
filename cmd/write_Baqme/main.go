package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"stage2024/pkg/gentopendata"
	"stage2024/pkg/protogen/bikes"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
)

const fetchdelay = time.Minute * 10
const url = "https://data.stad.gent/api/explore/v2.1/catalog/datasets/baqme-locaties-vrije-deelfietsen-gent/records"

type ApiData struct {
	Bikeid          string  `json:"bike_id"`
	Lat             float32 `json:"lat"`
	Lon             float32 `json:"lon"`
	Is_reserved     int32   `json:"is_reserved"`
	Is_disabled     int32   `json:"is_disabled"`
	Vehicle_type_id string  `json:"vehicle_type"`
	Rental_uris     string  `json:"rental_uris"`
	Geopoint        struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}
}

func main() {
	slog.Default()

	seed := flag.String("seedbroker", "localhost:19092", "brokers port to talk to")
	registry := flag.String("registry", "localhost:18081", "schema registry port to talk to")
	topic := flag.String("topic", "baqme-locations", "topic to produce to and consume from")

	slog.Info("Starting kafka client...", "seedbroker", *seed)
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(*seed),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		panic(err)
	}
	defer cl.Close()

	slog.Info("Starting schema registry client...", "host", *registry)
	rcl, err := sr.NewClient(sr.URLs(*registry))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	file, err := os.ReadFile(filepath.Join("./proto", bikes.File_bikes_Baqme_proto.Path()))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	sub := *topic + "-value"
	ss, err := rcl.CreateSchema(context.Background(), sub, sr.Schema{
		Schema: string(file),
		Type:   sr.TypeProtobuf,
	})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("Created or reusing schema", "subject", sub, "version", ss.Version, "id", ss.ID)

	var serde sr.Serde
	serde.Register(
		ss.ID,
		&bikes.BaqmeLocation{},
		sr.EncodeFn(func(a any) ([]byte, error) {
			return proto.Marshal(a.(*bikes.BaqmeLocation))
		}),
		sr.Index(1),
		sr.DecodeFn(func(b []byte, a any) error {
			return proto.Unmarshal(b, a.(*bikes.BaqmeLocation))
		}),
	)

	var wg sync.WaitGroup

	slog.Info("Producing records")

	for {
		allItems := gentopendata.Fetch[*bikes.BaqmeLocation](url,
			func(b []byte) *bikes.BaqmeLocation {
				var in ApiData
				json.Unmarshal(b, &in)

				out := bikes.BaqmeLocation{}
				out.BikeId = in.Bikeid
				out.Lat = in.Lat
				out.Lon = in.Lon
				out.IsReserved = in.Is_reserved
				out.IsDisabled = in.Is_disabled
				out.VehicleTypeId = in.Vehicle_type_id
				out.RentalUris = in.Rental_uris
				out.Geopoint = &bikes.Geopoint{
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
			if err != nil {
				slog.Error("Encoding", "error", err)
				os.Exit(1)
			}
			record := &kgo.Record{Topic: "baqme-locations", Value: itemByte}
			cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
				defer wg.Done()
				if err != nil {
					slog.Error("Produce error", "error", err)
				}
			})
		}
		wg.Wait()
		slog.Info("Uploaded all records")
		slog.Info("Sleeping for", "duration", fetchdelay)
		time.Sleep(fetchdelay)
	}
}
