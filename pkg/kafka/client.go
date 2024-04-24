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
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kversion"
	"github.com/twmb/franz-go/pkg/sasl/scram"
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
	for _, topic := range c.Config.Topics {
		if reflect.TypeOf(inputType) == reflect.TypeOf(topic.PType) {
			return topic.getName(""), true
		}
	}
	return "", false
}

func (c *KafkaClient) CreateTopics(ctx context.Context) {
	acl := kadm.NewClient(c.Kcl)
	topicDetails, err := acl.ListTopics(ctx)
	helper.MaybeDieErr(err)

	for _, topic := range c.Config.Topics {
		t := topic.getName("")
		if !topicDetails.Has(t) {
			a, err := acl.CreateTopic(ctx, -1, -1, nil, t)
			helper.MaybeDie(err, "Failed to create topic")
			slog.Info("topic created", "topic", a.Topic)
		}
	}
}

func getClient() *kgo.Client {
	seed := os.Getenv("SEED_BROKER")
	user := os.Getenv("EH_AUTH_USER")
	pw := os.Getenv("EH_AUTH_PASSWORD")

	slog.Info("Starting kafka client", "seedbrokers", seed)
	clientConfigs := []kgo.Opt{
		kgo.SeedBrokers(seed),
		kgo.FetchMaxBytes(5 * 1000 * 1000),
		kgo.MaxConcurrentFetches(12),
		kgo.MaxVersions(kversion.V2_6_0()),
		kgo.AllowAutoTopicCreation(),
	}

	if user != "" && pw != "" {
		clientConfigs = append(clientConfigs, kgo.SASL(scram.Auth{
			User: user,
			Pass: pw,
		}.AsSha512Mechanism()))
	}
	cl, err := kgo.NewClient(
		clientConfigs...,
	)

	helper.MaybeDie(err, "error while starting kafka client")
	err = cl.Ping(context.Background())
	helper.MaybeDie(err, "No ping")
	return cl
}

func getRepoClient() *sr.Client {
	registry := os.Getenv("REGISTRY")
	user := os.Getenv("EH_AUTH_USER")
	pw := os.Getenv("EH_AUTH_PASSWORD")

	slog.Info("Starting schema registry client", "host", registry)
	opts := []sr.Opt{
		sr.URLs(registry),
	}
	if user != "" && pw != "" {
		opts = append(opts, sr.BasicAuth(user, pw))
	}
	rcl, err := sr.NewClient(opts...)
	helper.MaybeDieErr(err)
	return rcl
}
