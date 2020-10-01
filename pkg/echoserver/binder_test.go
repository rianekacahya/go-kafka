package echoserver

import (
	"bytes"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Req struct {
	Name string `json:"name" validate:"required"`
}

func TestBind(t *testing.T) {
	cases := []struct {
		name     string
		req      string
		wantErr  bool
		wantData *Req
	}{
		{
			name:     "Fail on binding",
			wantErr:  true,
			req:      `"failed string"`,
			wantData: &Req{Name: ""},
		},
		{
			name:     "Success",
			req:      `{"name":"Hello"}`,
			wantData: &Req{Name: "Hello"},
		},
	}
	b := NewBinder()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "", bytes.NewBufferString(tt.req))
			req.Header.Set("Content-Type", "application/json")
			e := echo.New()
			e.Binder = NewBinder()
			c := e.NewContext(req, w)
			r := new(Req)
			err := b.Bind(r, c)
			assert.Equal(t, tt.wantData, r)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}