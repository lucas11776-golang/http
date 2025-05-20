package reader

import (
	"io/fs"

	"github.com/open2b/scriggo"
)

type TestingReader struct {
	files scriggo.Files
	cache scriggo.Files
}

// Comment
func NewTestingReader(files scriggo.Files) *TestingReader {
	return &TestingReader{
		files: files,
		cache: make(scriggo.Files),
	}
}

// Comment
func (ctx *TestingReader) Open(name string) (fs.File, error) {
	return ctx.files.Open(name)
}
