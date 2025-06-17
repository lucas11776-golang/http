package types

import (
	"github.com/open2b/scriggo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Headers map[string]string

type Query map[string]string

type File struct {
	Name string
	Type string
	Data []byte
}

type Fs scriggo.Files

// Comment
func (ctx Headers) Get(key string) string {
	header, ok := ctx[cases.Title(language.English).String(key)]

	if !ok {
		return ""
	}

	return header
}
