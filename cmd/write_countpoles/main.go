package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	h "stage2024/pkg/helper"
	"stage2024/pkg/protogen/poles"
	"strconv"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Historical data, so data is from 2023

// updates every 5 minutes
const fetchdelay = time.Minute * 5

type ApiData struct {
	Code          string    `json:"code"`
	Locatie       string    `json:"locatie"`
	Datum         string    `json:"datum"`
	Uur5minuten   string    `json:"uur5minuten"`
	Ordening      time.Time `json:"ordening"`
	Totaal        string    `json:"totaal"`
	Tegenrichting string    `json:"tegenrichting"`
	Hoofdrichting string    `json:"hoofdrichting"`
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		slog.Warn("unable to convert string to int", "string", s)
		return -1
	}
	return i
}

func main() {
	slog.Default()

	seed := flag.String("seedbroker", "localhost:19092", "brokers port to talk to")
	registry := flag.String("registry", "localhost:18081", "schema registry port to talk to")
	topic := flag.String("topic", "countpoles", "topic to produce to and consume from")

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

	file, err := os.ReadFile(filepath.Join("./proto", poles.File_poles_pole_proto.Path()))
	h.MaybeDieErr(err)

	sub := *topic + "-value"
	ss, err := rcl.CreateSchema(context.Background(), sub, sr.Schema{
		Schema: string(file),
		Type:   sr.TypeProtobuf,
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
		&poles.PoleData{},
		sr.EncodeFn(func(a any) ([]byte, error) {
			return proto.Marshal(a.(*poles.PoleData))
		}),
		sr.Index(0),
		sr.DecodeFn(func(b []byte, a any) error {
			return proto.Unmarshal(b, a.(*poles.PoleData))
		}),
	)

	var wg sync.WaitGroup

	slog.Info("Producing records")
	for {
		wg.Add(1)
		go func() {
			slog.Info("Sleeping", "time", fetchdelay)
			time.Sleep(fetchdelay)
			wg.Done()
		}()

		urls := h.GetPoleUrls()

		for _, url := range urls {
			allItems := h.Fetch(url,
				func(b []byte) *poles.PoleData {
					var in ApiData
					json.Unmarshal(b, &in)

					out := poles.PoleData{}
					out.Code = in.Code
					out.Locatie = in.Locatie
					out.Datum = in.Datum
					out.Uur5Minuten = in.Uur5minuten
					out.Ordening = timestamppb.New(in.Ordening)
					out.Totaal = int32(toInt(in.Totaal))
					out.Tegenrichting = int32(toInt(in.Tegenrichting))
					out.Hoofdrichting = int32(toInt(in.Hoofdrichting))

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
			slog.Info("fetched data", "src", url)
		}

		wg.Wait()
		slog.Info("uploaded all records")
	}
}