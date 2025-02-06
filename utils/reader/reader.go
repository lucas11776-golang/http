package reader

import (
	"io/fs"

	"github.com/open2b/scriggo"
)

type CacheReader interface {
	Open(name string) (fs.File, error)
	Cache(name string) (scriggo.Files, error)
}

// Comment
func ReadCache(reader CacheReader, cache scriggo.Files, view string) (scriggo.Files, error) {
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
