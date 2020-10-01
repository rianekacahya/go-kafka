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

func NewHandler(ctx context.Context, kafkago *gokafka.Gokafka, ordersUsecase usecase.OrdersUsecase) error {
	var (
		err       error
		transport = event{ordersUsecase}
	)

	// start order consumer
	if err = kafkago.Consumer(ctx, &gokafka.ConsumerConfig{
		Topic:           entity.TopicOrders,
		GroupID:         entity.OrdersConsumerGroups,
		SessionTimeout:  6000,
		AutoOffsetReset: gokafka.OffsetEarliest,
	}, transport.OrderConsumers); err != nil {
		return err
	}

	if err = kafkago.Consumer(ctx, &gokafka.ConsumerConfig{
		Topic:           "kuda",
		GroupID:         "kuda",
		SessionTimeout:  6000,
		AutoOffsetReset: gokafka.OffsetEarliest,
	}, transport.OrderConsumers); err != nil {
		return err
	}

	return nil
}

func (e *event) OrderConsumers(ctx context.Context, reader *gokafka.Reader) error {
	fmt.Println(string(reader.Message.Value))
	_, err := reader.Consumer.Commit()
	return err
}
