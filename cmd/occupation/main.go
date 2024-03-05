package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"stage2024/pkg/gentopendata"
	h "stage2024/pkg/helper"
	"stage2024/pkg/protogen/common"
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
	h.MaybeDieErr(err)
	defer cl.Close()

	slog.Info("Starting schema registry client", "host", *registry)
	rcl, err := sr.NewClient(sr.URLs(*registry))
	h.MaybeDieErr(err)

	file, err := os.ReadFile(filepath.Join("./proto", occupations.File_occupations_blue_bike_proto.Path()))
	h.MaybeDieErr(err)

	sub := *topic + "-value"
	ss, err := rcl.CreateSchema(context.Background(), sub, sr.Schema{
		Schema: string(file),
		Type:   sr.TypeProtobuf,
		References: []sr.SchemaReference{
			h.ReferenceLocation(rcl),
		},
	})
	h.MaybeDieErr(err)
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
		sr.Index(0),
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
				out.Location = &common.Location{
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
			h.MaybeDie(err, "Encoding error")
			record := &kgo.Record{Topic: *topic, Value: itemByte}
			cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
				defer wg.Done()
				h.MaybeDie(err, "Producing")
			})
		}

		wg.Wait()
		slog.Info("uploaded all records")
		slog.Info("Sleeping", "time", fetchdelay)
		time.Sleep(fetchdelay)
	}
}
