package main

import (
	"context"
	"fmt"
	"log/slog"
	"stage2024/pkg/database"
	"stage2024/pkg/helper"
	"stage2024/pkg/kafka"
	"stage2024/pkg/serde"

	"github.com/jackc/pgx/v5/pgtype"
)

func main() {
	slog.SetDefault(slog.New(slog.Default().Handler()))

	database.Init()
	temp()
}

func temp() {
	q := database.GetQueries()
	cl := kafka.GetClient()
	rcl := kafka.GetRepoClient()

	ctx := context.Background()
	c := make(map[int]*serde.Serde)
	for {
		fetches := cl.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			panic(fmt.Sprint(errs))
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			j, err := serde.HandleDecode(ctx, rcl, record.Value, c)
			helper.MaybeDieErr(err)
			//fmt.Println(string(j))
			q.CreateEvent(ctx, database.CreateEventParams{
				Data:      j,
				TopicName: pgtype.Text{String: record.Topic, Valid: true},
			})
			slog.Info("inserted record", "topic", record.Topic)
		}

	}
}
