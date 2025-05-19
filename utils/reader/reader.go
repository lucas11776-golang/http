package reader

import (
	"io"
	"io/fs"
	"os"
	"strings"
	"sync"

	"github.com/open2b/scriggo"
)

type TestingReader struct {
	files scriggo.Files
	cache scriggo.Files
	mutex sync.Mutex
}

// TODO: Something feel wrong here...
type Cache interface {
	Open(name string) (fs.File, error)
	Write(name string, data []byte) error
}

// Comment
func ReadCache(reader Cache, cache scriggo.Files, name string) (fs.File, error) {
	if _, ok := cache[name]; ok {
		return cache.Open(name)
	}

	file, err := os.Open(name)

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	reader.Write(strings.Trim(name, "/"), data)

	return cache.Open(strings.Trim(name, "/"))
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

// Comment
func (ctx *TestingReader) Write(name string, data []byte) error {
	ctx.mutex.Lock()
	ctx.cache[name] = data
	ctx.mutex.Unlock()

	return nil
}
