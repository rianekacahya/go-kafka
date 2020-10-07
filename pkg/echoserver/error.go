package echoserver

import (
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func Handler(err error, c echo.Context) {
	type errors struct {
		Code    string      `json:"code"`
		Message interface{} `json:"message"`
	}

	type response struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data,omitempty"`
		Error   *errors     `json:"error,omitempty"`
	}

	var (
		res  = new(response)
		code = http.StatusInternalServerError
		msg  interface{}
	)

	switch e := err.(type) {
	case *echo.HTTPError:
		code = e.Code
		msg = e.Message
		if e.Internal != nil {
			msg = fmt.Sprintf("%v, %v", err, e.Internal)
		}
	default:
		msg = http.StatusText(code)
	}

	if _, ok := msg.(string); ok {
		res.Error = &errors{
			Code:    fmt.Sprint(code),
			Message: msg,
		}
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == "HEAD" {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, res)
		}
		if err != nil {
			log.Printf("echo error wraper failed responding, error: %s", err)
		}
	}
}
