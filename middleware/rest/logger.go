package rest

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/rianekacahya/go-kafka/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			fields := []zapcore.Field{
				zap.String("remote_ip", c.RealIP()),
				zap.String("time", time.Since(start).String()),
				zap.String("host", req.Host),
				zap.String("request", fmt.Sprintf("%s %s", req.Method, req.RequestURI)),
				zap.Int("status", res.Status),
				zap.Int64("size", res.Size),
				zap.String("user_agent", req.UserAgent()),
			}

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
				fields = append(fields, zap.String("request_id", id))
			}

			n := res.Status
			switch {
			case n >= 500:
				logger.GetLogger().Error("Rest Logger", fields...)
			case n >= 400:
				logger.GetLogger().Warn("Rest Logger", fields...)
			case n >= 300:
				logger.GetLogger().Info("Rest Logger", fields...)
			default:
				logger.GetLogger().Info("Rest Logger", fields...)
			}

			return nil
		}
	}
}
