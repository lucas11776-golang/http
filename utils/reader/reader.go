package reader

import (
	"io"
	"io/fs"
	"sync"

	"github.com/open2b/scriggo"
)

type TestingReader struct {
	files scriggo.Files
	cache scriggo.Files
	mutex sync.Mutex
}

type CacheReader interface {
	Open(name string) (fs.File, error)
	Cache(name string) (scriggo.Files, error)
	Write(name string, data []byte) error
}

// Comment
func ReadCache(reader CacheReader, cache scriggo.Files, name string) (scriggo.Files, error) {
	if _, ok := cache[name]; ok {
		return cache, nil
	}

	file, err := reader.Open(name)

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	reader.Write(name, data)

	return cache, nil
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
func (ctx *TestingReader) Cache(name string) (scriggo.Files, error) {
	return ReadCache(ctx, ctx.cache, name)
}

// Comment
func (ctx *TestingReader) Write(name string, data []byte) error {
	ctx.mutex.Lock()
	ctx.cache[name] = data
	ctx.mutex.Unlock()

	return nil
}
