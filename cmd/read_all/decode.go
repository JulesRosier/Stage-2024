package main

import (
	"context"
	"errors"
	"fmt"
	"stage2024/pkg/serde"

	"github.com/twmb/franz-go/pkg/sr"
)

func handleDecode(ctx context.Context, cl *sr.Client, record []byte, serdeCache map[int]*serde.Serde) ([]byte, error) {
	var serdeHeader sr.ConfluentHeader
	id, toDecode, err := serdeHeader.DecodeID(record)
	if err != nil {
		return nil, errors.New("unable to decode the ID")
	}
	var rpkSerde *serde.Serde
	if cachedSerde, ok := serdeCache[id]; ok {
		rpkSerde = cachedSerde
	} else {
		schema, err := cl.SchemaByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("unable to get schema with index %v from the schema registry: %v", id, err)
		}
		rpkSerde, err = serde.NewSerde(ctx, cl, &schema, id, "")
		if err != nil {
			return nil, err
		}
		serdeCache[id] = rpkSerde
	}

	return rpkSerde.DecodeRecord(toDecode)
}
