package echoserver

import "testing"

func TestServer(t *testing.T) {
	e := NewServer(false)
	if e == nil {
		t.Errorf("Server should not be nil")
	}
}
