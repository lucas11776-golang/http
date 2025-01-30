package types

type Headers map[string]string

type Query map[string]string

type File struct {
	Name string
	Type string
	Data []byte
}
