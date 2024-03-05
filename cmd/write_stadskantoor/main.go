package main

import (
	"context"
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

const fetchdelay = time.Minute * 5
const url = "https://data.stad.gent/api/explore/v2.1/catalog/datasets/real-time-bezetting-fietsenstalling-stadskantoor-gent/records"

func main() {
	slog.Default()

	seed := flag.String("seedbroker", "localhost:19092", "brokers port to talk to")
	registry := flag.String("registry", "localhost:18081", "schema registry port to talk to")
	topic := flag.String("topic", "stadskantoor", "topic to produce to and consume from")

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

	file, err := os.ReadFile(filepath.Join("./proto", stalling.File_stalling_stadskantoor_proto.Path()))
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
		&stalling.StallingInfo{},
		sr.EncodeFn(func(a any) ([]byte, error) {
			return proto.Marshal(a.(*stalling.StallingInfo))
		}),
		sr.Index(1),
		sr.DecodeFn(func(b []byte, a any) error {
			return proto.Unmarshal(b, a.(*stalling.StallingInfo))
		}),
	)

	var wg sync.WaitGroup

	slog.Info("Producing records")

	for {
		stallingen := gentopendata.Fetch[stalling.StallingInfo](url)
		ctx := context.Background()
		for i := range stallingen {
			stalling := &stallingen[i]
			wg.Add(1)
			stallingByte, err := serde.Encode(stalling)
			if err != nil {
				slog.Error("Bike encoding", "error", err)
				os.Exit(1)
			}
			record := &kgo.Record{Topic: "stadskantoor", Value: stallingByte}
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
