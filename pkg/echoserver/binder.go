package echoserver

import "github.com/labstack/echo"

// CustomBinder struct
type CustomBinder struct {
	bind echo.Binder
}

func NewBinder() *CustomBinder {
	return &CustomBinder{bind: &echo.DefaultBinder{}}
}

// Bind tries to bind request into interface, and if it does then validate it
func (cb *CustomBinder) Bind(i interface{}, c echo.Context) error {
	if err := cb.bind.Bind(i, c); err != nil && err != echo.ErrUnsupportedMediaType {
		return err
	}

	return nil
}
