package kafka

import (
	"log/slog"
	"stage2024/pkg/database"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"gorm.io/gorm"
)

type OutboxListener struct {
	kc       *KafkaClient
	db       *database.DatabaseClient
	topicMap map[string]protoreflect.ProtoMessage
	queue    chan database.Outbox
}

func NewOutboxListener(kc *KafkaClient, db *database.DatabaseClient, topics []Topic) *OutboxListener {
	tm := map[string]protoreflect.ProtoMessage{}
	for _, topic := range topics {
		tm[topic.getName("")] = topic.PType
	}

	return &OutboxListener{
		kc:       kc,
		db:       db,
		topicMap: tm,
		queue:    make(chan database.Outbox, 1000),
	}
}

func (ol *OutboxListener) Start() {
	go ol.listen()
}

func (ol *OutboxListener) listen() {
	for o := range ol.queue {
		err := ol.db.DB.Transaction(func(tx *gorm.DB) error {
			err := tx.Delete(&o).Error
			if err != nil {
				return err
			}
			t := ol.topicMap[o.Topic]
			err = proto.Unmarshal(o.Payload, t)
			if err != nil {
				return err
			}
			err = ol.kc.Produce(t)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			slog.Warn("Outbox listener transaction failed", "error", err)
		}
	}
}

// Schould be called on a schedule.
// Pulls outbox record out of database and puts them on the consume queue.
func (ol *OutboxListener) FetchOutboxData() {
	rows := []database.Outbox{}
	result := ol.db.DB.Order("created_at desc").Find(&rows)
	if result.Error != nil {
		slog.Warn("Failed to fetch Outbox data", "error", result.Error)
	}
	slog.Debug("Fetched all outbox records", "count", result.RowsAffected)

	for _, row := range rows {
		ol.queue <- row
	}
}
