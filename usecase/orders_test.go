package usecase

import (
	"context"
	"errors"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/domain/mocks/invoker"
	"github.com/rianekacahya/go-kafka/domain/mocks/repository"
	"github.com/rianekacahya/go-kafka/pkg/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestSubmitOrders(t *testing.T) {
	mockOrdersRepository := new(repository.Orders)
	mockInvokerEvent := new(invoker.Event)
	mockInvokerSync := new(invoker.Sync)
	mockOrder := entity.Orders{
		OrderHeaders: &entity.OrderHeaders{
			CustomerEmail: entity.StrNull("rian.eka.cahya@gmail.com"),
		},
		OrderItems: []entity.OrderItems{
			{
				SKU: entity.StrNull("HYUGUYGUGB-KK"),
				Quantity: entity.Int32Null(20),
			},
			{
				SKU: entity.StrNull("NJNKGTRETR-RT"),
				Quantity: entity.Int32Null(5),
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		tempMockOrders := mockOrder

		tempMockOrders.PurchaseDate = entity.TimeNull(time.Now())
		mockInvokerEvent.On("OrderProducers", mock.Anything, mock.AnythingOfType("*entity.Orders")).Return(nil).Once()

		us := NewOrdersUsecase(mockOrdersRepository, mockInvokerEvent, mockInvokerSync)

		err := us.SubmitOrders(context.TODO(), &tempMockOrders)
		assert.NoError(t, err)
		mockInvokerEvent.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		tempMockOrders := mockOrder

		tempMockOrders.PurchaseDate = entity.TimeNull(time.Now())
		mockInvokerEvent.On("OrderProducers", mock.Anything, mock.AnythingOfType("*entity.Orders")).Return(errors.New("Unexpected Error")).Once()

		us := NewOrdersUsecase(mockOrdersRepository, mockInvokerEvent, mockInvokerSync)

		err := us.SubmitOrders(context.TODO(), &tempMockOrders)
		assert.Error(t, err)
		mockInvokerEvent.AssertExpectations(t)
	})
}

func TestSaveOrders(t *testing.T) {
	mockOrdersRepository := new(repository.Orders)
	mockInvokerEvent := new(invoker.Event)
	mockInvokerSync := new(invoker.Sync)
	mockOrder := entity.Orders{
		OrderHeaders: &entity.OrderHeaders{
			CustomerEmail: entity.StrNull("rian.eka.cahya@gmail.com"),
		},
		OrderItems: []entity.OrderItems{
			{
				SKU: entity.StrNull("HYUGUYGUGB-KK"),
				Quantity: entity.Int32Null(20),
			},
			{
				SKU: entity.StrNull("NJNKGTRETR-RT"),
				Quantity: entity.Int32Null(5),
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		tempMockOrders := mockOrder
		tempMockOrders.PurchaseCode = entity.StrNull(helper.RandomString(20))
		mockOrdersRepository.On("SaveOrder", mock.Anything, mock.AnythingOfType("*entity.Orders")).Return(nil).Once()
		mockInvokerSync.On("EventWhisper", mock.Anything, mock.AnythingOfType("string")).Return(nil)

		us := NewOrdersUsecase(mockOrdersRepository, mockInvokerEvent, mockInvokerSync)

		err := us.SaveOrders(context.TODO(), &tempMockOrders)
		assert.NoError(t, err)
		mockOrdersRepository.AssertExpectations(t)
		mockInvokerSync.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		tempMockOrders := mockOrder
		tempMockOrders.PurchaseCode = entity.StrNull(helper.RandomString(20))
		mockOrdersRepository.On("SaveOrder", mock.Anything, mock.AnythingOfType("*entity.Orders")).Return(errors.New("Unexpected Error")).Once()
		mockInvokerSync.On("EventWhisper", mock.Anything, mock.AnythingOfType("string")).Return(errors.New("Unexpected Error"))

		us := NewOrdersUsecase(mockOrdersRepository, mockInvokerEvent, mockInvokerSync)

		err := us.SaveOrders(context.TODO(), &tempMockOrders)
		assert.Error(t, err)
		mockOrdersRepository.AssertExpectations(t)
		mockInvokerSync.AssertExpectations(t)
	})
}
