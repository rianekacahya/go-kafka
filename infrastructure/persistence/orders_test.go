package persistence

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/pkg/helper"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
	"time"
)

func TestSaveOrderCommit(t *testing.T) {
	order := &entity.Orders{
		OrderHeaders: &entity.OrderHeaders{
			CustomerEmail: entity.StrNull("rian.eka.cahya@gmail.com"),
			PurchaseCode:  entity.StrNull(helper.RandomString(20)),
			PurchaseDate:  entity.TimeNull(time.Now()),
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

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectBegin()
	bind := []driver.Value{order.CustomerEmail, order.PurchaseCode, order.PurchaseDate}
	mock.ExpectExec("INSERT INTO order_headers").
		WithArgs(order.CustomerEmail, order.PurchaseCode, order.PurchaseDate).
		WillReturnResult(sqlmock.NewResult(1, 1))

	bind = []driver.Value{}
	for _, v := range order.OrderItems {
		bind = append(bind, order.PurchaseCode, v.SKU, v.Quantity)
	}
	mock.ExpectExec("INSERT INTO order_items").
		WithArgs(bind...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	or := NewOrdersRepository(db)

	if err := or.SaveOrder(context.TODO(), order); err != nil {
		t.Errorf("error was not expected while saving order: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveOrderRollback(t *testing.T) {
	order := &entity.Orders{
		OrderHeaders: &entity.OrderHeaders{
			CustomerEmail: entity.StrNull("rian.eka.cahya@gmail.com"),
			PurchaseCode:  entity.StrNull(helper.RandomString(20)),
			PurchaseDate:  entity.TimeNull(time.Now()),
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

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectBegin()
	bind := []driver.Value{order.CustomerEmail, order.PurchaseCode, order.PurchaseDate}
	mock.ExpectExec("INSERT INTO order_headers").
		WithArgs(order.CustomerEmail, order.PurchaseCode, order.PurchaseDate).
		WillReturnResult(sqlmock.NewResult(1, 1))

	bind = []driver.Value{}
	for _, v := range order.OrderItems {
		bind = append(bind, order.PurchaseCode, v.SKU, v.Quantity)
	}
	mock.ExpectExec("INSERT INTO order_items").
		WithArgs(bind...).
		WillReturnError(fmt.Errorf("saving data order header failed"))

	mock.ExpectRollback()

	or := NewOrdersRepository(db)

	if err := or.SaveOrder(context.TODO(), order); err == nil {
		t.Errorf("was expecting an error while saving order: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
