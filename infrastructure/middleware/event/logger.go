package event

import (
	"context"
	"encoding/json"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/pkg/gokafka"
	"go.uber.org/zap"
)

func Logger(log *zap.Logger) gokafka.MiddlewareFunc {
	return func (next gokafka.HandlerFunc) gokafka.HandlerFunc {
		return func(ctx context.Context, reader *gokafka.Reader) error {
			log.Info("Kafka Consumer Logger",
				zap.String("topic", *reader.Message.TopicPartition.Topic),
				zap.Any("key", json.RawMessage(reader.Message.Key)),
				zap.Any("value", json.RawMessage(reader.Message.Value)),
				zap.String("request_id", gokafka.FindHeaders(entity.ContextRequestID, reader.Message.Headers)),
			)

			if err := next(ctx, reader); err != nil {
				log.Info("Kafka Consumer Error Logger",
					zap.String("topic", *reader.Message.TopicPartition.Topic),
					zap.Any("key", json.RawMessage(reader.Message.Key)),
					zap.Any("value", json.RawMessage(reader.Message.Value)),
					zap.String("request_id", gokafka.FindHeaders(entity.ContextRequestID, reader.Message.Headers)),
					zap.Any("message", err),
				)

				return err
			}

			return nil
		}
	}
}