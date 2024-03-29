package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"stage2024/pkg/helper"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

type Client struct {
	Kcl    *kgo.Client
	Rcl    *sr.Client
	Serde  *sr.Serde
	Config Config
}
type Config struct {
	Topics []Topic
}

func NewClient(config Config) *Client {
	rcl := getRepoClient()
	c := Client{
		Kcl:    getClient(),
		Rcl:    rcl,
		Serde:  getSerde(rcl, config.Topics),
		Config: config,
	}
	return &c
}

func (c Client) Produce(item any) error {
	topic, ok := c.findTopicByType(item)
	if !ok {
		return fmt.Errorf("no topic found for type %s", reflect.TypeOf(item))
	}
	itemBytes, err := c.Serde.Encode(item)
	if err != nil {
		return err
	}
	record := &kgo.Record{Topic: topic, Value: itemBytes, Key: []byte(strings.SplitN(topic, "_", 2)[0])}
	c.Kcl.Produce(context.Background(), record, func(r *kgo.Record, err error) {
		if err != nil {
			slog.Warn("Produce failed", "error", err)
		}
		slog.Debug("Produced event", "topic", topic, "offset", r.Offset)
	})
	return nil
}

func (c Client) findTopicByType(inputType any) (string, bool) {
	// FIXME: could be contant time with hashmap
	for _, topic := range c.Config.Topics {
		if reflect.TypeOf(inputType) == reflect.TypeOf(topic.PType) {
			return topic.getName(""), true
		}
	}
	return "", false
}

func getClient() *kgo.Client {
	seed := os.Getenv("SEED_BROKER")

	slog.Info("Starting kafka client", "seedbrokers", seed)
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seed),
		kgo.AllowAutoTopicCreation(),
		// we only need to produce
		// kgo.ConsumeRegex(),
		// kgo.ConsumeTopics("^[A-Za-z].*"),
		// kgo.ConsumerGroup("Testing"),
	)
	helper.MaybeDie(err, "error while starting kafka client")

	return cl
}

func getRepoClient() *sr.Client {
	registry := os.Getenv("REGISTRY")

	slog.Info("starting schema registry client", "host", registry)
	rcl, err := sr.NewClient(sr.URLs(registry))
	helper.MaybeDieErr(err)
	return rcl
}