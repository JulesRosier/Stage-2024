package writedonkey

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"stage2024/pkg/gentopendata"
	h "stage2024/pkg/helper"
	"stage2024/pkg/kafka"
	"stage2024/pkg/protogen/common"
	"stage2024/pkg/protogen/occupations"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
)

// updates every 10 minutes
const fetchdelay = time.Minute * 10
const url = "https://data.stad.gent/api/explore/v2.1/catalog/datasets/donkey-republic-beschikbaarheid-deelfietsen-per-station/records"

type ApiData struct {
	Station_id          string `json:"station_id"`
	Num_bikes_available int32  `json:"num_bikes_available"`
	Num_docks_available int32  `json:"num_docks_available"`
	Is_renting          int32  `json:"is_renting"`
	Is_installed        int32  `json:"is_installed"`
	Is_returning        int32  `json:"is_returning"`
	Last_reported       string `json:"last_reported"`
	Geopunt             struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}
	Name string `json:"name"`
}

func WriteDonkey(cl *kgo.Client, rcl *sr.Client) {
	slog.Default()

	topic := "donkey-locations"

	// Proto file wordt gelezen, dit moet telkens aangepast worden naar de juiste file
	file, err := os.ReadFile(filepath.Join("./proto", occupations.File_occupations_donkey_proto.Path()))
	h.MaybeDieErr(err)

	// schema ophalen
	ss := kafka.GetSchema(topic, rcl, file)

	var serde sr.Serde
	serde.Register(
		ss.ID,
		&occupations.DonkeyLocation{},
		sr.EncodeFn(func(a any) ([]byte, error) {
			return proto.Marshal(a.(*occupations.DonkeyLocation))
		}),
		sr.Index(0),
		sr.DecodeFn(func(b []byte, a any) error {
			return proto.Unmarshal(b, a.(*occupations.DonkeyLocation))
		}),
	)

	var wg sync.WaitGroup

	for {
		slog.Info("Producing to", "topic", topic)
		allItems := gentopendata.Fetch(url,
			func(b []byte) *occupations.DonkeyLocation {
				var in ApiData
				err := json.Unmarshal(b, &in)
				if err != nil {
					slog.Warn("Failed to unmarshal", "error", err)
				}
				out := occupations.DonkeyLocation{}
				out.StationId = in.Station_id
				out.NumBikesAvailable = in.Num_bikes_available
				out.NumDocksAvailable = in.Num_docks_available
				out.IsRenting = in.Is_renting
				out.IsInstalled = in.Is_installed
				out.IsReturning = in.Is_returning
				out.LastReported = in.Last_reported
				out.Location = &common.Location{
					Lon: in.Geopunt.Lon,
					Lat: in.Geopunt.Lat,
				}
				out.Name = in.Name

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
				h.MaybeDie(err, "Produce error")
			})
		}
		wg.Wait()
		slog.Info("Uploaded all records", "topic", topic)
		slog.Info("Sleeping", "time", fetchdelay)
		time.Sleep(fetchdelay)
	}
}
