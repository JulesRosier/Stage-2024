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
	"stage2024/pkg/protogen/stalling"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
)

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

func main() {
	slog.Default()

	seed := flag.String("seedbroker", "localhost:19092", "brokers port to talk to")
	registry := flag.String("registry", "localhost:18081", "schema registry port to talk to")
	topic := flag.String("topic", "stalling_gent", "topic to produce to and consume from")

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

	file, err := os.ReadFile(filepath.Join("./proto", stalling.File_stalling_gent_proto.Path()))
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
		&stalling.GentStallingInfo{},
		sr.EncodeFn(func(a any) ([]byte, error) {
			return proto.Marshal(a.(*stalling.GentStallingInfo))
		}),
		sr.Index(1),
		sr.DecodeFn(func(b []byte, a any) error {
			return proto.Unmarshal(b, a.(*stalling.GentStallingInfo))
		}),
	)

	var wg sync.WaitGroup

	slog.Info("Producing records")

	for {
		allItems := gentopendata.Fetch[*stalling.GentStallingInfo](url,
			func(b []byte) *stalling.GentStallingInfo {
				var in ApiData
				json.Unmarshal(b, &in)

				out := stalling.GentStallingInfo{}
				out.Time = in.Time
				out.Facilityname = in.Facilityname
				out.Id = in.Id
				out.Totalplaces = int32(in.Totalplaces)
				out.Freeplaces = int32(in.Freeplaces)
				out.Occupiedplaces = int32(in.Occupiedplaces)
				out.Bezetting = int32(in.Bezetting)
				out.GeoPoint_2D = &stalling.GeoPoint2{
					Lon: in.Geopoint.Lon,
					Lat: in.Geopoint.Lat,
				}

				return &out
			},
		)
		ctx := context.Background()
		for _, item := range allItems {
			wg.Add(1)
			stallingByte, err := serde.Encode(item)
			if err != nil {
				slog.Error("Encoding", "error", err)
				os.Exit(1)
			}
			record := &kgo.Record{Topic: "stalling_gent", Value: stallingByte}
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
