package kafka

import (
	"context"
	"log/slog"
	"os"
	h "stage2024/pkg/helper"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

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

func ConnectSchemaRegistry() *sr.Client {
	registry := os.Getenv("REGISTRY")

	slog.Info("Starting schema registry client...", "host", registry)
	rcl, err := sr.NewClient(sr.URLs(registry))
	h.MaybeDieErr(err)

	return rcl
}

func GetSchema(topic string, rcl *sr.Client, file []byte) sr.SubjectSchema {
	sub := topic + "-value"
	ss, err := rcl.CreateSchema(context.Background(), sub, sr.Schema{
		Schema:     string(file),
		Type:       sr.TypeProtobuf,
		References: []sr.SchemaReference{h.ReferenceLocation(rcl)},
	})
	h.MaybeDieErr(err)
	slog.Info("Created or reusing schema", "sub", sub)

	return ss
}
