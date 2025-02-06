package http

import (
	"io/fs"

	"github.com/open2b/scriggo"
)

type StaticReader interface {
	Open(name string) (fs.File, error)
	Statics(name string) (scriggo.Files, error)
}

type Static struct {
	fs StaticReader
}

type DefaultStaticReader struct {
}

// Comment
func (ctx *DefaultStaticReader) Open(name string) (fs.File, error) {
	return nil, nil
}

// Comment
func (ctx *DefaultStaticReader) Statics(name string) (scriggo.Files, error) {
	return nil, nil
}

// Comment
func InitStatic(fs StaticReader) *Static {
	return &Static{fs: fs}
}

// Comment
func (ctx *Static) Read(file string) ([]byte, error) {

	return nil, nil
}
