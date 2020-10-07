package invoker

import (
	"context"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/stretchr/testify/mock"
)

type Event struct {
	mock.Mock
}

func (e *Event) OrderProducers(ctx context.Context, payload *entity.Orders) error {
	ret := e.Called(ctx, payload)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Orders) error); ok {
		r0 = rf(ctx, payload)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
