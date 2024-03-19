package main

import (
	"fmt"
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

func main() {
	cl := kafka.Connect()
	defer cl.Close()
	rcl := kafka.ConnectSchemaRegistry()

	go writebaqme.WriteBaqme(cl, rcl)
	go writebluebike.WriteBluebike(cl, rcl)
	go writebolt.WriteBolt(cl, rcl)
	go writecountpoles.WriteCountples(cl, rcl)
	go writedonkey.WriteDonkey(cl, rcl)
	go writestallinggent.WriteStallingGent(cl, rcl)
	go writestallingstadskantoor.WriteStallingStadskantoor(cl, rcl)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	fmt.Println("\nReceived an interrupt signal, exiting...")
	os.Exit(0)
}
