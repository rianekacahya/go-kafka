package response

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/rianekacahya/go-kafka/pkg/crashy"
	"net/http"
)

type Errors struct {
	Code    string      `json:"code"`
	Message interface{} `json:"message"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Errors     `json:"error,omitempty"`
}

func Render(c echo.Context, data interface{}) error {
	var response = new(Response)

	response.Success = true
	response.Data = data

	return c.JSON(http.StatusOK, response)
}

func Error(c echo.Context, err error) error {
	var (
		httpcode int
		code     string
		message  string
		response = new(Response)
	)

	switch err.(type) {
	case *crashy.Error:
		switch err.(*crashy.Error).Code() {
		default:
			httpcode = http.StatusInternalServerError
		case crashy.ErrCodeNetBuild, crashy.ErrCodeNetConnect:
			httpcode = http.StatusBadGateway
		case crashy.ErrCodeValidation, crashy.ErrCodeFormatting:
			httpcode = http.StatusBadRequest
		case crashy.ErrCodeDataRead, crashy.ErrCodeDataWrite:
			httpcode = http.StatusServiceUnavailable
		case crashy.ErrCodeNoResult:
			httpcode = http.StatusNotFound
		case crashy.ErrCodeUnauthorized, crashy.ErrCodeExpired:
			httpcode = http.StatusUnauthorized
		case crashy.ErrCodeForbidden:
			httpcode = http.StatusForbidden
		case crashy.ErrCodeTooManyRequest:
			httpcode = http.StatusTooManyRequests
		case crashy.ErrCodeDataIncomplete:
			httpcode = http.StatusPartialContent
		case crashy.ErrCodeSend:
			httpcode = http.StatusRequestTimeout
		}
		code = err.(*crashy.Error).Code()
		message = err.(*crashy.Error).Message()
		if c.Echo().Debug {
			message = err.(*crashy.Error).Unwrap().Error()
		}
	default:
		httpcode = http.StatusInternalServerError
		code = fmt.Sprint(http.StatusInternalServerError)
		message = http.StatusText(httpcode)
		if c.Echo().Debug {
			message = err.Error()
		}
	}

	response.Error = &Errors{
		Code:    code,
		Message: message,
	}

	return c.JSON(httpcode, response)
}
