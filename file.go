package http

import (
	"io"
	"mime/multipart"
)

type File struct {
	header *multipart.FileHeader
	file   multipart.File
}

func NewFile(file multipart.File, header *multipart.FileHeader) *File {
	return &File{
		file:   file,
		header: header,
	}
}

// Comment
func (ctx *File) Name() string {
	return ctx.header.Filename
}

// Comment
func (ctx *File) Mime() string {
	return ctx.header.Header.Get("Content-Type")
}

// Comment
func (ctx *File) Read() ([]byte, error) {
	return io.ReadAll(ctx.file)
}

// Comment
func (ctx *File) Size() int64 {
	return ctx.header.Size
}
