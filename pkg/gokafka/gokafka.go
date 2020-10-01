package gokafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rianekacahya/go-kafka/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	KafkaOffset string

	Gokafka struct {
		AddressFamily string
		Address       string
		MaxRetry      int
	}

	HandlerFunc func(context.Context, Reader) error

	Reader struct {
		Gokafka *Gokafka
		Message *kafka.Message
	}

	ConsumerConfig struct {
		Topic           string
		GroupID         string
		SessionTimeout  int
		AutoRebalance   bool
		AutoCommit      bool
		PartitionEOF    bool
		AutoOffsetReset KafkaOffset
	}

	ProducerConfig struct {
		Idempotence bool
		Headers     []kafka.Header
		Topic       string
		Value       []byte
	}
)

const (
	OffsetEarliest KafkaOffset = "earliest"
	OffsetLatest   KafkaOffset = "latest"
	OffsetNone     KafkaOffset = "none"
)

const (
	HeadersRequestID     = "request_id"
)

func New(address, addressfamily string, maxretry int) *Gokafka {
	return &Gokafka{address, addressfamily, maxretry}
}

func (k *Gokafka) Consumer(ctx context.Context, config *ConsumerConfig, action HandlerFunc) error {
	conf, err := config.configurating(k)
	if err != nil {
		return err
	}

	c, err := kafka.NewConsumer(&conf)
	if err != nil {
		return err
	}

	if err := c.SubscribeTopics([]string{config.Topic}, nil); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case ev := <-c.Events():
			switch e := ev.(type) {
			case *kafka.Message:
				err := action(ctx, Reader{k, e})
				k.error(err, e)
			case kafka.Error:
				k.error(e, nil)
			case kafka.PartitionEOF:
			}
		}
	}
}

func (cc *ConsumerConfig) configurating(gokafka *Gokafka) (kafka.ConfigMap, error) {
	var conf = make(kafka.ConfigMap)

	if cc.Topic == "" {
		return nil, errors.New("topic can't be null")
	}

	if err := conf.SetKey("bootstrap.servers", gokafka.Address); err != nil {
		return nil, err
	}

	if err := conf.SetKey("broker.address.family", gokafka.AddressFamily); err != nil {
		return nil, err
	}

	if err := conf.SetKey("enable.auto.commit", cc.AutoCommit); err != nil {
		return nil, err
	}

	if err := conf.SetKey("enable.partition.eof", cc.PartitionEOF); err != nil {
		return nil, err
	}

	if err := conf.SetKey("go.application.rebalance.enable", cc.AutoRebalance); err != nil {
		return nil, err
	}

	if err := conf.SetKey("group.id", cc.GroupID); err != nil {
		return nil, err
	}

	if cc.SessionTimeout != 0 {
		if err := conf.SetKey("session.timeout.ms", cc.SessionTimeout); err != nil {
			return nil, err
		}
	}

	if cc.AutoOffsetReset != "" {
		if err := conf.SetKey("auto.offset.reset", cc.AutoOffsetReset); err != nil {
			return nil, err
		}
	} else {
		if err := conf.SetKey("auto.offset.reset", OffsetNone); err != nil {
			return nil, err
		}
	}

	return conf, nil
}

func (k *Gokafka) Publish(ctx context.Context, config *ProducerConfig) error {
	var conf = make(kafka.ConfigMap)

	if err := conf.SetKey("bootstrap.servers", k.Address); err != nil {
		return err
	}

	if err := conf.SetKey("broker.address.family", k.AddressFamily); err != nil {
		return err
	}

	if err := conf.SetKey("enable.idempotence", config.Idempotence); err != nil {
		return err
	}

	p, err := kafka.NewProducer(&conf)
	if err != nil {
		return err
	}

	delivery := make(chan kafka.Event)
	defer close(delivery)

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &config.Topic, Partition: kafka.PartitionAny},
		Value:          config.Value,
		Headers:        config.Headers,
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
			if v.Key == HeadersRequestID {
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
