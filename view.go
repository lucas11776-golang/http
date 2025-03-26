package http

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/lucas11776-golang/http/utils/path"
	"github.com/lucas11776-golang/http/utils/reader"
	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/native"
)

type ViewData map[string]interface{}

type viewWriter struct {
	parsed []byte
}

type defaultViewReader struct {
	dir   string
	cache scriggo.Files
}

type View struct {
	fs        reader.CacheReader
	extension string
}

type ViewInterface interface {
	Read(view string, data ViewData) ([]byte, error)
}

// Comment
func (ctx *viewWriter) Write(p []byte) (n int, err error) {
	ctx.parsed = append(ctx.parsed, p...)

	return len(ctx.parsed), nil
}

// Comment
func (ctx *viewWriter) Parsed() []byte {
	return ctx.parsed
}

// Comment
func (ctx *defaultViewReader) Open(name string) (fs.File, error) {
	return os.Open(path.Path(ctx.dir, name))
}

// Comment
func (ctx *defaultViewReader) Cache(name string) (scriggo.Files, error) {
	return reader.ReadCache(ctx, ctx.cache, name)
}

// Comment
func (ctx *defaultViewReader) Write(name string, data []byte) error {
	return nil
}

// Comment
func DefaultViewReader(views string) *defaultViewReader {
	wd, err := os.Getwd()

	if err != nil {
		log.Fatalf("Failed to get current working dir in view reader: %s", err.Error())
	}

	return &defaultViewReader{
		dir:   path.Path(wd, views),
		cache: make(scriggo.Files),
	}
}

// Comment
func InitView(fs reader.CacheReader, extension string) *View {
	return &View{
		fs:        fs,
		extension: extension,
	}
}

// Comment
func (ctx *View) Read(view string, data ViewData) ([]byte, error) {
	globals := native.Declarations{}

	for key, value := range data {
		globals[key] = value
	}

	vw := strings.Join([]string{strings.ReplaceAll(view, ".", "/"), ctx.extension}, ".")

	views, err := ctx.fs.Cache(vw)

	if err != nil {
		return nil, err
	}

	template, err := scriggo.BuildTemplate(views, vw, &scriggo.BuildOptions{
		Globals: globals,
	})

	if err != nil {
		return nil, err
	}

	writer := viewWriter{}

	if err := template.Run(&writer, nil, nil); err != nil {
		return nil, err
	}

	return []byte(strings.ReplaceAll(string(writer.parsed), "\r\n\r\n", "\r\n")), nil
}

type ViewReaderTest struct {
	mutex sync.Mutex
	Files scriggo.Files
	cache scriggo.Files
}

// Comment
func (ctx *ViewReaderTest) Open(name string) (fs.File, error) {
	return ctx.Files.Open(name)
}

// Comment
func (ctx *ViewReaderTest) Cache(name string) (scriggo.Files, error) {

	fmt.Println("CACHE NAME", name)

	return reader.ReadCache(ctx, ctx.cache, name)
}

// Comment
func (ctx *ViewReaderTest) Write(name string, data []byte) error {
	ctx.mutex.Lock()
	ctx.cache[name] = data
	ctx.mutex.Unlock()
	return nil
}
