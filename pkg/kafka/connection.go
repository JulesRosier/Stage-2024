package kafka

import (
	"log/slog"
	"os"
	"stage2024/pkg/helper"

	_ "github.com/joho/godotenv/autoload"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

func GetClient() *kgo.Client {
	seed := os.Getenv("SEED_BROKER")

	slog.Info("Starting kafka client", "seedbrokers", seed)
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seed),
		kgo.ConsumeRegex(),
		kgo.ConsumeTopics("^[A-Za-z].*"),
		// kgo.ConsumerGroup("Testing"),
	)
	helper.MaybeDie(err, "error while starting kafka client")

	return cl
}

func GetRepoClient() *sr.Client {
	registry := os.Getenv("REGISTRY")

	slog.Info("starting schema registry client", "host", registry)
	rcl, err := sr.NewClient(sr.URLs(registry))
	helper.MaybeDieErr(err)
	return rcl
}
