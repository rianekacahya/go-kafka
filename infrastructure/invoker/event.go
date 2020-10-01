package invoker

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/helper"
	"github.com/rianekacahya/go-kafka/pkg/crashy"
	"github.com/rianekacahya/go-kafka/pkg/gokafka"
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

	// tes publish message
	err = e.kafkago.Publish(ctx, &gokafka.ProducerConfig{
		Topic: entity.TopicOrders,
		Value: body,
		Headers: []kafka.Header{
			{
				Key:   gokafka.HeadersRequestID,
				Value: []byte(helper.GetContextString(ctx, entity.ContextRequestID)),
			},
		},
	})

	if err != nil {
		return crashy.Wrap(err, crashy.ErrCodeSend, "sending kafka message failed")
	}

	payload.CustomerEmail.Scan("saya@kuda.com")
	body1, err := json.Marshal(payload)
	if err != nil {
		return crashy.Wrap(err, crashy.ErrCodeFormatting, "error when encode request body to json")
	}

	// tes publish message
	err = e.kafkago.Publish(ctx, &gokafka.ProducerConfig{
		Topic: "kuda",
		Value: body1,
		Headers: []kafka.Header{
			{
				Key:   gokafka.HeadersRequestID,
				Value: []byte(helper.GetContextString(ctx, entity.ContextRequestID)),
			},
		},
	})

	if err != nil {
		return crashy.Wrap(err, crashy.ErrCodeSend, "sending kafka message failed")
	}

	return nil
}
