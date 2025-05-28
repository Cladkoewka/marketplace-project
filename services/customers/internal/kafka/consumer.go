package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/Cladkoewka/marketplace-project/services/customers/internal/kafka/event"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(broker, topic, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		GroupID: groupID,
		Topic:   topic,
	})
	return &Consumer{reader: reader}
}

func (c *Consumer) Consume(ctx context.Context, handler func(context.Context, event.OrderPlacedEvent) error) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var e event.OrderPlacedEvent
		if err := json.Unmarshal(m.Value, &e); err != nil {
			slog.Error("failed to unmarshal message", "error", err)
			continue
		}

		if err := handler(ctx, e); err != nil {
			slog.Error("failed to handle event", "error", err)
		} else {
			slog.Info("event handled", "event", e)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

