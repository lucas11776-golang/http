package http

import (
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/lucas11776-golang/http/utils/extensions"
	"github.com/lucas11776-golang/http/utils/path"
	"github.com/lucas11776-golang/http/utils/reader"
	"github.com/open2b/scriggo"
)

type Static struct {
	fs reader.CacheReader
}

type StaticReader struct {
	dir string
}

type defaultStaticReader struct {
	dir   string
	cache scriggo.Files
}

// Comment
func (ctx *defaultStaticReader) Open(name string) (fs.File, error) {
	return os.Open(path.Path(ctx.dir, name))
}

// Comment
func (ctx *defaultStaticReader) Cache(name string) (scriggo.Files, error) {
	// TODO Should I set max cache size.
	return reader.ReadCache(ctx, ctx.cache, name)
}

// Comment
func DefaultStaticReader(dir string) *defaultStaticReader {
	wd, err := os.Getwd()

	if err != nil {
		log.Fatalf("Failed to get current working dir in view reader: %s", err.Error())
	}

	return &defaultStaticReader{
		dir:   path.Path(wd, dir),
		cache: make(scriggo.Files),
	}
}

// Comment
func InitStatic(fs reader.CacheReader) *Static {
	return &Static{fs: fs}
}

// Comment
func (ctx *Static) Read(name string) ([]byte, error) {
	files, err := ctx.fs.Cache(name)

	if err != nil {
		return nil, err
	}

	file, err := files.Open(name)

	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()

	if err != nil {
		return nil, err
	}

	data := make([]byte, stat.Size())

	_, err = file.Read(data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// Comment
func (ctx *Static) HandleRequest(req *Request) (*Response, error) {
	u, err := url.Parse(req.Path())

	if err != nil {
		return nil, fmt.Errorf("File does not exists in statics (%s)", req.Path())
	}

	req.Response._Body, err = ctx.Read(u.Path)

	if err != nil {
		return nil, err
	}

	p := strings.Split(u.Path, ".")

	req.Response.SetHeader("content-type", extensions.ContentType(p[len(p)-1]))

	return req.Response, nil
}
