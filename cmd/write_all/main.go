package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	writebaqme "stage2024/cmd/write_baqme"
	writebluebike "stage2024/cmd/write_bluebike"
	writebolt "stage2024/cmd/write_bolt"
	writecountpoles "stage2024/cmd/write_countpoles"
	writedonkey "stage2024/cmd/write_donkey"
	writestallinggent "stage2024/cmd/write_stalling_gent"
	writestallingstadskantoor "stage2024/cmd/write_stalling_stadskantoor"
	"stage2024/pkg/kafka"
	"syscall"
)

// writes all the data to kafka
func main() {
	slog.SetDefault(slog.New(slog.Default().Handler()))

	// start Kafka client
	cl := kafka.Connect()
	defer cl.Close()

	// start schema registry client and initialize all schemas
	rcl := kafka.ConnectSchemaRegistry()
	serde := kafka.GetSerde(rcl)

	// start all the writers
	go writebaqme.WriteBaqme(cl, serde)
	go writebluebike.WriteBluebike(cl, serde)
	go writebolt.WriteBolt(cl, serde)
	go writecountpoles.WriteCountples(cl, serde)
	go writedonkey.WriteDonkey(cl, serde)
	go writestallinggent.WriteStallingGent(cl, serde)
	go writestallingstadskantoor.WriteStallingStadskantoor(cl, serde)

	// wait for interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	fmt.Println("\nReceived an interrupt signal, exiting...")
	os.Exit(0)
}
