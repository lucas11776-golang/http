package headers

import (
	"strings"
	"testing"

	"github.com/lucas11776-golang/http/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestHeader(t *testing.T) {
	t.Run("TestToHeader", func(t *testing.T) {
		name := "content-type"
		value := "application/json"

		headers := types.Headers{name: value}

		header := ToHeader(headers)

		h, ok := header[cases.Title(language.English).String(name)]

		if !ok {
			t.Fatalf("Failed to get header (%s)", value)
		}

		if strings.Join(h, ",") != value {
			t.Fatalf("Expected header to be (%s) but got (%s)", value, h)
		}
	})
}
