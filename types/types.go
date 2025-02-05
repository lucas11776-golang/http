package types

import "github.com/open2b/scriggo"

type Headers map[string]string

type Query map[string]string

type File struct {
	Name string
	Type string
	Data []byte
}

type Fs scriggo.Files
