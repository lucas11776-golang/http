package http

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/lucas11776-golang/http/utils/extensions"
	"github.com/lucas11776-golang/http/utils/path"
)

type Static struct {
	fs fs.FS
}

type DefaultStaticReader struct {
	fs fs.FS
}

// Comment
func (ctx *DefaultStaticReader) Open(name string) (fs.File, error) {
	return ctx.fs.Open(strings.Split(name, "?")[0])
}

// Comment
func NewDefaultStaticReader(dir string) *DefaultStaticReader {
	wd, err := os.Getwd()

	if err != nil {
		log.Fatalf("Failed to get current working dir in view reader: %s", err.Error())
	}

	return &DefaultStaticReader{
		fs: os.DirFS(path.Path(wd, dir)),
	}
}

// Comment
func InitStatic(fs fs.FS) *Static {
	return &Static{fs: fs}
}

// Comment
func (ctx *Static) Read(name string) ([]byte, error) {
	file, err := ctx.fs.Open(name)

	if err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}

// Comment
func (ctx *Static) HandleRequest(req *Request) (*Response, error) {
	u, err := url.Parse(req.Path())

	if err != nil {
		return nil, fmt.Errorf("file does not exists in statics (%s)", req.Path())
	}

	body, err := ctx.Read(u.Path)

	if err != nil {
		return nil, err
	}

	req.Response.SetBody(body)

	p := strings.Split(u.Path, ".")

	req.Response.SetHeader("content-type", extensions.ContentType(p[len(p)-1]))

	return req.Response, nil
}
