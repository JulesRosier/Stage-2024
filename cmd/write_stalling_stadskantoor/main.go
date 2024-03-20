package writestallingstadskantoor

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"stage2024/pkg/gentopendata"
	h "stage2024/pkg/helper"
	"stage2024/pkg/kafka"
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

func WriteStallingStadskantoor(cl *kgo.Client, rcl *sr.Client) {
	slog.Default()

	topic := "stalling-stadskantoor"

	// Proto file wordt gelezen, dit moet dus ieder keer aangepast worden naar de juiste file
	file, err := os.ReadFile(filepath.Join("./proto", occupations.File_occupations_stadskantoor_proto.Path()))
	h.MaybeDieErr(err)

	// schema ophalen
	ss := kafka.GetSchema(topic, rcl, file)

	// Hier telkens anders, struct van gegenereerde file wordt hier gebruikt
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

	for {
		slog.Info("Producing to", "topic", topic)
		allItems := gentopendata.Fetch(url,
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

			record := &kgo.Record{Topic: topic, Value: stallingByte}
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
