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
	"stage2024/pkg/helper"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/common"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
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

func main() {
	slog.Default()

	seed := flag.String("seedbroker", "localhost:19092", "brokers port to talk to")
	registry := flag.String("registry", "localhost:18081", "schema registry port to talk to")
	topic := flag.String("topic", "donkey-locations", "topic to produce to and consume from")

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

	file, err := os.ReadFile(filepath.Join("./proto", bikes.File_bikes_donkey_proto.Path()))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	sub := *topic + "-value"
	ss, err := rcl.CreateSchema(context.Background(), sub, sr.Schema{
		Schema:     string(file),
		Type:       sr.TypeProtobuf,
		References: []sr.SchemaReference{helper.ReferenceLocation(rcl)},
	})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("Created or reusing schema", "subject", sub, "version", ss.Version, "id", ss.ID)

	var serde sr.Serde
	serde.Register(
		ss.ID,
		&bikes.DonkeyLocation{},
		sr.EncodeFn(func(a any) ([]byte, error) {
			return proto.Marshal(a.(*bikes.DonkeyLocation))
		}),
		sr.Index(0),
		sr.DecodeFn(func(b []byte, a any) error {
			return proto.Unmarshal(b, a.(*bikes.DonkeyLocation))
		}),
	)

	var wg sync.WaitGroup

	slog.Info("Producing records")

	for {
		allItems := gentopendata.Fetch[*bikes.DonkeyLocation](url,
			func(b []byte) *bikes.DonkeyLocation {
				var in ApiData
				json.Unmarshal(b, &in)

				out := bikes.DonkeyLocation{}
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
			itemByte, err := serde.Encode(item)
			if err != nil {
				slog.Error("Encoding", "error", err)
				os.Exit(1)
			}
			record := &kgo.Record{Topic: "donkey-locations", Value: itemByte}
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
