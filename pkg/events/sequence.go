package events

import (
	"context"
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/kafka"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/common"
	"stage2024/pkg/protogen/stations"
	"stage2024/pkg/protogen/users"

	h "stage2024/pkg/helper"

	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func RunSequence(kc *kafka.Client) {
	slog.Info("Starting sequence")
	db := database.GetDb()

	user := database.User{}
	db.Order("random()").Find(&user)

	ctx := context.TODO()

	itemByte, err := kc.Serde.Encode(&bikes.BikeReserved{
		TimeStamp: timestamppb.Now(),
		Bike: &bikes.Bike{
			Bike: &bikes.BikeIdentification{
				Id:    "aaa",
				Model: "bbb",
			},
			IsElectric: false,
		},
		Station: &stations.StationIdentification{
			Id: "aaa",
			Location: &common.Location{
				Latitude:  1.1,
				Longitude: 1.2,
			},
			Name: "cool",
		},
		User: &users.UserIdentification{
			Id:           "aa",
			UserName:     "aaa",
			EmailAddress: "bbb",
		},
	})
	h.MaybeDie(err, "Encoding error")

	topic := "bike_reserved"
	record := &kgo.Record{Topic: topic, Value: itemByte}
	kc.Kcl.Produce(ctx, record, func(_ *kgo.Record, err error) {
		h.MaybeDie(err, "Produce error")
	})
}
