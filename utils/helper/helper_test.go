package helper

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHelper(t *testing.T) {
	t.Run("TestUrl", func(t *testing.T) {
		host := "http://localhost:8080/"

		if err := os.Setenv("APP_URL", host); err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("%s/products/1", strings.TrimRight(host, "/"))
		actual := Url("products", 1)

		if expected != actual {
			t.Fatalf("expected url to be (%s) but got (%s)", expected, actual)
		}
	})

	t.Run("TestSubdomain", func(t *testing.T) {
		host := "http://example.com/"

		if err := os.Setenv("APP_URL", host); err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprint("http://api.example.com/products/1")
		actual := Subdomain("api", "products", 1)

		if expected != actual {
			t.Fatalf("expected subdomain to be (%s) but got (%s)", expected, actual)
		}
	})

	t.Run("TestSubdomain", func(t *testing.T) {
		tm := time.Now()

		expected := tm.Format(time.DateTime)
		actual := Format(tm, time.DateTime)

		if expected != actual {
			t.Fatalf("expected format to be (%s) but got (%s)", expected, actual)
		}
	})
}
