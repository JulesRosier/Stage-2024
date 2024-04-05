package events

import (
	"stage2024/pkg/helper"
	"stage2024/pkg/kafka"
)

type EventClient struct {
	Kc      *kafka.Client
	Channel chan helper.Change
}

func NewEventClient(kc *kafka.Client) *EventClient {

	return &EventClient{
		Kc:      kc,
		Channel: make(chan helper.Change, 100),
	}
}
