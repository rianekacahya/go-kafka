package repository

import (
	"context"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/stretchr/testify/mock"
)

type Orders struct {
	mock.Mock
}

func (o *Orders) SaveOrder(ctx context.Context, order *entity.Orders) error {
	ret := o.Called(ctx, order)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Orders) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
