package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"stage2024/pkg/helper"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

type KafkaClient struct {
	Kcl    *kgo.Client
	Rcl    *sr.Client
	Serde  *sr.Serde
	Config Config
}

type Config struct {
	Topics []Topic
}

func NewClient(config Config) *KafkaClient {
	rcl := getRepoClient()
	c := KafkaClient{
		Kcl:    getClient(),
		Rcl:    rcl,
		Serde:  getSerde(rcl, config.Topics),
		Config: config,
	}
	return &c
}

func (c KafkaClient) Produce(item any, timeStamp time.Time) error {
	topic, ok := c.findTopicByType(item)
	if !ok {
		return fmt.Errorf("no topic found for type %s", reflect.TypeOf(item))
	}
	itemBytes, err := c.Serde.Encode(item)
	if err != nil {
		return err
	}
	record := &kgo.Record{
		Topic:     topic,
		Value:     itemBytes,
		Key:       []byte(strings.SplitN(topic, "_", 2)[0]),
		Timestamp: timeStamp,
	}
	rs := c.Kcl.ProduceSync(context.Background(), record)
	for _, r := range rs {
		if r.Err != nil {
			return r.Err
		}
	}
	return nil
}

func (c KafkaClient) findTopicByType(inputType any) (string, bool) {
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
