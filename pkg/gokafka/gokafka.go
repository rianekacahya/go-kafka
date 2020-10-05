package gokafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type (
	Gokafka struct {
		Address    string
		Middleware []MiddlewareFunc
	}

	HandlerFunc func(context.Context, *Reader) error
	MiddlewareFunc func(HandlerFunc) HandlerFunc

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

func New(address string, middleware ...MiddlewareFunc) *Gokafka {
	return &Gokafka{address, middleware}
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
				log.Println("kafka have new consumer")
				if err := c.Assign(e.Partitions); err != nil {
					log.Printf("got an error while assign new consumer, error: %v\n", err)
				}
			case kafka.RevokedPartitions:
				log.Println("kafka have revoked consumer")
				if err := c.Unassign(); err != nil {
					log.Printf("got an error while revoke consumer, error: %v\n", err)
				}
			case *kafka.Message:
				h := action
				for i := len(k.Middleware) - 1; i >= 0; i-- {
					h = k.Middleware[i](h)
				}
				if err := h(ctx, &Reader{k, c, e}); err != nil {
					log.Printf("got an error while processing message, error: %v\n", err)
				}
			case kafka.Error:
				log.Printf("got an error while processing event, error: %v\n", e)
			default:
				log.Printf("Unhandled event, ignored: %v\n", e)
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
