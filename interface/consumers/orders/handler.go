package orders

import (
	"context"
	"fmt"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/domain/usecase"
	"github.com/rianekacahya/go-kafka/pkg/gokafka"
)

type event struct {
	ordersUsecase usecase.OrdersUsecase
}

func NewHandler(ctx context.Context, kafkago *gokafka.Gokafka, ordersUsecase usecase.OrdersUsecase) {
	transport := event{ordersUsecase}

	// start order consumer
	go kafkago.Consumer(ctx, &gokafka.ConsumerConfig{
		Topic:           entity.TopicOrders,
		GroupID:         entity.OrdersConsumerGroups,
		AutoOffsetReset: gokafka.OffsetEarliest,
		AutoRebalance:   true,
	}, transport.OrderConsumers)

	go kafkago.Consumer(ctx, &gokafka.ConsumerConfig{
		Topic:           "kuda",
		GroupID:         "kuda",
		AutoRebalance:   true,
		AutoOffsetReset: gokafka.OffsetEarliest,
	}, transport.OrderConsumers)
}

func (e *event) OrderConsumers(ctx context.Context, reader *gokafka.Reader) error {
	fmt.Println(string(reader.Message.Value))
	_, err := reader.Consumer.Commit()
	return err
}
