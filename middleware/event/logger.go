package event

import (
	"context"
	"encoding/json"
	"github.com/rianekacahya/go-kafka/pkg/gokafka"
	"github.com/rianekacahya/go-kafka/pkg/logger"
	"go.uber.org/zap"
)

func Logger() gokafka.MiddlewareFunc {
	return func (next gokafka.HandlerFunc) gokafka.HandlerFunc {
		return func(ctx context.Context, reader *gokafka.Reader) error {
			logger.GetLogger().Info("Kafka Consumer Logger",
				zap.String("topic", *reader.Message.TopicPartition.Topic),
				zap.Any("key", json.RawMessage(reader.Message.Key)),
				zap.Any("value", json.RawMessage(reader.Message.Value)),
				zap.String("request_id", gokafka.FindHeaders("request_id", reader.Message.Headers)),
			)

			if err := next(ctx, reader); err != nil {
				logger.GetLogger().Info("Kafka Consumer Error Logger",
					zap.String("topic", *reader.Message.TopicPartition.Topic),
					zap.Any("key", json.RawMessage(reader.Message.Key)),
					zap.Any("value", json.RawMessage(reader.Message.Value)),
					zap.String("request_id", gokafka.FindHeaders("request_id", reader.Message.Headers)),
					zap.Any("message", err),
				)
			}

			return nil
		}
	}
}