package orders

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/domain/usecase"
	"github.com/rianekacahya/go-kafka/pkg/crashy"
	"github.com/rianekacahya/go-kafka/pkg/gokafka"
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
	var req = new(entity.Orders)

	defer reader.Consumer.Commit()

	// serialize message from kafka
	if err := json.Unmarshal(reader.Message.Value, req); err != nil {
		return crashy.Wrap(err, crashy.ErrCodeFormatting, "Failed when serialize data")
	}

	if err := e.ordersUsecase.SaveOrders(ctx, req); err != nil {
		return err
	}

	return nil
}
