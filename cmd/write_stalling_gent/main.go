package writestallinggent

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"stage2024/pkg/gentopendata"
	h "stage2024/pkg/helper"
	"stage2024/pkg/protogen/common"
	"stage2024/pkg/protogen/occupations"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
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

const Topic = "stalling-gent"

func WriteStallingGent(cl *kgo.Client, serde *sr.Serde) {
	slog.Default()

	var wg sync.WaitGroup

	for {
		allItems := gentopendata.Fetch(url,
			func(b []byte) *occupations.StallingGent {
				var in ApiData
				err := json.Unmarshal(b, &in)
				if err != nil {
					slog.Warn("Failed to unmarshal", "error", err)
				}
				out := occupations.StallingGent{}
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
			record := &kgo.Record{Topic: Topic, Value: stallingByte}
			cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
				defer wg.Done()
				h.MaybeDie(err, "Produce error")
			})
		}
		wg.Wait()
		slog.Info("Sleeping for", "duration", fetchdelay, "topic", Topic)
		time.Sleep(fetchdelay)
	}
}
