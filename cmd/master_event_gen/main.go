package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"stage2024/pkg/database"
	"stage2024/pkg/events"
	h "stage2024/pkg/helper"
	"stage2024/pkg/kafka"
	"stage2024/pkg/opendata"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/stations"
	"stage2024/pkg/protogen/users"
	"stage2024/pkg/scheduler"
	"time"

	"gorm.io/gorm"
)

const maxUser = 100

func main() {
	fmt.Println("Starting...")
	logLevel := h.GetLogLevel()
	fmt.Printf("LOG_LEVEL = %s\n", logLevel)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)
	database.Init()

	db := database.GetDb()

	topics := []kafka.Topic{
		{ProtoFile: users.File_users_user_registered_proto, PType: &users.UserRegistered{}},

		{ProtoFile: bikes.File_bikes_bike_abandoned_proto, PType: &bikes.BikeAbandoned{}},
		{ProtoFile: bikes.File_bikes_bike_defect_reported_proto, PType: &bikes.BikeDefectReported{}},
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

	i := 0
	for i < 10 {
		i++
		go events.RunSequence(kc)
	}

	changesCh := make(chan h.Change, 100)

	s := scheduler.NewScheduler()

	CreateUsers(db, kc)
	s.Schedule(time.Minute*5, func() { opendata.Bolt(changesCh) })
	s.Schedule(time.Minute*10, func() { opendata.Baqme(changesCh) })
	s.Schedule(time.Minute*5, func() { opendata.BlueBike(changesCh) })
	s.Schedule(time.Minute*10, func() { opendata.Donkey(changesCh) })
	s.Schedule(time.Minute*1, func() { opendata.StorageGhent(changesCh) })
	s.Schedule(time.Minute*5, func() { opendata.StorageTownHall(changesCh) })

	go func() {
		for change := range changesCh {
			ChangeDetected(kc, change)
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
