package http

import (
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/lucas11776-golang/http/utils/extensions"
	"github.com/lucas11776-golang/http/utils/path"
	"github.com/lucas11776-golang/http/utils/reader"
	"github.com/open2b/scriggo"
)

type Static struct {
	fs reader.Cache
}

type defaultStaticReader struct {
	mutex sync.Mutex
	dir   string
	fs    scriggo.Files
}

// Comment
func (ctx *defaultStaticReader) Open(name string) (fs.File, error) {
	return reader.ReadCache(ctx, ctx.fs, fmt.Sprintf("%s/%s", strings.Trim("/", ctx.dir), strings.Trim("/", name)))
}

// Comment
func (ctx *defaultStaticReader) Write(name string, data []byte) error {
	ctx.mutex.Lock()
	ctx.fs[name] = data
	ctx.mutex.Unlock()

	return nil
}

// Comment
func DefaultStaticReader(dir string) *defaultStaticReader {
	wd, err := os.Getwd()

	if err != nil {
		log.Fatalf("Failed to get current working dir in view reader: %s", err.Error())
	}

	return &defaultStaticReader{
		dir: path.Path(wd, dir),
		fs:  make(scriggo.Files),
	}
}

// Comment
func InitStatic(fs reader.Cache) *Static {
	return &Static{fs: fs}
}

// Comment
func (ctx *Static) Read(name string) ([]byte, error) {

	file, err := ctx.fs.Open(name)

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
