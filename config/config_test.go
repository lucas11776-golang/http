package config

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("TestSetConfig", func(t *testing.T) {
		config := Init()

		random := strconv.Itoa(int(rand.Float64() * 10000))

		c := config.Set("SEED", random)

		if c.Get("SEED") != random {
			t.Fatalf("Expected config setting to be (%s) but got (%s)", random, c.Get("random"))
		}

		if c.Get("JWT_KEY") != "" {
			t.Fatalf("Expected config setting to be empty but got (%s)", c.Get("random"))
		}
	})
}
