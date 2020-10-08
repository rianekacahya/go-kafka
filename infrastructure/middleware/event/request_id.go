package event

import (
	"context"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/pkg/gokafka"
)

func RequestID() gokafka.MiddlewareFunc {
	return func(next gokafka.HandlerFunc) gokafka.HandlerFunc {
		return func(ctx context.Context, reader *gokafka.Reader) error {
			ctx = context.WithValue(ctx, entity.ContextRequestID, gokafka.FindHeaders(entity.ContextRequestID, reader.Message.Headers))
			return next(ctx, reader)
		}
	}
}
