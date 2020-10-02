package orders

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/domain/usecase"
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
	fmt.Println(string(reader.Message.Value))
	_, err := reader.Consumer.Commit()
	return err
}
