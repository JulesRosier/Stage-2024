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
		{ProtoFile: bikes.File_bikes_bike_deployed_proto, PType: &bikes.BikeDeployed{}},
		{ProtoFile: bikes.File_bikes_bike_immobilized_proto, PType: &bikes.BikeImmobilized{}},
		{ProtoFile: bikes.File_bikes_bike_picked_up_proto, PType: &bikes.BikePickedUp{}},
		{ProtoFile: bikes.File_bikes_bike_reserved_proto, PType: &bikes.BikeReserved{}},
		{ProtoFile: bikes.File_bikes_bike_returned_proto, PType: &bikes.BikeReturned{}},
		{ProtoFile: bikes.File_bikes_bike_stored_proto, PType: &bikes.BikeStored{}},

		{ProtoFile: stations.File_stations_station_occupation_increased_proto, PType: &stations.StationOccupationDecreased{}},
		{ProtoFile: stations.File_stations_station_full_proto, PType: &stations.StationFull{}},
		{ProtoFile: stations.File_stations_station_occupation_decreased_proto, PType: &stations.StationOccupationIncreased{}},
		{ProtoFile: stations.File_stations_station_created_proto, PType: &stations.StationCreated{}},
		{ProtoFile: stations.File_stations_station_deprecated_proto, PType: &stations.StationDeprecated{}},
	}

	kc := kafka.NewClient(kafka.Config{
		Topics: topics,
	})

	ol := kafka.NewOutboxListener(kc, db, topics)
	ol.Start()

	i := 0
	for i < 10 {
		i++
		//go events.RunSequence(kc)
	}

	s := scheduler.NewScheduler()

	ec := events.NewEventClient(kc)

	CreateUsers(db, kc)
	s.Schedule(time.Minute*5, func() { opendata.Bolt(ec.Channel) })
	s.Schedule(time.Minute*10, func() { opendata.Baqme(ec.Channel) })
	s.Schedule(time.Minute*5, func() { opendata.BlueBike(ec.Channel) })
	s.Schedule(time.Minute*10, func() { opendata.Donkey(ec.Channel) })
	s.Schedule(time.Minute*1, func() { opendata.StorageGhent(ec.Channel) })
	s.Schedule(time.Minute*5, func() { opendata.StorageTownHall(ec.Channel) })

	s.Schedule(time.Second*30, ol.FetchOutboxData)

	go func() {
		for change := range ec.Channel {
			ec.ChangeDetected(change)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("Received an interrupt signal, exiting...")

	s.Stop()

	slog.Info("Exiting... Goodbye!")
}

func CreateUsers(db *gorm.DB, kc *kafka.KafkaClient) {
	var userCount int64
	db.Model(&database.User{}).Count(&userCount)

	for userCount < maxUser {
		events.CreateRandomUser(db, kc)
		userCount++
	}
}
