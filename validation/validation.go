package validation

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type File struct {
	header *multipart.FileHeader
	file   multipart.File
}

type Errors map[string]string

type Values map[string]string

type Files map[string]*File

type Data struct {
	Files  Files
	Values Values
}

type Validator struct {
	errors    Errors
	validated *Data
	request   *http.Request
	rules     RulesBag
}

type Rule interface{}

type Rules []Rule

type RulesBag map[string]Rules

// Comment
func (ctx Files) Get(key string) *File {
	file, ok := ctx[key]

	if !ok {
		return nil
	}

	return file
}

type RuleValidation interface {
	Validate(validator *Validator, field string, value interface{}, args ...string) error
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

// Comment
func Validation(req *http.Request, rules RulesBag) *Validator {
	return &Validator{
		errors:  make(Errors),
		request: req,
		rules:   rules,
		validated: &Data{
			Files:  make(Files),
			Values: make(Values),
		},
	}
}

// Comment
func (ctx *Validator) Validated() (fields Values, files Files) {
	return ctx.validated.Values, ctx.validated.Files
}

// Comment
func (ctx *Validator) Values() Values {
	return ctx.validated.Values
}

// Comment
func (ctx *Validator) Value(key string) string {
	value, ok := ctx.validated.Values[key]

	if !ok {
		return ""
	}

	return value
}

// Comment
func (ctx *Validator) Files() Files {
	return ctx.validated.Files
}

// Comment
func (ctx *Validator) File(key string) *File {
	file, ok := ctx.validated.Files[key]

	if !ok {
		return nil
	}

	return file
}

// Comment
func (ctx *Validator) Errors() Errors {
	return ctx.errors
}

// Comment
func (ctx *Validator) Error(key string) string {
	return ctx.errors[key]
}

// Comment
func (ctx *Validator) getValue(key string) interface{} {
	if value := ctx.request.FormValue(key); value != "" {
		return value
	}

	if file, header, err := ctx.request.FormFile(key); err == nil {
		return &File{
			file:   file,
			header: header,
		}
	}

	return nil
}

// Comment
func (ctx *Validator) FormValue(key string) string {
	return ctx.request.FormValue(key)
}

// Comment
func (ctx *Validator) FormFile(key string) *File {
	if file, header, err := ctx.request.FormFile(key); err == nil {
		return &File{
			file:   file,
			header: header,
		}
	}

	return nil
}

// Comment
func (ctx *Validator) addValue(key string, value interface{}) {
	switch value.(type) {
	case *File:
		ctx.validated.Files[key] = value.(*File)

	case string:
		ctx.validated.Values[key] = value.(string)
	}
}

// Comment
func (ctx *Validator) call(callback RuleValidation, field string, args ...string) error {
	value := ctx.getValue(field)

	if value == nil {
		value = ""
	}

	if err := callback.Validate(ctx, field, value, args...); err != nil {
		return err
	}

	ctx.addValue(field, value)

	return nil
}

// Comment
func (ctx *Validator) validate(field string, _rules Rules) error {
	for _, rule := range _rules {

		switch rule.(type) {
		case RuleValidation:
			return ctx.call(rule.(RuleValidation), field)

		case string:
			_field := strings.Split(rule.(string), ":")

			if len(_field) > 2 {
				return fmt.Errorf("invalid args %s", _field[1:])
			}

			_args := []string{}

			if len(_field) == 2 {
				_args = strings.Split(_field[1], ",")
			}

			_rule, ok := rules[_field[0]]

			if !ok {
				return fmt.Errorf("rule %s does not exist", rule)
			}

			return ctx.call(_rule, field, _args...)

		default:
			return fmt.Errorf("rule %s does not exist", rule)
		}
	}

	return nil
}

// Comment
func (ctx *Validator) Validate() bool {
	for field, rules := range ctx.rules {
		if err := ctx.validate(field, rules); err != nil {
			ctx.errors[field] = strings.ToUpper(err.Error()[:1]) + err.Error()[1:]
		}
	}

	return len(ctx.errors) == 0
}

// Comment
func (ctx *Validator) Reset() *Validator {
	ctx.errors = make(Errors)
	ctx.validated = &Data{
		Files:  make(Files),
		Values: make(Values),
	}

	return ctx
}
