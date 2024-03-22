package helper

import (
	"context"
	"sync"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
)

func Produce(serde *sr.Serde, cl *kgo.Client, wg *sync.WaitGroup, item any, ctx context.Context, topic string) {
	itemByte, err := serde.Encode(item)
	MaybeDie(err, "Encoding error")
	record := &kgo.Record{Topic: topic, Value: itemByte}
	cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
		defer wg.Done()
		MaybeDie(err, "Produce error")
	})
}
