package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"stage2024/pkg/gentopendata"
	h "stage2024/pkg/helper"
	"stage2024/pkg/protogen/common"
	"stage2024/pkg/protogen/events"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
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

func main() {
	slog.Default()

	seed := flag.String("seedbroker", "localhost:19092", "brokers port to talk to")
	registry := flag.String("registry", "localhost:18081", "schema registry port to talk to")
	topic := flag.String("topic", "bike.droppedoff", "topic to produce to and consume from")

	slog.Info("Starting kafka client...", "seedbroker", *seed)
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(*seed),
		kgo.AllowAutoTopicCreation(),
	)
	h.MaybeDieErr(err)
	defer cl.Close()

	slog.Info("Starting schema registry client...", "host", *registry)
	rcl, err := sr.NewClient(sr.URLs(*registry))
	h.MaybeDieErr(err)

	file, err := os.ReadFile(filepath.Join("./proto", events.File_events_bike_pickup_proto.Path()))
	h.MaybeDieErr(err)

	sub := *topic + "-value"
	ss, err := rcl.CreateSchema(context.Background(), sub, sr.Schema{
		Schema:     string(file),
		Type:       sr.TypeProtobuf,
		References: []sr.SchemaReference{h.ReferenceLocation(rcl)},
	})
	h.MaybeDieErr(err)
	slog.Info("Created or reusing schema", "subject", sub, "version", ss.Version, "id", ss.ID)

	var serde sr.Serde
	serde.Register(
		ss.ID,
		&events.BikeDropOff{},
		sr.EncodeFn(func(a any) ([]byte, error) {
			return proto.Marshal(a.(*events.BikeDropOff))
		}),
		sr.Index(0),
		sr.DecodeFn(func(b []byte, a any) error {
			return proto.Unmarshal(b, a.(*events.BikeDropOff))
		}),
	)

	var wg sync.WaitGroup

	slog.Info("Producing records")

	var prevItems []ApiData

	for {
		allItems := gentopendata.Fetch(url,
			func(b []byte) ApiData {
				var in ApiData
				err := json.Unmarshal(b, &in)
				if err != nil {
					slog.Warn("Failed to unmarshal", "error", err)
				}
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

				e := events.BikeDropOff{
					Reported:  timestamppb.New(time.Unix(unixTime, 0)),
					StationId: item.Station_id,
					Location: &common.Location{
						Lon: item.Geopunt.Lon,
						Lat: item.Geopunt.Lat,
					},
					Company:        "donkey-republic",
					Dif:            bikesDiv,
					DocksAvailable: item.Num_docks_available,
					BikesAvailable: item.Num_bikes_available,
				}
				ctx := context.TODO()

				wg.Add(1)
				itemByte, err := serde.Encode(&e)
				h.MaybeDie(err, "Encoding error")
				record := &kgo.Record{Topic: *topic, Value: itemByte}
				cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
					defer wg.Done()
					h.MaybeDie(err, "Produce error")
				})

			}
		}

		prevItems = allItems

		wg.Wait()
		slog.Info("Uploaded all records")
		slog.Info("Sleeping for", "duration", fetchdelay)
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
