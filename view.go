package http

import (
	"io"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/lucas11776-golang/http/utils/path"
	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/native"
)

type ViewData map[string]interface{}

type ViewWriter struct {
	parsed []byte
}

type DefaultViewReader struct {
	dir       string
	extension string
	files     scriggo.Files
	fs        fs.FS
}

type View struct {
	fs fs.FS
}

type ViewInterface interface {
	Read(view string, data ViewData) ([]byte, error)
}

// Comment
func (ctx *ViewWriter) Write(p []byte) (n int, err error) {
	ctx.parsed = append(ctx.parsed, p...)

	return len(ctx.parsed), nil
}

// Comment
func (ctx *ViewWriter) Parsed() []byte {
	return ctx.parsed
}

// Comment
func (ctx *DefaultViewReader) Open(name string) (fs.File, error) {
	if file, err := ctx.files.Open(path.FileRealPath(name, ctx.extension)); err == nil {
		return file, nil
	}

	file, err := ctx.fs.Open(path.FileRealPath(name, ctx.extension))

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	ctx.files[path.FileRealPath(name, ctx.extension)] = data

	return ctx.files.Open(path.FileRealPath(name, ctx.extension))
}

// Comment
func NewDefaultViewReader(views string, extension string) *DefaultViewReader {
	wd, err := os.Getwd()

	if err != nil {
		log.Fatalf("Failed to get current working dir in view reader: %s", err.Error())
	}

	return &DefaultViewReader{
		dir:       path.Path(wd, views),
		extension: extension,
		files:     make(scriggo.Files),
		fs:        os.DirFS(views),
	}
}

// Comment
func NewView(fs fs.FS) *View {
	return &View{
		fs: fs,
	}
}

// Comment
func (ctx *View) Read(view string, data ViewData) ([]byte, error) {
	globals := native.Declarations{}

	for key, value := range data {
		globals[key] = value
	}

	template, err := scriggo.BuildTemplate(ctx.fs, view, &scriggo.BuildOptions{
		Globals: globals,
	})

	if err != nil {
		return nil, err
	}

	writer := ViewWriter{}

	if err := template.Run(&writer, nil, nil); err != nil {
		return nil, err
	}

	return []byte(strings.ReplaceAll(string(writer.parsed), "\r\n\r\n", "\r\n")), nil
}
