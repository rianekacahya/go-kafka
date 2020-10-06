package orders

import (
	"github.com/labstack/echo"
	newrelic "github.com/newrelic/go-agent"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/domain/usecase"
	"github.com/rianekacahya/go-kafka/pkg/crashy"
	"github.com/rianekacahya/go-kafka/pkg/echoserver/response"
	"github.com/rianekacahya/go-kafka/pkg/telemetry"
)

type rest struct {
	ordersUsecase usecase.OrdersUsecase
}

func NewHandler(echo *echo.Group, ordersUsecase usecase.OrdersUsecase) {
	transport := rest{ordersUsecase}

	routes := echo.Group("/orders")
	routes.POST("", transport.submit)
}

func (r *rest) submit(c echo.Context) error {
	var (
		err error
		req = new(entity.Orders)
		ntx = telemetry.ContextTelemetry(c.Request().Context())
	)

	// bind request body
	if err = c.Bind(req); err != nil {
		return response.Error(c, crashy.Wrap(err, crashy.ErrCodeFormatting, "request body not valid"))
	}

	// newrelic submit order segment
	defer newrelic.StartSegment(ntx, "Hanlder.HTTP.Submit").End()

	// call usecase
	if err = r.ordersUsecase.SubmitOrders(c.Request().Context(), req); err != nil {
		return response.Error(c, err)
	}

	return response.Render(c, "success")
}
