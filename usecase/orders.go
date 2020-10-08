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
	syncInvoker      invoker.SyncInvoker
}

func NewOrdersUsecase(ordersRepository repository.OrdersRepository, eventInvoker invoker.EventInvoker, syncInvoker invoker.SyncInvoker) *orders {
	return &orders{ordersRepository, eventInvoker, syncInvoker}
}

func (o *orders) SubmitOrders(ctx context.Context, req *entity.Orders) error {
	// set purchase date
	req.PurchaseDate = entity.TimeNull(time.Now())

	// publish event order
	if err := o.eventInvoker.OrderProducers(ctx, req); err != nil {
		return err
	}

	// start gossip between producer and event consumer
	if err := <- o.syncInvoker.EventListener(ctx, helper.GetContextString(ctx, entity.ContextRequestID)); err != nil {
		return err
	}

	return nil
}

func (o *orders) SaveOrders(ctx context.Context, req *entity.Orders) error {
	// set purchase code
	req.PurchaseCode = entity.StrNull(helper.RandomString(20))
	if err := o.ordersRepository.SaveOrder(ctx, req); err != nil {
		return err
	}

	// if everything is oke, send gossip to whisper
	if err := o.syncInvoker.EventWhisper(ctx, helper.GetContextString(ctx, entity.ContextRequestID)); err != nil {
		return err
	}

	return nil
}
