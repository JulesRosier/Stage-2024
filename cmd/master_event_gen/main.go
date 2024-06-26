package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"stage2024/pkg/database"
	"stage2024/pkg/events"
	"stage2024/pkg/helper"
	"stage2024/pkg/kafka"
	"stage2024/pkg/opendata"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/stations"
	"stage2024/pkg/protogen/users"
	"stage2024/pkg/scheduler"
	"stage2024/pkg/settings"
	"time"
)

const fakeBikeFrequency = 60
const maxUsers = 500
const maxBikes = 500

const banner = `
███████╗██╗   ██╗███████╗███╗   ██╗████████╗     ██████╗ ███████╗███╗   ██╗
██╔════╝██║   ██║██╔════╝████╗  ██║╚══██╔══╝    ██╔════╝ ██╔════╝████╗  ██║
█████╗  ██║   ██║█████╗  ██╔██╗ ██║   ██║       ██║  ███╗█████╗  ██╔██╗ ██║
██╔══╝  ╚██╗ ██╔╝██╔══╝  ██║╚██╗██║   ██║       ██║   ██║██╔══╝  ██║╚██╗██║
███████╗ ╚████╔╝ ███████╗██║ ╚████║   ██║       ╚██████╔╝███████╗██║ ╚████║
╚══════╝  ╚═══╝  ╚══════╝╚═╝  ╚═══╝   ╚═╝        ╚═════╝ ╚══════╝╚═╝  ╚═══╝
→ https://github.com/JulesRosier/Stage-2024
███████████████████████████████████████████████████████████████████████████╗
╚══════════════════════════════════════════════════════════════════════════╝`

func main() {
	fmt.Println(banner)
	set, err := settings.Load()
	helper.MaybeDie(err, "Failed to load configs")

	logLevel := helper.GetLogLevel(set.Logger)
	fmt.Printf("LOG_LEVEL = %s\n", logLevel)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	dbc := database.NewDatabase(set.Database)

	topics := []kafka.ProtoDefinition{
		{ProtoFile: users.File_users_user_registered_proto, PType: &users.UserRegistered{}},

		{ProtoFile: bikes.File_bikes_bike_abandoned_proto, PType: &bikes.BikeAbandoned{}},
		{ProtoFile: bikes.File_bikes_bike_defect_reported_proto, PType: &bikes.BikeDefectReported{}},
		{ProtoFile: bikes.File_bikes_bike_deployed_proto, PType: &bikes.BikeDeployed{}},
		{ProtoFile: bikes.File_bikes_bike_immobilized_proto, PType: &bikes.BikeImmobilized{}},
		{ProtoFile: bikes.File_bikes_bike_picked_up_proto, PType: &bikes.BikePickedUp{}},
		{ProtoFile: bikes.File_bikes_bike_reserved_proto, PType: &bikes.BikeReserved{}},
		{ProtoFile: bikes.File_bikes_bike_returned_proto, PType: &bikes.BikeReturned{}},
		{ProtoFile: bikes.File_bikes_bike_stored_proto, PType: &bikes.BikeStored{}},
		{ProtoFile: bikes.File_bikes_bike_repaired_proto, PType: &bikes.BikeRepaired{}},

		{ProtoFile: stations.File_stations_station_occupation_increased_proto, PType: &stations.StationOccupationDecreased{}},
		{ProtoFile: stations.File_stations_station_full_proto, PType: &stations.StationFull{}},
		{ProtoFile: stations.File_stations_station_occupation_decreased_proto, PType: &stations.StationOccupationIncreased{}},
		{ProtoFile: stations.File_stations_station_created_proto, PType: &stations.StationCreated{}},
		{ProtoFile: stations.File_stations_station_deprecated_proto, PType: &stations.StationDeprecated{}},
	}

	kc := kafka.NewClient(kafka.Config{
		ProtoDefinition: topics,
		Settings:        set.Kafka,
	})
	kc.CreateTopics(context.Background())

	ol := kafka.NewOutboxListener(kc, dbc, topics)
	ol.Start()

	s := scheduler.NewScheduler()

	database.CreateUsersBikes(dbc.DB, maxUsers, maxBikes)

	// Schedule the data fetching functions
	s.Schedule(time.Minute*5, func() { opendata.BlueBike(dbc.DB) })
	s.Schedule(time.Minute*10, func() { opendata.Donkey(dbc.DB) })
	s.Schedule(time.Minute*1, func() { opendata.StorageGhent(dbc.DB) })
	s.Schedule(time.Minute*5, func() { opendata.StorageTownHall(dbc.DB) })

	s.Schedule(time.Second*30, ol.FetchOutboxData)

	s.Schedule(time.Hour*12, func() { database.BikeCleanUpTransaction(dbc.DB) })

	s.Schedule(time.Minute*fakeBikeFrequency, func() { events.BikeEventGen(dbc.DB) })

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("Received an interrupt signal, exiting...")

	s.Stop()

	slog.Info("Exiting... Goodbye!")
}
