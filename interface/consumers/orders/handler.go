package orders

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	newrelic "github.com/newrelic/go-agent"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/domain/usecase"
	"github.com/rianekacahya/go-kafka/pkg/crashy"
	"github.com/rianekacahya/go-kafka/pkg/gokafka"
	"github.com/rianekacahya/go-kafka/pkg/telemetry"
)

type event struct {
	ordersUsecase usecase.OrdersUsecase
}

func NewHandler(ctx context.Context, kafkago *gokafka.Gokafka, ordersUsecase usecase.OrdersUsecase) {
	transport := event{ordersUsecase}

	go kafkago.Consumer(ctx, []string{entity.TopicOrders}, kafka.ConfigMap{
		"group.id":                        entity.OrdersConsumerGroups,
		"go.application.rebalance.enable": true,
		"enable.partition.eof":            false,
		"enable.auto.commit":              false,
		"auto.offset.reset":               gokafka.OffsetEarliest,
	}, transport.OrderConsumers)
}

func (e *event) OrderConsumers(ctx context.Context, reader *gokafka.Reader) error {
	var (
		req = new(entity.Orders)
		ntx = telemetry.ContextTelemetry(ctx)
	)

	defer reader.Consumer.Commit()

	// serialize message from kafka
	if err := json.Unmarshal(reader.Message.Value, req); err != nil {
		return crashy.Wrap(err, crashy.ErrCodeFormatting, "Failed when serialize data")
	}

	// newrelic submit order segment
	s := newrelic.StartSegment(ntx, "Hanlder.Event.OrderConsumers")

	if err := e.ordersUsecase.SaveOrders(ctx, req); err != nil {
		return err
	}

	if err := s.End(); err != nil {
		return crashy.Wrap(err, crashy.ErrCodeSend, "ending newrelic segment failed")
	}

	return nil
}
