package strings

import (
	"math/rand/v2"
	"testing"
)

func TestStrings(t *testing.T) {

	t.Run("TestJoinPath", func(t *testing.T) {
		path := JoinPath("/api", "users", "/1")

		if path != "api/users/1" {
			t.Fatalf("Expected joined string path to be (%s) but go (%s)", "api/users/1", path)
		}

		path = JoinPath("/api", "/products", "/", "")

		if path != "api/products" {
			t.Fatalf("Expected joined string path to be (%s) but go (%s)", "api/products", path)
		}

		path = JoinPath("/")

		if path != "" {
			t.Fatalf("Expected joined string path to be (%s) but go (%s)", "", path)
		}
	})

	t.Run("TestRandom", func(t *testing.T) {
		size := int(rand.Float64() * 1000)
		str := Random(size)

		if len(str) != size {
			t.Fatalf("Expected string size to be (%d) but go (%d)", size, len(str))
		}
	})
}
