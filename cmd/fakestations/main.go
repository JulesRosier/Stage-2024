package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
	"stage2024/pkg/scheduler"
	"stage2024/pkg/settings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Starting...")
	set, err := settings.Load()
	helper.MaybeDie(err, "Failed to load configs")

	logLevel := helper.GetLogLevel(set.Logger)
	fmt.Printf("LOG_LEVEL = %s\n", logLevel)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	dbc := database.NewDatabase(set.Database)

	s := scheduler.NewScheduler()
	station := &database.Station{
		Id:          gofakeit.UUID(),
		OpenDataId:  "fakest",
		Name:        "Fake station",
		Lat:         52.3676,
		Lon:         4.9041,
		MaxCapacity: 10,
		Occupation:  5,
		IsActive:    sql.NullBool{Bool: true, Valid: true},
	}
	s.Schedule(time.Second*1, func() { fakeStationData(dbc.DB, station) })

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("Received an interrupt signal, exiting...")

	s.Stop()

	slog.Info("Exiting... Goodbye!")
}

func fakeStationData(db *gorm.DB, station *database.Station) {
	// Create a station
	station.Occupation = int32(gofakeit.Number(0, int(station.MaxCapacity)))
	database.UpdateStation([]*database.Station{station}, db)

	slog.Info("updated", "station", station.Id, "occupation", station.Occupation)
}
