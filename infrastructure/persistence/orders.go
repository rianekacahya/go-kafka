package persistence

import (
	"context"
	"database/sql"
	newrelic "github.com/newrelic/go-agent"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/pkg/crashy"
	"github.com/rianekacahya/go-kafka/pkg/helper"
	"github.com/rianekacahya/go-kafka/pkg/mysql"
	"github.com/rianekacahya/go-kafka/pkg/telemetry"
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
		ntx   = telemetry.ContextTelemetry(ctx)
	)

	tx, err := o.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// saving batch order items
	query.WriteString(`INSERT INTO order_headers (customer_email, purchase_code, purchase_date) VALUES (?, ?, ?)`)
	bind = []interface{}{order.CustomerEmail, order.PurchaseCode, order.PurchaseDate}

	// newrelic datastore
	s := newrelic.DatastoreSegment{
		StartTime:          newrelic.StartSegmentNow(ntx),
		Product:            newrelic.DatastoreMySQL,
		Collection:         "order_headers",
		Operation:          "INSERT",
		ParameterizedQuery: mysql.TransformQuery(query.String(), bind),
	}

	// saving order headers
	_, err = tx.ExecContext(ctx, query.String(), bind...)

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return crashy.Wrap(err, crashy.ErrCodeDataWrite, "rollback failed")
		}

		return crashy.Wrap(err, crashy.ErrCodeDataWrite, "saving data order header failed")
	}

	if err := s.End(); err != nil {
		return crashy.Wrap(err, crashy.ErrCodeSend, "ending newrelic segment failed")
	}

	// reset strings builder query
	query.Reset()
	bind = nil

	// saving batch order items
	query.WriteString(`INSERT INTO order_items (purchase_code, sku, quantity) VALUES `)

	for _, v := range order.OrderItems {
		query.WriteString(`(?, ?, ?),`)
		bind = append(bind, order.PurchaseCode, v.SKU, v.Quantity)
	}

	s = newrelic.DatastoreSegment{
		StartTime:          newrelic.StartSegmentNow(ntx),
		Product:            newrelic.DatastoreMySQL,
		Collection:         "order_items",
		Operation:          "INSERT",
		ParameterizedQuery: mysql.TransformQuery(query.String(), bind),
	}

	_, err = tx.ExecContext(ctx, helper.TrimLastChar(query.String()), bind...)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return crashy.Wrap(err, crashy.ErrCodeDataWrite, "rollback failed")
		}

		return crashy.Wrap(err, crashy.ErrCodeDataWrite, "saving data order header failed")
	}

	// commit data
	if err := tx.Commit(); err != nil {
		return crashy.Wrap(err, crashy.ErrCodeDataWrite, "commit failed")
	}

	if err := s.End(); err != nil {
		return crashy.Wrap(err, crashy.ErrCodeSend, "ending newrelic segment failed")
	}

	return nil
}
