package usecase

import (
	"context"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/domain/invoker"
	"github.com/rianekacahya/go-kafka/domain/repository"
)

type orders struct {
	ordersRepository repository.OrdersRepository
	eventInvoker     invoker.EventInvoker
}

func NewOrdersUsecase(ordersRepository repository.OrdersRepository, eventInvoker invoker.EventInvoker) *orders {
	return &orders{ordersRepository, eventInvoker}
}

func (o *orders) SubmitOrders(ctx context.Context, req *entity.Orders) error {
	return o.eventInvoker.OrderProducers(ctx, req)
}
