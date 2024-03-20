package writestallinggent

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

const fetchdelay = time.Minute * 1
const url = "https://data.stad.gent/api/explore/v2.1/catalog/datasets/real-time-bezettingen-fietsenstallingen-gent/records"

type ApiData struct {
	Time           string `json:"time"`
	Facilityname   string `json:"facilityname"`
	Id             string `json:"id"`
	Totalplaces    int32  `json:"totalplaces"`
	Freeplaces     int32  `json:"freeplaces"`
	Occupiedplaces int32  `json:"occupiedplaces"`
	Bezetting      int32  `json:"bezetting"`
	Geopoint       struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}
}

func WriteStallingGent(cl *kgo.Client, rcl *sr.Client) {
	slog.Default()

	topic := "stalling-gent"

	// Proto file wordt gelezen, dit moet dus ieder keer aangepast worden naar de juiste file
	file, err := os.ReadFile(filepath.Join("./proto", occupations.File_occupations_gentstalling_proto.Path()))
	h.MaybeDieErr(err)

	// schema ophalen
	ss := kafka.GetSchema(topic, rcl, file)

	//truct van gegenereerde file wordt hier gebruikt
	var serde sr.Serde
	serde.Register(
		ss.ID,
		&occupations.GentStallingInfo{},
		sr.EncodeFn(func(a any) ([]byte, error) {
			return proto.Marshal(a.(*occupations.GentStallingInfo))
		}),
		sr.Index(0),
		sr.DecodeFn(func(b []byte, a any) error {
			return proto.Unmarshal(b, a.(*occupations.GentStallingInfo))
		}),
	)

	var wg sync.WaitGroup

	for {
		slog.Info("Producing to", "topic", topic)
		allItems := gentopendata.Fetch(url,
			func(b []byte) *occupations.GentStallingInfo {
				var in ApiData
				err := json.Unmarshal(b, &in)
				if err != nil {
					slog.Warn("Failed to unmarshal", "error", err)
				}
				out := occupations.GentStallingInfo{}
				out.Time = in.Time
				out.Facilityname = in.Facilityname
				out.Id = in.Id
				out.Totalplaces = int32(in.Totalplaces)
				out.Freeplaces = int32(in.Freeplaces)
				out.Occupiedplaces = int32(in.Occupiedplaces)
				out.Bezetting = int32(in.Bezetting)
				out.Location = &common.Location{
					Lon: in.Geopoint.Lon,
					Lat: in.Geopoint.Lat,
				}

				return &out
			},
		)
		ctx := context.Background()
		for _, item := range allItems {
			wg.Add(1)
			stallingByte, err := serde.Encode(item)
			h.MaybeDie(err, "Encoding error")
			record := &kgo.Record{Topic: topic, Value: stallingByte}
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
