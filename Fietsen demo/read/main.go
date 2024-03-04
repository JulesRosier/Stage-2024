package main

import (
	"context"
	"demo-read/pb"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
)

func main() {
	slog.Default()

	seed := flag.String("seedbroker", "localhost:19092", "brokers port to talk to")
	registry := flag.String("registry", "localhost:18081", "schema registry port to talk to")
	topic := flag.String("topic", "bolt-test", "topic to produce to and consume from")

	slog.Info("Starting kafka client", "seedbrokers", *seed)
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(*seed),
		kgo.ConsumeTopics("bolt-test"),
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

	file, err := os.ReadFile("../bike.proto")
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
		&pb.Bike{},
		sr.EncodeFn(func(a any) ([]byte, error) {
			return proto.Marshal(a.(*pb.Bike))
		}),
		sr.Index(1),
		sr.DecodeFn(func(b []byte, a any) error {
			return proto.Unmarshal(b, a.(*pb.Bike))
		}),
	)

	for {
		fetches := cl.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			panic(fmt.Sprint(errs))
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			var bike pb.Bike
			err := serde.Decode(record.Value, &bike)
			if err != nil {
				slog.Error(err.Error())
				continue
			}
			slog.Info("Pulled bike",
				"bike_id", bike.GetBikeId(),
			)
		}

	}
}
