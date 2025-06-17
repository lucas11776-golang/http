package validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

const (
	REQUIRED       string = "the %s is required"
	MINIMUM_STRING string = "the %s must have minimum length of %s character"
	MINIMUM_FILE   string = ""
	MAXIMUM_STRING string = "the %s must have maximum length of %s character"
	MAXIMUM_FILE   string = ""
)

// TODO: refactor to this will all fun(valid bool, message ErrorMessage, params MessageErrorParameters) error
type ErrorMessage struct {
	String string
	File   string
}

type MessageErrorParameters map[string]string

var (
	MinimumErrorMessage *ErrorMessage = &ErrorMessage{
		String: "the :name must have minimum length of :value character",
		File:   "",
	}
)

var (
	ErrValueNotSupport = errors.New("the value is not support")
)

// Comment
func FormatName(field string) string {
	return strings.ReplaceAll(field, "_", " ")
}

type Required struct{}

// Comment
func (ctx *Required) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	if value == nil {
		return fmt.Errorf("the %s is required", FormatName(field))
	}

	switch value.(type) {
	case string:
		if value == "" {
			return fmt.Errorf("the %s is required", FormatName(field))
		}

	case *File:
		// TODO: add file

	default:
		return ErrValueNotSupport
	}

	return nil
}

type MinimumLenght struct{}

// Comment
func (ctx *MinimumLenght) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	if len(args) != 1 {
		return errors.New("minimun lenght requires at least one argument of size")
	}

	switch value.(type) {
	case string:
		if len(value.(string)) < cast.ToInt(args[0]) {
			return fmt.Errorf(MINIMUM_STRING, FormatName(field), args[0])
		}
	case *File:
		// TODO: add file

	default:
		return ErrValueNotSupport
	}

	return nil
}

type MaximumLenght struct{}

// Comment
func (ctx *MaximumLenght) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	if len(args) != 1 {
		return errors.New("minimun lenght requires at least one argument of size")
	}

	if value == nil {
		return fmt.Errorf(MAXIMUM_STRING, FormatName(field), args[0])
	}

	switch value.(type) {
	case string:
		if len(value.(string)) > cast.ToInt(args[0]) {
			return fmt.Errorf(MAXIMUM_STRING, FormatName(field), args[0])
		}
	case *File:
		// TODO: add file

	default:
		return ErrValueNotSupport
	}

	return nil
}

var ValidatorsRules = map[string]RuleValidation{
	"required": &Required{},
	"min":      &MinimumLenght{},
	"max":      &MaximumLenght{},
}
