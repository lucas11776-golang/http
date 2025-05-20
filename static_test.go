package http

import (
	"io"
	"io/fs"
	"strings"
	"sync"
	"testing"

	"github.com/lucas11776-golang/http/types"
	"github.com/open2b/scriggo"
)

func TestStatic(t *testing.T) {
	var cssName = "assets/css/main.css"

	var cssContent = strings.Join([]string{
		"body { margin: 0px !important; padding: 0px !important; background-color: limegreen; }",
	}, "\r\n")

	static := InitStatic(&StaticReaderTest{
		Files: scriggo.Files{
			cssName: []byte(cssContent),
		},
	})

	t.Run("TestGetStyles", func(t *testing.T) {
		data, err := static.Read(cssName)

		if err != nil {
			t.Fatalf("Something went wrong when reading file: %s", err.Error())
		}

		if string(data) != cssContent {
			t.Fatalf("Expected styles to be but (%s) but go (%s)", cssContent, string(data))
		}
	})

	t.Run("TestHandleRequest", func(t *testing.T) {
		request, err := NewRequest("GET", cssName, "HTTP/1.1", types.Headers{"Accept": "text/css"}, strings.NewReader(""))

		if err != nil {
			t.Errorf("Something went wrong when trying to create request: %s", err.Error())
		}

		response, err := static.HandleRequest(request)

		if err != nil {
			t.Fatalf("Something went wrong when getting static: %s", err.Error())
		}

		if response.GetHeader("content-type") != "text/css" {
			t.Fatalf("Expected content type to be but (%s) but go (%s)", "text/css", response.GetHeader("content-type"))
		}

		tBody, err := io.ReadAll(response.Body)

		if err != nil {
			t.Fatalf("Something went wrong went trying to read body: %s", tBody)
		}

		if cssContent != string(tBody) {
			t.Fatalf("Expected content body to be but (%s) but go (%s)", cssContent, string(tBody))
		}
	})
}

type StaticReaderTest struct {
	mutex sync.Mutex
	Files scriggo.Files
}

// Comment
func (ctx *StaticReaderTest) Open(name string) (fs.File, error) {
	return ctx.Files.Open(name)
}

// Comment
func (ctx *StaticReaderTest) Write(name string, data []byte) error {
	ctx.mutex.Lock()

	ctx.Files[name] = data

	ctx.mutex.Unlock()

	return nil
}
