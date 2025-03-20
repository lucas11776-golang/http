package reader

import (
	"io/fs"

	"github.com/open2b/scriggo"
)

type CacheReader interface {
	Open(name string) (fs.File, error)
	Cache(name string) (scriggo.Files, error)
	Write(name string, data []byte) error
}

// Comment
func ReadCache(reader CacheReader, cache scriggo.Files, name string) (scriggo.Files, error) {
	_, ok := cache[name]

	if ok {
		return cache, nil
	}

	file, err := reader.Open(name)

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

	cache[name] = data

	return cache, nil
}
