package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"stage2024/pkg/protogen/stalling"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
)

func main() {
	slog.Default()
	seed := flag.String("seedbroker", "localhost:19092", "brokers port to talk to")
	registry := flag.String("registry", "localhost:18081", "schema registry port to talk to")
	topic := flag.String("topic", "stadskantoor", "topic to produce to and consume from")

	slog.Info("Starting kafka client", "seedbrokers", *seed)
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(*seed),
		kgo.ConsumeTopics("stadskantoor"),
		kgo.ConsumerGroup("Testing"),
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
	//stalling.File_stalling_stadskantoor_proto.Path()
	fmt.Println(stalling.File_stalling_stadskantoor_proto.Path())
	file, err := os.ReadFile("./proto/stalling/stadskantoor.proto")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	sub := *topic + "-value"
	ss, err := rcl.LookupSchema(context.TODO(), sub, sr.Schema{
		Schema: string(file),
		Type:   sr.TypeProtobuf,
	})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("Using schema",
		"subject", ss.Subject,
		"version", ss.Version,
		"id", ss.ID,
	)

	ctx := context.Background()

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
	for {
		fetches := cl.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			panic(errs)
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			var stalling stalling.StallingInfo
			err := serde.Decode(record.Value, &stalling)
			if err != nil {
				slog.Error(err.Error())
				continue
			}
			slog.Info("Pulled record", "record", stalling.GetName())
		}
	}
}
