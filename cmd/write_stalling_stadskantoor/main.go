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
)

const fetchdelay = time.Minute * 5
const url = "https://data.stad.gent/api/explore/v2.1/catalog/datasets/real-time-bezetting-fietsenstalling-stadskantoor-gent/records"

type ApiData struct {
	Name            string  `json:"name"`
	Parkingcapacity float32 `json:"parkingCapacity"`
	Vacantspaces    float32 `json:"vacantSpaces"`
	Naam            string  `json:"naam"`
	Parking         string  `json:"parking"`
	Occupation      int32   `json:"occupation"`
	Infotekst       string  `json:"infotekst"`
	Enginfotekst    string  `json:"enginfotekst"`
	Frinfotekst     string  `json:"frinfotekst"`
	Locatie         struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"locatie"`
}

func main() {
	slog.Default()

	seed := flag.String("seedbroker", "localhost:19092", "brokers port to talk to")
	registry := flag.String("registry", "localhost:18081", "schema registry port to talk to")
	topic := flag.String("topic", "stalling-stadskantoor", "topic to produce to and consume from")

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

	file, err := os.ReadFile(filepath.Join("./proto", occupations.File_occupations_stadskantoor_proto.Path()))
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
		&occupations.StallingInfo{},
		sr.EncodeFn(func(a any) ([]byte, error) {
			return proto.Marshal(a.(*occupations.StallingInfo))
		}),
		sr.Index(0),
		sr.DecodeFn(func(b []byte, a any) error {
			return proto.Unmarshal(b, a.(*occupations.StallingInfo))
		}),
	)

	var wg sync.WaitGroup

	slog.Info("Producing records")

	for {
		allItems := gentopendata.Fetch[*occupations.StallingInfo](url,
			func(b []byte) *occupations.StallingInfo {
				var in ApiData
				err := json.Unmarshal(b, &in)
				if err != nil {
					slog.Warn("Failed to unmarshal", "error", err)
					return nil
				}
				out := occupations.StallingInfo{}
				out.Name = in.Name
				out.Parkingcapacity = int32(in.Parkingcapacity)
				out.Vacantspaces = int32(in.Vacantspaces)
				out.Naam = in.Naam
				out.Parking = in.Parking
				out.Occupation = in.Occupation
				out.Infotekst = in.Infotekst
				out.Enginfotekst = in.Enginfotekst
				out.Frinfotekst = in.Frinfotekst
				out.Location = &common.Location{
					Lon: in.Locatie.Lon,
					Lat: in.Locatie.Lat,
				}

				return &out
			},
		)
		ctx := context.Background()
		for _, item := range allItems {
			wg.Add(1)
			stallingByte, err := serde.Encode(item)
			h.MaybeDie(err, "Encoding")

			record := &kgo.Record{Topic: "stalling-stadskantoor", Value: stallingByte}
			cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
				defer wg.Done()
				h.MaybeDie(err, "Produce error")
			})
		}
		wg.Wait()
		slog.Info("Uploaded all records")
		slog.Info("Sleeping for", "duration", fetchdelay)
		time.Sleep(fetchdelay)
	}
}
