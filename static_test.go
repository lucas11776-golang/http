package http

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/lucas11776-golang/http/utils/reader"
	"github.com/open2b/scriggo"
)

func TestStatic(t *testing.T) {
	t.Run("TestGetStyles", func(t *testing.T) {
		// css := "assets/css/main.css"

		// static := InitStatic(&staticReaderTest{})

	})
}

var staticReaderTestFS = scriggo.Files{
	"assets/css/main.css": []byte(strings.Join([]string{
		"body { margin: 0px !important; padding: 0px !important; background-color: limegreen; }",
	}, "\r\n")),
}

type staticReaderTest struct {
	cache scriggo.Files
}

// Comment
func (ctx *staticReaderTest) Open(name string) (fs.File, error) {
	return responseReaderTestFS.Open(name)
}

// Comment
func (ctx *staticReaderTest) Cache(name string) (scriggo.Files, error) {
	return reader.ReadCache(ctx, ctx.cache, name)
}
