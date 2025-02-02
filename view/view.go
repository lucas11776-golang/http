package view

import (
	"io/fs"
	"os"
	"strings"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/native"
)

type Data map[string]interface{}

type viewFs struct {
	dir string
}

type Reader interface {
	Open(name string) (fs.File, error)
}

type viewWriter struct {
	parsed []byte
}

type View struct {
	fs        fs.FS
	extension string
}

// Comment
func (ctx *viewFs) cleanPath(path string) string {
	return strings.Trim(strings.Trim(path, "\\"), "/")
}

// Comment
func (ctx *viewFs) viewPath(view string) string {
	return strings.ReplaceAll(strings.Join([]string{ctx.cleanPath(ctx.dir), ctx.cleanPath(view)}, "\\"), "/", "\\")
}

// Comment
func (ctx *viewFs) Open(view string) (fs.File, error) {
	return os.OpenFile(ctx.viewPath(view), os.O_RDONLY, fs.ModeExclusive)
}

// Comment
func FileSystem(views string) fs.FS {
	return &viewFs{
		dir: views,
	}
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
func Init(fsys Reader, extension string) *View {
	return &View{
		fs:        fsys,
		extension: extension,
	}
}

// Comment
func (ctx *View) Read(view string, data Data) ([]byte, error) {
	globals := native.Declarations{}

	if data != nil {
		for key, value := range data {
			globals[key] = value
		}
	}

	template, err := scriggo.BuildTemplate(ctx.fs, strings.Join([]string{view, ctx.extension}, "."), &scriggo.BuildOptions{
		Globals: globals,
	})

	if err != nil {
		return nil, err
	}

	writer := viewWriter{}

	template.Run(&writer, nil, nil)

	return []byte(strings.ReplaceAll(string(writer.parsed), "\r\n\r\n", "\r\n")), nil
}
