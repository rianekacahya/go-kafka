package rest

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	newrelic "github.com/newrelic/go-agent"
	"github.com/rianekacahya/go-kafka/pkg/telemetry"
	"strings"
)

func Telemetry(app newrelic.Application) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			txnName := fmt.Sprintf("%s%s", strings.ToLower(c.Request().Method), c.Path())
			txn := app.StartTransaction(txnName, c.Response().Writer, c.Request())
			defer txn.End()

			ctx := context.WithValue(c.Request().Context(), telemetry.ContextKeyTelemetry, txn)
			c.SetRequest(c.Request().WithContext(ctx))
			if err := next(c); err != nil {
				return err
			}

			return nil
		}
	}
}
