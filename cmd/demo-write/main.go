package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"stage2024/pkg/protogen/bikes"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
)

type BikeData struct {
	TotalCount int                  `json:"total_count"`
	Results    []bikes.BoltLocation `json:"results"`
}

const maxRequestCount = 100

// updates every 5 minutes
const fetchdelay = time.Minute * 5
const url = "https://data.stad.gent/api/explore/v2.1/catalog/datasets/bolt-deelfietsen-gent/records"

func fetchData(url string, offset int) (BikeData, error) {
	slog.Info("Makking request", "offset", offset)
	resp, err := http.Get(fmt.Sprintf("%s?limit=%d&offset=%d", url, maxRequestCount, offset))
	if err != nil {
		return BikeData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return BikeData{}, err
	}

	var data BikeData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return BikeData{}, err
	}
	slog.Info("Request result", "bike_count", len(data.Results))

	return data, nil
}

func fetchBikes() []bikes.BoltLocation {
	offset := 0
	totalCount := 1

	allBikes := make([]bikes.BoltLocation, 0)

	for offset < totalCount {
		data, err := fetchData(url, offset)
		if err != nil {
			fmt.Println("Error fetching data:", err)
			break
		}
		totalCount = data.TotalCount
		allBikes = append(allBikes, data.Results...)
		offset += len(data.Results)
	}

	slog.Info("Total expented bikes", "count", totalCount)
	slog.Info("Total bikes fetched", "count", len(allBikes))
	return allBikes
}

func main() {
	slog.Default()

	seed := flag.String("seedbroker", "localhost:19092", "brokers port to talk to")
	registry := flag.String("registry", "localhost:18081", "schema registry port to talk to")
	topic := flag.String("topic", "bolt-test", "topic to produce to and consume from")

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

	file, err := os.ReadFile("./proto/bikes/Bolt.proto")
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
		&bikes.BoltLocation{},
		sr.EncodeFn(func(a any) ([]byte, error) {
			return proto.Marshal(a.(*bikes.BoltLocation))
		}),
		sr.Index(1),
		sr.DecodeFn(func(b []byte, a any) error {
			return proto.Unmarshal(b, a.(*bikes.BoltLocation))
		}),
	)

	var wg sync.WaitGroup

	slog.Info("Producing records")
	for {
		allBikes := fetchBikes()
		ctx := context.Background()
		for i := range allBikes {
			bike := &allBikes[i]
			wg.Add(1)
			bikeByte, err := serde.Encode(bike)
			if err != nil {
				slog.Error("Bike encoding", "error", err)
				return
			}
			record := &kgo.Record{Topic: "bolt-test", Value: bikeByte}
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
