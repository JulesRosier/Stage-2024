package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"stage2024/pkg/database"
	h "stage2024/pkg/helper"
	"stage2024/pkg/kafka"
	"stage2024/pkg/opendata"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/stations"
	"stage2024/pkg/protogen/users"
	"stage2024/pkg/scheduler"
	"time"

	"github.com/brianvoe/gofakeit"
	"gorm.io/gorm"
)

const maxUser = 100
const fakeBikefrequency = 60

func main() {
	fmt.Println("Starting...")
	logLevel := h.GetLogLevel()
	fmt.Printf("LOG_LEVEL = %s\n", logLevel)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	dbc := database.NewDatabase()

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

	ol := kafka.NewOutboxListener(kc, dbc, topics)
	ol.Start()

	s := scheduler.NewScheduler()

	CreateUsers(dbc.DB, kc)

	//TODO: Uncomment
	s.Schedule(time.Minute*5, func() { opendata.BlueBike(dbc.DB) })
	s.Schedule(time.Minute*10, func() { opendata.Donkey(dbc.DB) })
	s.Schedule(time.Minute*1, func() { opendata.StorageGhent(dbc.DB) })
	s.Schedule(time.Minute*5, func() { opendata.StorageTownHall(dbc.DB) })

	// s.Schedule(time.Second*10, ol.FetchOutboxData)

	//Generate bike events every x minutes
	s.Schedule(time.Minute*fakeBikefrequency, func() { BikeEventGen(dbc.DB, fakeBikefrequency) })

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
	users := []*database.User{}
	for userCount < maxUser {
		users = append(users, &database.User{
			Id:           gofakeit.UUID(),
			UserName:     gofakeit.Username(),
			EmailAddress: gofakeit.Email(),
		})

		userCount++
	}
	database.UpdateUser(users, db)
}
