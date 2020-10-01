package usecase

import (
	"context"
	"github.com/rianekacahya/go-kafka/domain/entity"
)

type OrdersUsecase interface {
	SubmitOrders(context.Context, *entity.Orders) error
}