package usecase

import (
	"context"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/domain/invoker"
	"github.com/rianekacahya/go-kafka/domain/repository"
	"github.com/rianekacahya/go-kafka/pkg/helper"
	"time"
)

type orders struct {
	ordersRepository repository.OrdersRepository
	eventInvoker     invoker.EventInvoker
}

func NewOrdersUsecase(ordersRepository repository.OrdersRepository, eventInvoker invoker.EventInvoker) *orders {
	return &orders{ordersRepository, eventInvoker}
}

func (o *orders) SubmitOrders(ctx context.Context, req *entity.Orders) error {
	// set purchase date
	req.PurchaseDate = entity.TimeNull(time.Now())

	return o.eventInvoker.OrderProducers(ctx, req)
}

func (o *orders) SaveOrders(ctx context.Context, req *entity.Orders) error {
	// set purchase code
	req.PurchaseCode = entity.StrNull(helper.RandomString(20))

	return o.ordersRepository.SaveOrder(ctx, req)
}