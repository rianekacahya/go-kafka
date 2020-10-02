package gokafka

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rianekacahya/go-kafka/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

type (
	Gokafka struct {
		Address string
	}

	HandlerFunc func(context.Context, *Reader) error

	Payload struct {
		Topic   string
		Value   []byte
		Key     []byte
		Headers []kafka.Header
	}

	Reader struct {
		Gokafka  *Gokafka
		Consumer *kafka.Consumer
		Message  *kafka.Message
	}
)

const (
	OffsetEarliest = "earliest"
	OffsetLatest   = "latest"
	OffsetNone     = "none"
)

func New(address string) *Gokafka {
	return &Gokafka{address}
}

func (k *Gokafka) Consumer(ctx context.Context, topic []string, config kafka.ConfigMap, action HandlerFunc) {
	config["bootstrap.servers"] = k.Address

	c, err := kafka.NewConsumer(&config)
	if err != nil {
		log.Fatalf("got an error while bootsraping consumer, error: %s", err)
	}

	if err := c.SubscribeTopics(topic, nil); err != nil {
		log.Fatalf("got an error while bootsraping consumer error: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				if err := c.Assign(e.Partitions); err != nil {
					k.error(err, nil)
				}
			case kafka.RevokedPartitions:
				if err := c.Unassign(); err != nil {
					k.error(err, nil)
				}
			case *kafka.Message:
				if err := action(ctx, &Reader{k, c, e}); err != nil {
					k.error(err, e)
				}
			case kafka.Error:
				k.error(e, nil)
			default:
			}
		}
	}
}

func (k *Gokafka) Publish(ctx context.Context, payload *Payload, config kafka.ConfigMap) error {
	config["bootstrap.servers"] = k.Address

	p, err := kafka.NewProducer(&config)
	if err != nil {
		return err
	}

	delivery := make(chan kafka.Event)
	defer close(delivery)

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &payload.Topic, Partition: kafka.PartitionAny},
		Key:            payload.Key,
		Value:          payload.Value,
		Headers:        payload.Headers,
	}, delivery)
	if err != nil {
		return err
	}

	e := <-delivery
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}

func (k *Gokafka) error(err error, msg *kafka.Message) {
	fields := []zapcore.Field{
		zap.Any("error", err),
	}

	if msg != nil {
		var requestID string

		for _, v := range msg.Headers {
			if v.Key == "request_id" {
				requestID = string(v.Value)
			}

		}

		fields = append(fields,
			zap.String("topic", *msg.TopicPartition.Topic),
			zap.Any("value", json.RawMessage(msg.Value)),
			zap.String("request_id", requestID),
		)
	}

	logger.GetLogger().Error("Kafka Error Logger", fields...)
}
