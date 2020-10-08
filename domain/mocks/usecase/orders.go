package usecase

import (
	"context"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/stretchr/testify/mock"
)

type Orders struct {
	mock.Mock
}

func (o *Orders) SubmitOrders(ctx context.Context, req *entity.Orders) error {
	ret := o.Called(ctx, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Orders) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (o *Orders) SaveOrders(ctx context.Context, req *entity.Orders) error {
	ret := o.Called(ctx, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Orders) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}