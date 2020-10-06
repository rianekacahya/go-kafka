package event

import (
	"context"
	"fmt"
	newrelic "github.com/newrelic/go-agent"
	"github.com/rianekacahya/go-kafka/pkg/gokafka"
	"github.com/rianekacahya/go-kafka/pkg/telemetry"
)

func Telemetry(app newrelic.Application) gokafka.MiddlewareFunc {
	return func(next gokafka.HandlerFunc) gokafka.HandlerFunc {
		return func(ctx context.Context, reader *gokafka.Reader) error {
			txnName := fmt.Sprintf("%s/%s", "event", *reader.Message.TopicPartition.Topic)
			txn := app.StartTransaction(txnName, nil, nil)
			defer txn.End()

			ctx = context.WithValue(ctx, telemetry.ContextKeyTelemetry, txn)
			if err := next(ctx, reader); err != nil {
				return err
			}

			return nil
		}
	}
}
