package http

import (
	"math/rand"
	"testing"

	"github.com/lucas11776-golang/http/config"
)

func TestServer(t *testing.T) {
	serve := Server("127.0.0.1", 0)

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

	serve.Close()
}
