package strings

import "testing"

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
}
