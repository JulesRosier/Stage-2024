package main

import (
	"log/slog"
	"os"
	"os/signal"
	"stage2024/pkg/database"
	"stage2024/pkg/kafka"
	//"stage2024/pkg/server"
)

func main() {
	slog.SetDefault(slog.New(slog.Default().Handler()))

	//server := server.NewServer()
	//server.RegisterRoutes()
	//server.ApplyMiddleware()

	database.Init()

	// server.Start()
	go kafka.EventImporter()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("Received an interrupt signal, exiting...")

	// server.Stop()
}
