package kafka

import (
	"log/slog"
	"os"
	h "stage2024/pkg/helper"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

// Creates a new Kafka client
func Connect() *kgo.Client {
	seed := os.Getenv("SEED_BROKER")

	// seed := flag.String("seedbroker", "localhost:19092", "brokers port to talk to")
	// registry := flag.String("registry", "localhost:18081", "schema registry port to talk to")
	// topic := flag.String("topic", topicname, "topic to produce to and consume from")

	slog.Info("Starting kafka client...", "seedbroker", seed)
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seed),
		kgo.AllowAutoTopicCreation(),
	)
	h.MaybeDieErr(err)

	return cl
}

// Creates a new schema registry client
func ConnectSchemaRegistry() *sr.Client {
	registry := os.Getenv("REGISTRY")

	slog.Info("Starting schema registry client...", "host", registry)
	rcl, err := sr.NewClient(sr.URLs(registry))
	h.MaybeDieErr(err)

	return rcl
}
