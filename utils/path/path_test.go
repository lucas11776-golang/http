package path

import "testing"

func TestPath(t *testing.T) {
	t.Run("TestFileRealPath", func(t *testing.T) {
		if result := FileRealPath("index", "html"); result != "index.html" {
			t.Fatalf("Expected result to be (%s) but got (%s)", "index.html", result)
		}

		if result := FileRealPath("views/index", "html"); result != "views/index.html" {
			t.Fatalf("Expected result to be (%s) but got (%s)", "views/index.html", result)
		}
	})
}
