package server

import "testing"

func TestServer(t *testing.T) {
	serve, err := Serve("127.0.0.1", 0)

	if err != nil {
		t.Fatalf("Failed to start the server: %s", err.Error())
	}

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

	t.Run("TestConfig", func(t *testing.T) {
		serve.SetConfig("VIEW_PATH", "views")

		if serve.GetConfig("VIEW_PATH") != "views" {
			t.Fatalf("Expected view path to be (%s) but go (%s)", "views", serve.GetConfig("VIEW_PATH"))
		}
	})

	err = serve.Close()

	if err != nil {
		t.Fatalf("Failed to close server %s", err.Error())
	}
}
