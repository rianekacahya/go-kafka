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
	mockInvoker := new(invoker.Event)
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
		mockInvoker.On("OrderProducers", mock.Anything, mock.AnythingOfType("*entity.Orders")).Return(nil).Once()

		us := NewOrdersUsecase(mockOrdersRepository, mockInvoker)

		err := us.SubmitOrders(context.TODO(), &tempMockOrders)
		assert.NoError(t, err)
		mockInvoker.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		tempMockOrders := mockOrder

		tempMockOrders.PurchaseDate = entity.TimeNull(time.Now())
		mockInvoker.On("OrderProducers", mock.Anything, mock.AnythingOfType("*entity.Orders")).Return(errors.New("Unexpected Error")).Once()

		us := NewOrdersUsecase(mockOrdersRepository, mockInvoker)

		err := us.SubmitOrders(context.TODO(), &tempMockOrders)
		assert.Error(t, err)
		mockInvoker.AssertExpectations(t)
	})
}

func TestSaveOrders(t *testing.T) {
	mockOrdersRepository := new(repository.Orders)
	mockInvoker := new(invoker.Event)
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

		us := NewOrdersUsecase(mockOrdersRepository, mockInvoker)

		err := us.SaveOrders(context.TODO(), &tempMockOrders)
		assert.NoError(t, err)
		mockOrdersRepository.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		tempMockOrders := mockOrder
		tempMockOrders.PurchaseCode = entity.StrNull(helper.RandomString(20))
		mockOrdersRepository.On("SaveOrder", mock.Anything, mock.AnythingOfType("*entity.Orders")).Return(errors.New("Unexpected Error")).Once()

		us := NewOrdersUsecase(mockOrdersRepository, mockInvoker)

		err := us.SaveOrders(context.TODO(), &tempMockOrders)
		assert.Error(t, err)
		mockOrdersRepository.AssertExpectations(t)
	})
}