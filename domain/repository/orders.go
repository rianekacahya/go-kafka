package repository

import (
	"context"
	"github.com/rianekacahya/go-kafka/domain/entity"
)

type OrdersRepository interface {
	SaveOrder(context.Context, *entity.Orders) error
}
