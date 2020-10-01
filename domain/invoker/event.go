package invoker

import (
	"context"
	"github.com/rianekacahya/go-kafka/domain/entity"
)

type EventInvoker interface {
	OrderProducers(context.Context, *entity.Orders) error
}
