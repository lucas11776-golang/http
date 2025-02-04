package server

import (
	"math/rand/v2"
	"testing"

	"github.com/lucas11776-golang/http/config"
)

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

	t.Run("TestDependency", func(t *testing.T) {
		dependency := rand.Float64() * 10000

		serve.Set("random", dependency)

		if serve.Get("random") != dependency {
			t.Fatalf("Expected dependency to be (%f) but got (%f)", dependency, serve.Get("random"))
		}
		if serve.Get("asset") != nil {
			t.Fatalf("Expected dependency to be nil but got (%s)", serve.Get("asset"))
		}
	})

	t.Run("TestConfig", func(t *testing.T) {
		jwt := "eye.hdlfashflahkh"

		config := serve.Get("config").(*config.Config).Set("JWT", jwt)

		jwtConfig := config.Get("JWT")

		if jwtConfig != jwt {
			t.Fatalf("Expected view path to be (%s) but go (%s)", jwt, jwtConfig)
		}
	})

	err = serve.Close()

	if err != nil {
		t.Fatalf("Failed to close server %s", err.Error())
	}
}
