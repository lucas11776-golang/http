package server

import (
	"testing"
)

func TestServer(t *testing.T) {
	serve := Serve("127.0.0.1", 0)

	t.Run("TestServe", func(t *testing.T) {
		if serve.Address() != "127.0.0.1" {
			t.Fatalf("Failed to start the server address %s", "127.0.0.1")
		}

		if serve.Port() == 0 {
			t.Fatalf("Server can not run in port %d", 0)
		}

		if serve.listener == nil {
			t.Fatalf("Server listener is not defined")
		}
	})

	serve.Close()
}
