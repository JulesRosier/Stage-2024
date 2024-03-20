package writebluebike

import (
	"context"
	"encoding/json"
	"log/slog"
	"stage2024/pkg/gentopendata"
	h "stage2024/pkg/helper"
	"stage2024/pkg/protogen/common"
	"stage2024/pkg/protogen/occupations"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const Topic = "bluebike-locations"

// updates every 5 minutes
const fetchdelay = time.Minute * 5

type ApiData struct {
	LastSeen       time.Time `json:"last_seen"`
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	BikesInUse     int       `json:"bikes_in_use"`
	BikesAvailable int       `json:"bikes_available"`
	Longitude      string    `json:"longitude"`
	Latitude       string    `json:"latitude"`
	GeoPoint       struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"geopoint"`
	Type string `json:"type"`
}

func WriteBluebike(cl *kgo.Client, serde *sr.Serde) {
	slog.Default()

	urls := []string{"https://data.stad.gent/api/explore/v2.1/catalog/datasets/blue-bike-deelfietsen-gent-sint-pieters-m-hendrikaplein/records",
		"https://data.stad.gent/api/explore/v2.1/catalog/datasets/blue-bike-deelfietsen-gent-dampoort/records",
		"https://data.stad.gent/api/explore/v2.1/catalog/datasets/blue-bike-deelfietsen-gent-sint-pieters-st-denijslaan/records",
		"https://data.stad.gent/api/explore/v2.1/catalog/datasets/blue-bike-deelfietsen-merelbeke-drongen-wondelgem/records"}

	var wg sync.WaitGroup

	for {
		wg.Add(1)
		go func() {
			slog.Info("Sleeping for", "duration", fetchdelay, "topic", Topic)
			time.Sleep(fetchdelay)
			wg.Done()
		}()
		for _, url := range urls {
			allItems := gentopendata.Fetch(url,
				func(b []byte) *occupations.BlueBikeOccupation {
					var in ApiData
					json.Unmarshal(b, &in)

					out := occupations.BlueBikeOccupation{}
					out.LastSeen = timestamppb.New(in.LastSeen)
					out.Id = int32(in.Id)
					out.Name = in.Name
					out.BikesInUse = int32(in.BikesInUse)
					out.BikesAvailable = int32(in.BikesAvailable)
					out.Location = &common.Location{
						Lon: in.GeoPoint.Lon,
						Lat: in.GeoPoint.Lat,
					}
					out.Type = in.Type

					return &out
				},
			)
			ctx := context.Background()
			for _, item := range allItems {
				wg.Add(1)
				itemByte, err := serde.Encode(item)
				h.MaybeDie(err, "Encoding error")
				record := &kgo.Record{Topic: Topic, Value: itemByte}
				cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
					defer wg.Done()
					h.MaybeDie(err, "Producing")
				})
			}
		}
		wg.Wait()
	}
}
