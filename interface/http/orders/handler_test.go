package orders

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/rianekacahya/go-kafka/domain/entity"
	"github.com/rianekacahya/go-kafka/domain/mocks/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSubmit(t *testing.T) {
	mockUsecase := new(usecase.Orders)
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

	j, err := json.Marshal(mockOrder)
	assert.NoError(t, err)

	mockUsecase.On("SubmitOrders", mock.Anything, mock.AnythingOfType("*entity.Orders")).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/v1/orders", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v1/orders")

	handler := rest {
		ordersUsecase: mockUsecase,
	}

	err = handler.submit(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUsecase.AssertExpectations(t)
}
