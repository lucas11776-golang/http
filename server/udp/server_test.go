package udp

import (
	"testing"
)

func TestServer(t *testing.T) {
	serve := Serve("127.0.0.1", 0)

	t.Run("TestServe", func(t *testing.T) {
		if serve.Address() != "127.0.0.1" {
			t.Fatalf("Expected address to be (%s) but got (%s)", "127.0.0.1", serve.Address())
		}
	})

	serve.Close()
}
