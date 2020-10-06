package invoker

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/pkg/crashy"
	"github.com/rianekacahya/go-kafka/pkg/gokafka"
	"github.com/rianekacahya/go-kafka/pkg/helper"
)

type event struct {
	kafkago *gokafka.Gokafka
}

func EventInitialize(kafkago *gokafka.Gokafka) *event {
	return &event{kafkago}
}

func (e *event) OrderProducers(ctx context.Context, payload *entity.Orders) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return crashy.Wrap(err, crashy.ErrCodeFormatting, "error when encode request body to json")
	}

	err = e.kafkago.Publish(ctx, &gokafka.Payload{
		Topic: entity.TopicOrders,
		Value: body,
		Headers: []kafka.Header{
			{
				Key:   entity.ContextRequestID,
				Value: []byte(helper.GetContextString(ctx, entity.ContextRequestID)),
			},
		},
	}, kafka.ConfigMap{})

	return nil
}
