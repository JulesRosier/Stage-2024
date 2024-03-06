package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"stage2024/pkg/helper"
	"stage2024/pkg/serde"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

func main() {
	slog.Default()

	seed := flag.String("seedbroker", "localhost:19092", "brokers port to talk to")
	registry := flag.String("registry", "localhost:18081", "schema registry port to talk to")
	// topic := flag.String("topic", "bolt-test", "topic to produce to and consume from")

	slog.Info("Starting kafka client", "seedbrokers", *seed)
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(*seed),
		// kgo.ConsumeTopics("baqme-locations"),
		kgo.ConsumeTopics(".*-locations$"),
		kgo.ConsumeRegex(),
		kgo.ConsumerGroup("Testing"),
	)
	if err != nil {
		panic(err)
	}
	defer cl.Close()

	slog.Info("Starting schema registry client", "host", *registry)
	rcl, err := sr.NewClient(sr.URLs(*registry))
	helper.MaybeDieErr(err)

	ctx := context.Background()
	c := make(map[int]*serde.Serde)
	for {
		fetches := cl.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			panic(fmt.Sprint(errs))
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			j, err := handleDecode(ctx, rcl, record.Value, c)
			helper.MaybeDieErr(err)
			fmt.Println(string(j))
			// var bike bikes.BoltLocation
			// err := serde.Decode(record.Value, &bike)
			// if err != nil {
			// 	slog.Error(err.Error())
			// 	continue
			// }
			// slog.Info("Pulled bike",
			// 	"bike_id", bike.GetBikeId(),
			// )
		}

	}
}
