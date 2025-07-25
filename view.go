package http

import (
	"io/fs"
	"os"
	"strings"

	"github.com/lucas11776-golang/http/utils/helper"
	"github.com/lucas11776-golang/http/utils/path"
	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/native"
)

type ViewData map[string]interface{}

type ViewWriter struct {
	parsed []byte
}

type View struct {
	fs           fs.FS
	extension    string
	declarations []native.Declarations
}

type ViewInterface interface {
	Read(view string, data ViewData) ([]byte, error)
}

// Comment
func (ctx *ViewWriter) Write(p []byte) (n int, err error) {
	ctx.parsed = append(ctx.parsed, p...)

	return len(ctx.parsed), nil
}

// Comment
func (ctx *ViewWriter) Parsed() []byte {
	return ctx.parsed
}

type DefaultViewReader struct {
	fs.FS
}

// Comment
func NewDefaultViewReader(views string) fs.FS {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return &DefaultViewReader{os.DirFS(path.Path(wd, views))}
}

// Comment
func NewView(fs fs.FS, extension string, declarations ...native.Declarations) *View {
	return &View{
		fs:           fs,
		extension:    extension,
		declarations: declarations,
	}
}

// Comment
func (ctx *View) buildTemplate(view string, options *scriggo.BuildOptions) (*scriggo.Template, error) {
	return scriggo.BuildTemplate(ctx.fs, path.FileRealPath(view, ctx.extension), options)
}

// Comment
func viewDeclarationsWithHelpers(req *Request) native.Declarations {
	return native.Declarations{
		"url":           helper.Url,
		"subdomain":     helper.Subdomain,
		"format":        helper.Format,
		"session":       SessionValue(req),
		"has":           SessionHas(req),
		"error":         SessionError(req),
		"errors":        SessionErrors(req),
		"csrfName":      SessionCsrfName(req),
		"csrfToken":     SessionCsrfToken(req),
		"old":           SessionOld(req),
		"methodName":    func() string { return RequestFormMethodName },
		"request":       func() *Request { return req },
		"replace":       strings.ReplaceAll,
		"truncate":      helper.Truncate,
		"cast":          func() *helper.Cast { return &helper.Cast{} },
		"queryToString": helper.QueryToString,
		"current":       func() string { return helper.Url(req.Path()) }, // TODO: create url cast e.g url().Current(), url().To("login")...
	}
}

// Comment
func (ctx *View) Read(view string, data ViewData, req *Request) ([]byte, error) {
	globals := viewDeclarationsWithHelpers(req)

	for _, declarations := range ctx.declarations {
		for k, v := range declarations {
			globals[k] = v
		}
	}

	for key, value := range data {
		globals[key] = value
	}

	template, err := ctx.buildTemplate(view, &scriggo.BuildOptions{Globals: globals})

	if err != nil {
		return nil, err
	}

	writer := &ViewWriter{}

	if err := template.Run(writer, nil, nil); err != nil {
		return nil, err
	}

	// TODO: Find out what is this for (\r\n\r\n)
	return []byte(strings.ReplaceAll(string(writer.parsed), "\r\n\r\n", "\r\n")), nil
}
