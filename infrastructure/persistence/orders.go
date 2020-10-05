package persistence

import (
	"context"
	"database/sql"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/helper"
	"github.com/rianekacahya/go-kafka/pkg/crashy"
	"strings"
)

type orders struct {
	database *sql.DB
}

func NewOrdersRepository(database *sql.DB) *orders {
	return &orders{database}
}

func (o *orders) SaveOrder(ctx context.Context, order *entity.Orders) error {
	var (
		query strings.Builder
		bind  []interface{}
	)

	tx, err := o.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// saving order headers
	_, err = tx.ExecContext(ctx, `INSERT INTO order_headers (customer_email, purchase_code, purchase_date) VALUES (?, ?, ?)`,
		order.CustomerEmail, order.PurchaseCode, order.PurchaseDate,
	)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return crashy.Wrap(err, crashy.ErrCodeDataWrite, "rollback failed")
		}

		return crashy.Wrap(err, crashy.ErrCodeDataWrite, "saving data order header failed")
	}

	// saving batch order items
	query.WriteString(`INSERT INTO order_items (purchase_code, sku, quantity) VALUES `)

	for _, v := range order.OrderItems {
		query.WriteString(`(?, ?, ?),`)
		bind = append(bind, order.PurchaseCode, v.SKU, v.Quantity)
	}

	_, err = tx.ExecContext(ctx, helper.TrimLastChar(query.String()), bind...)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return crashy.Wrap(err, crashy.ErrCodeDataWrite, "rollback failed")
		}

		return crashy.Wrap(err, crashy.ErrCodeDataWrite, "saving data order header failed")
	}

	// commit data
	if  err := tx.Commit();  err != nil {
		return crashy.Wrap(err, crashy.ErrCodeDataWrite, "commit failed")
	}

	return nil
}
