package writestallingstadskantoor

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

const Topic = "stalling-stadskantoor"

func WriteStallingStadskantoor(cl *kgo.Client, serde *sr.Serde) {
	slog.Default()

	var wg sync.WaitGroup

	for {
		allItems := gentopendata.Fetch(url,
			func(b []byte) *occupations.StallingStadskantoor {
				var in ApiData
				err := json.Unmarshal(b, &in)
				h.MaybeDie(err, "Failed to unmarshal")
				out := occupations.StallingStadskantoor{}
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
			h.Produce(serde, cl, &wg, item, ctx, Topic)
		}
		wg.Wait()
		slog.Info("Sleeping for", "duration", fetchdelay, "topic", Topic)
		time.Sleep(fetchdelay)
	}
}
