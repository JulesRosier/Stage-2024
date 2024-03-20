package writecountpoles

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	h "stage2024/pkg/helper"
	"stage2024/pkg/kafka"
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
	Code             string    `json:"code"`
	Locatie          string    `json:"locatie"`
	Datum            string    `json:"datum"`
	Uur5minuten      string    `json:"uur5minuten"`
	Ordening         time.Time `json:"ordening"`
	Totaal           string    `json:"totaal"`
	Tegenrichting    string    `json:"tegenrichting"`
	Hoofdrichting    string    `json:"hoofdrichting"`
	DatumUur5minuten time.Time `json:"datum5minuten"`
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		slog.Warn("unable to convert string to int", "string", s)
		return -1
	}
	return i
}

func WriteCountples(cl *kgo.Client, rcl *sr.Client) {
	slog.Default()

	topic := "countpoles"

	// Proto file wordt gelezen, dit moet dus ieder keer aangepast worden naar de juiste file
	file, err := os.ReadFile(filepath.Join("./proto", poles.File_poles_pole_proto.Path()))
	h.MaybeDieErr(err)

	// schema ophalen
	ss := kafka.GetSchema(topic, rcl, file)

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

	for {
		slog.Info("Producing to", "topic", topic)
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
					h.MaybeDieErr(err)

					out := poles.PoleData{}
					out.Code = in.Code
					out.Locatie = in.Locatie
					out.Datum = in.Datum
					out.Uur5Minuten = in.Uur5minuten
					out.Ordening = timestamppb.New(in.Ordening)
					out.Totaal = int32(toInt(in.Totaal))
					out.Tegenrichting = int32(toInt(in.Tegenrichting))
					out.Hoofdrichting = int32(toInt(in.Hoofdrichting))

					datum5min, err := time.Parse("2006-01-02 15:04:05", in.Datum+" "+in.Uur5minuten)
					if err != nil {
						slog.Warn("unable to parse date", "date", in.Datum+" "+in.Uur5minuten, "error", err.Error())
						datum5min = time.Time{}
					}
					out.Datum5Minuten = timestamppb.New(datum5min)

					return &out
				},
			)
			ctx := context.Background()
			for _, item := range allItems {
				wg.Add(1)
				itemByte, err := serde.Encode(item)
				h.MaybeDie(err, "Encoding error")
				record := &kgo.Record{Topic: topic, Value: itemByte}
				cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
					defer wg.Done()
					h.MaybeDie(err, "Producing")
				})
			}
		}
		slog.Info("uploaded all records", "topic", topic)
		wg.Wait()
	}
}
