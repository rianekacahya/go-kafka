package usecase

import (
	"context"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/domain/invoker"
	"github.com/rianekacahya/go-kafka/domain/repository"
	"github.com/rianekacahya/go-kafka/helper"
	"github.com/rianekacahya/go-kafka/pkg/crashy"
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
	if err := req.PurchaseDate.Scan(time.Now()); err != nil {
		return crashy.Wrap(err, crashy.ErrCodeFormatting, "Errror when binding purchase date")
	}
	return o.eventInvoker.OrderProducers(ctx, req)
}

func (o *orders) SaveOrders(ctx context.Context, req *entity.Orders) error {
	// set purchase code
	if err := req.PurchaseCode.Scan(helper.RandomString(20)); err != nil {
		return crashy.Wrap(err, crashy.ErrCodeFormatting, "Errror when binding purchase code")
	}

	if err := o.ordersRepository.SaveOrder(ctx, req); err != nil {
		return err
	}

	return nil
}