package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"stage2024/pkg/database"
	"stage2024/pkg/events"
	"stage2024/pkg/helper"
	h "stage2024/pkg/helper"
	"stage2024/pkg/kafka"
	"stage2024/pkg/opendata"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/common"
	"stage2024/pkg/protogen/stations"
	"stage2024/pkg/protogen/users"
	"stage2024/pkg/scheduler"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/gorm"
)

const maxUser = 100

func main() {
	logLevel := h.GetLogLevel()
	fmt.Printf("LOG_LEVEL = %s\n", logLevel)
	slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	database.Init()

	db := database.GetDb()

	topics := []kafka.Topic{
		{ProtoFile: users.File_users_user_registered_proto, PType: &users.UserRegistered{}},

		{ProtoFile: bikes.File_bikes_bike_abandoned_proto, PType: &bikes.AbandonedBike{}},
		{ProtoFile: bikes.File_bikes_bike_defetect_reported_proto, PType: &bikes.BikeDefectReported{}},
		{ProtoFile: bikes.File_bikes_bike_brought_out_proto, PType: &bikes.BikeBroughtOut{}},
		{ProtoFile: bikes.File_bikes_bike_immobilized_proto, PType: &bikes.BikeImmobilized{}},
		{ProtoFile: bikes.File_bikes_bike_picked_up_proto, PType: &bikes.BikePickedUp{}},
		{ProtoFile: bikes.File_bikes_bike_reserved_proto, PType: &bikes.BikeReserved{}},
		{ProtoFile: bikes.File_bikes_bike_returned_proto, PType: &bikes.BikeReturned{}},
		{ProtoFile: bikes.File_bikes_bike_stored_proto, PType: &bikes.BikeStored{}},

		{ProtoFile: stations.File_stations_station_capacity_decreased_proto, PType: &stations.StationCapacityDecreased{}},
		{ProtoFile: stations.File_stations_station_capacity_exhausted_proto, PType: &stations.StationCapacityExhausted{}},
		{ProtoFile: stations.File_stations_station_capacity_increased_proto, PType: &stations.StationCapacityIncreased{}},
		{ProtoFile: stations.File_stations_station_created_proto, PType: &stations.StationCreated{}},
		{ProtoFile: stations.File_stations_station_deprecated_proto, PType: &stations.StationDeprecated{}},
	}

	kc := kafka.NewClient(kafka.Config{
		Topics: topics,
	})

	events.RunSequence(kc)

	changesCh := make(chan []string, 100)

	s := scheduler.NewScheduler()

	CreateUsers(db, kc)
	// s.Schedule(time.Minute*5, func() { opendata.Bolt(db, changesCh) })
	test1 := struct {
		BikeId             string `json:"bike_id"`
		CurrentRangeMeters int32  `json:"current_range_meters"`
		PricingPlanId      string `json:"pricing_plan_id"`
		VehicleTypeId      string `json:"vehicle_type_id"`
		IsReserved         int32  `json:"is_reserved"`
		IsDisabled         int32  `json:"is_disabled"`
		RentalUris         string `json:"rental_uris"`
		Loc                struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		}
	}{}

	model := "bolt"
	p := opendata.Puller[*database.Bike]{
		Model: model,
		Url:   "https://data.stad.gent/api/explore/v2.1/catalog/datasets/bolt-deelfietsen-gent/records",
		Transform: func(b []byte) *database.Bike {
			faker := gofakeit.New(0) // seed to get same random bool values each time
			in := test1
			err := json.Unmarshal(b, &in)
			helper.MaybeDieErr(err)

			out := &database.Bike{}
			out.Id = in.BikeId
			out.BikeModel = model
			out.IsElectric = sql.NullBool{Bool: faker.Bool(), Valid: true} //random bool value
			out.Location = fmt.Sprintf("%v", &common.Location{
				Latitude:  in.Loc.Lat,
				Longitude: in.Loc.Lon,
			})
			out.IsImmobilized = sql.NullBool{Valid: false} //fake
			out.IsAbandoned = sql.NullBool{Valid: false}   //fake
			out.IsAvailable = sql.NullBool{Bool: in.IsDisabled == 0, Valid: true}
			out.IsInStorage = sql.NullBool{Valid: false} //fake
			out.IsReserved = sql.NullBool{Bool: in.IsReserved != 0, Valid: true}
			out.IsDefect = sql.NullBool{Valid: false} //fake
			return out
		},
	}
	s.Schedule(time.Second*30, func() { p.Pul(db, changesCh) })

	go func() {
		for item := range changesCh {
			slog.Info("[" + strings.Join(item, `, `) + `]`)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("Received an interrupt signal, exiting...")

	s.Stop()

	slog.Info("Exiting... Goodbye!")
}

func CreateUsers(db *gorm.DB, kc *kafka.Client) {
	var userCount int64
	db.Model(&database.User{}).Count(&userCount)

	for userCount < maxUser {
		events.CreateRandomUser(db, kc)
		userCount++
	}
}
