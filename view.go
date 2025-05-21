package http

import (
	"io/fs"
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

type View struct {
	fs        fs.FS
	extension string
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

type DefaultViewReader struct {
	fs.FS
}

// Comment
func NewDefaultViewReader(views string) fs.FS {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return &DefaultViewReader{os.DirFS(path.Path(wd, views))}
}

// Comment
func NewView(fs fs.FS, extension string) *View {
	return &View{
		fs:        fs,
		extension: extension,
	}
}

// Comment
func (ctx *View) buildTemplate(view string, options *scriggo.BuildOptions) (*scriggo.Template, error) {
	return scriggo.BuildTemplate(ctx.fs, path.FileRealPath(view, ctx.extension), options)
}

// Comment
func (ctx *View) Read(view string, data ViewData) ([]byte, error) {
	globals := native.Declarations{}

	for key, value := range data {
		globals[key] = value
	}

	template, err := ctx.buildTemplate(view, &scriggo.BuildOptions{Globals: globals})

	if err != nil {
		return nil, err
	}

	writer := &ViewWriter{}

	if err := template.Run(writer, nil, nil); err != nil {
		return nil, err
	}

	return []byte(strings.ReplaceAll(string(writer.parsed), "\r\n\r\n", "\r\n")), nil
}
