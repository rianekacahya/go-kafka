package rest

import (
	"context"
	"github.com/labstack/echo"
	"github.com/rianekacahya/go-kafka/domain/entity"
)

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Request().Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = c.Response().Header().Get(echo.HeaderXRequestID)
			}

			ctx := context.WithValue(c.Request().Context(), entity.ContextRequestID, id)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
