package echoserver

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func NewServer(debug bool) *echo.Echo {
	// init echo
	e := echo.New()

	// Hide banner
	e.HideBanner = true

	// Set debug status parameter
	e.Debug = debug

	// init default adapter
	e.Use(
		middleware.RequestID(),
		middleware.RecoverWithConfig(middleware.RecoverConfig{
			DisableStackAll:   true,
			DisablePrintStack: false,
		}),
	)

	// custom error handler
	e.HTTPErrorHandler = Handler

	// Custom binder
	e.Binder = &CustomBinder{bind: &echo.DefaultBinder{}}

	return e
}
