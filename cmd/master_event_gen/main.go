package main

import (
	"log/slog"
	"os"
	"os/signal"
	"stage2024/pkg/database"
	"stage2024/pkg/opendata"
	"stage2024/pkg/scheduler"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/gorm"
)

const maxUser = 100

func main() {
	database.Init()

	db := database.GetDb()

	CreateUsers(db)

	changesCh := make(chan []string, 100)

	s := scheduler.NewScheduler()

	s.Schedule(time.Minute*5, func() { opendata.Bolt(db, changesCh) })

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

	slog.Info("Exiting... Goodbye!🧌")
}

func CreateUsers(db *gorm.DB) {
	var userCount int64
	db.Model(&database.User{}).Count(&userCount)

	for userCount < maxUser {
		CreateRandomUser(db)
		userCount++
	}
}

func CreateRandomUser(db *gorm.DB) {
	u := &database.User{
		Id:           gofakeit.UUID(),
		UserName:     gofakeit.Username(),
		EmailAddress: gofakeit.Email(),
	}
	db.Create(u)
}
