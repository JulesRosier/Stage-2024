package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"stage2024/pkg/gentopendata"
	"stage2024/pkg/protogen/occupations"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// updates every 5 minutes
const fetchdelay = time.Minute * 5

const url = "https://data.stad.gent/api/explore/v2.1/catalog/datasets/blue-bike-deelfietsen-gent-sint-pieters-m-hendrikaplein/records"

type ApiData struct {
	LastSeen       time.Time `json:"last_seen"`
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	BikesInUse     int       `json:"bikes_in_use"`
	BikesAvailable int       `json:"bikes_available"`
	Longitude      string    `json:"longitude"`
	Latitude       string    `json:"latitude"`
	GeoPoint       struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"geopoint"`
	Type string `json:"type"`
}

func main() {
	slog.Default()

	seed := flag.String("seedbroker", "localhost:19092", "brokers port to talk to")
	registry := flag.String("registry", "localhost:18081", "schema registry port to talk to")
	topic := flag.String("topic", "bluebike-test", "topic to produce to and consume from")

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

	file, err := os.ReadFile("./proto/occupations/BlueBike.proto")
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
	slog.Info("created or reusing schema",
		"subject", ss.Subject,
		"version", ss.Version,
		"id", ss.ID,
	)

	var serde sr.Serde
	serde.Register(
		ss.ID,
		&occupations.BlueBikeOccupation{},
		sr.EncodeFn(func(a any) ([]byte, error) {
			return proto.Marshal(a.(*occupations.BlueBikeOccupation))
		}),
		sr.Index(1),
		sr.DecodeFn(func(b []byte, a any) error {
			return proto.Unmarshal(b, a.(*occupations.BlueBikeOccupation))
		}),
	)

	var wg sync.WaitGroup

	slog.Info("Producing records")
	for {
		allItems := gentopendata.Fetch[*occupations.BlueBikeOccupation](url,
			func(b []byte) *occupations.BlueBikeOccupation {
				var in ApiData
				json.Unmarshal(b, &in)

				out := occupations.BlueBikeOccupation{}
				out.LastSeen = timestamppb.New(in.LastSeen)
				out.Id = int32(in.Id)
				out.Name = in.Name
				out.BikesInUse = int32(in.BikesInUse)
				out.BikesAvailable = int32(in.BikesAvailable)
				out.Geopoint = &occupations.GeoPoint{
					Lon: in.GeoPoint.Lon,
					Lat: in.GeoPoint.Lat,
				}
				out.Type = in.Type

				return &out
			},
		)
		ctx := context.Background()
		for _, item := range allItems {
			wg.Add(1)
			itemByte, err := serde.Encode(item)
			if err != nil {
				slog.Error("Encoding", "error", err)
				return
			}
			record := &kgo.Record{Topic: *topic, Value: itemByte}
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
