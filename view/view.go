package view

import (
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/native"
)

type Data map[string]interface{}

type viewWriter struct {
	parsed []byte
}

type viewReader struct {
	dir   string
	cache scriggo.Files
}

type Reader interface {
	Open(name string) (fs.File, error)
	Views(name string) (scriggo.Files, error)
}

type View struct {
	fs        Reader
	extension string
}

// Comment
func Path(path ...string) string {
	pth := []string{}

	for _, p := range path {
		pth = append(pth, strings.Trim(strings.ReplaceAll(p, "/", "\\"), "\\"))
	}

	return strings.Join(pth, "\\")
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
func ReadViewCache(reader Reader, cache scriggo.Files, view string) (scriggo.Files, error) {
	_, ok := cache[view]

	if ok {
		return cache, nil
	}

	file, err := reader.Open(view)

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

	cache[view] = data

	return cache, nil
}

// Comment
func (ctx *viewReader) Open(name string) (fs.File, error) {
	return os.Open(Path(ctx.dir, name))
}

// Comment
func (ctx *viewReader) Views(name string) (scriggo.Files, error) {
	return ReadViewCache(ctx, ctx.cache, name)
}

// Comment
func ViewReader(views string) *viewReader {
	wd, err := os.Getwd()

	if err != nil {
		log.Fatalf("Failed to get current working dir in view reader: %s", err.Error())
	}

	return &viewReader{
		dir:   Path(wd, views),
		cache: make(scriggo.Files),
	}
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

	vw := strings.Join([]string{view, ctx.extension}, ".")

	views, err := ctx.fs.Views(vw)

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

	template.Run(&writer, nil, nil)

	return []byte(strings.ReplaceAll(string(writer.parsed), "\r\n\r\n", "\r\n")), nil
}
