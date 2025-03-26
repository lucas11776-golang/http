package tcp

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	serve := Serve("127.0.0.1", 0)

	t.Run("TestServe", func(t *testing.T) {
		if serve.Address() != "127.0.0.1" {
			t.Fatalf("Expected address to be (%s) but got (%s)", "127.0.0.1", serve.Address())
		}

		if serve.Port() == 0 {
			t.Fatalf("Expected port not to be (0) but got (%d)", 0)
		}

		if serve.listener.Addr().String() != fmt.Sprintf("%s:%d", serve.Address(), serve.Port()) {
			t.Fatalf(
				"Expected listener address to be (%s) but got (%s)",
				fmt.Sprintf("%s:%d", serve.Address(), serve.Port()),
				serve.listener.Addr().String(),
			)
		}
	})

	serve.Close()
}
