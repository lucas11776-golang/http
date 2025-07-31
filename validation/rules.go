package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/lucas11776-golang/orm"
	"github.com/spf13/cast"
)

type ErrorMessage struct {
	Value string
	File  string
}

var (
	ErrValueNotSupport = errors.New("the value is not support")
)

const (
	NullableFlag = ":nullable:"
)

var (
	RequiredErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s is required",
	}
	MinimumErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s must have minimum length of %s character",
		File:  "",
	}
	MaximumErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s must have maximum length of %s character",
		File:  "",
	}
	ConfirmedErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s does not match %s confirmation",
	}
	EmailErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s is invalid",
	}
	ExistsErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s does not exists in %s",
	}
	UniqueErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s already exists in %s",
	}
	DatetimeErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s is invalid datetime",
	}
	DateErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s is invalid date",
	}
	StringErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s is not a string",
	}
	IntegerErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s is not a integer",
	}
	FloatErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s is not a float",
	}
	NumberErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s is not a number",
	}
	AcceptedErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s are accepted",
	}
	NullableErrorMessage *ErrorMessage = &ErrorMessage{
		Value: NullableFlag + "/",
		File:  NullableFlag + "/",
	}
	FileErrorMessage *ErrorMessage = &ErrorMessage{
		Value: "the %s is not a file",
	}
)

// Comment
func FormatName(field string) string {
	return strings.ReplaceAll(field, "_", " ")
}

type ValidateRuleCallback func() bool

type TypeValidation struct {
	Value ValidateRuleCallback
	File  ValidateRuleCallback
}

// Comment
func FormattedErrorMessage(field string, err string, args ...string) string {
	if len(args) != 0 {
		return fmt.Sprintf(err, []any{FormatName(field), args[0]}...)
	}

	return fmt.Sprintf(err, FormatName(field))
}

// Comment
func CallRuleValidation(field string, value interface{}, errorMessage *ErrorMessage, validation *TypeValidation, args ...string) error {
	switch value.(type) {
	case string:
		if validation.Value() {
			return nil
		}

		return errors.New(FormattedErrorMessage(field, errorMessage.Value, args...))

	case *File:
		if validation.File() {
			return nil
		}

		return errors.New(FormattedErrorMessage(field, errorMessage.File, args...))

	default:
		return ErrValueNotSupport
	}
}

/********************************** RequiredRule **********************************/
type RequiredRule struct{}

// Comment
func (ctx *RequiredRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	return CallRuleValidation(
		field,
		value,
		RequiredErrorMessage,
		&TypeValidation{
			Value: func() bool { return value.(string) != "" },
			File:  func() bool { return value != nil },
		},
		args...,
	)
}

/********************************** MinimumRule **********************************/
type MinimumRule struct{}

// Comment
func (ctx *MinimumRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	return CallRuleValidation(
		field,
		value,
		MinimumErrorMessage,
		&TypeValidation{
			Value: func() bool { return len(value.(string)) >= cast.ToInt(args[0]) },
			File:  func() bool { return false },
		},
		args...,
	)
}

/********************************** MaximumRule **********************************/
type MaximumRule struct{}

// Comment
func (ctx *MaximumRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	return CallRuleValidation(
		field,
		value,
		MaximumErrorMessage,
		&TypeValidation{
			Value: func() bool { return len(value.(string)) <= cast.ToInt(args[0]) },
			File:  func() bool { return false },
		},
		args...,
	)
}

/********************************** ConfirmedRule **********************************/
type ConfirmedRule struct{}

// Comment
func (ctx *ConfirmedRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	return CallRuleValidation(
		field,
		value,
		ConfirmedErrorMessage,
		&TypeValidation{
			Value: func() bool { return value.(string) == validator.FormValue(fmt.Sprintf("%s_confirmation", field)) },
			File:  func() bool { return false },
		},
		append(args, field)...,
	)
}

/********************************** EmailRule **********************************/
type EmailRule struct{}

// Comment
func (ctx *EmailRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	return CallRuleValidation(
		field,
		value,
		EmailErrorMessage,
		&TypeValidation{
			Value: func() bool {
				return regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(value.(string))
			},
			File: func() bool { return false },
		},
		args...,
	)
}

/********************************** ExistsRule **********************************/
type ExistsRule struct{}

// Comment
func (ctx *ExistsRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	if len(args) < 2 {
		return errors.New("exists expect at least 2 arguments")
	}

	db := orm.DB.Database(args[1])

	if db == nil {
		return fmt.Errorf("connection %s does not exist in database", args[1])
	}

	if len(args) > 2 {
		field = args[2]
	}

	return CallRuleValidation(
		field,
		value,
		ExistsErrorMessage,
		&TypeValidation{
			Value: func() bool {
				count, err := db.Count(&orm.Statement{
					Table: args[0],
					Where: []interface{}{&orm.Where{
						Key:      field,
						Operator: orm.EQUALS,
						Value:    value,
					}},
				})

				if err != nil {
					return false
				}

				return count != 0
			},
			File: func() bool { return false },
		},
		append(args, field)...,
	)
}

/********************************** Exists **********************************/
type UniqueRule struct{}

// Comment
func (ctx *UniqueRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	if len(args) < 2 {
		return errors.New("exists expect at least 2 arguments")
	}

	db := orm.DB.Database(args[1])

	if db == nil {
		return fmt.Errorf("connection %s does not exist in database", args[1])
	}

	if len(args) > 2 {
		field = args[2]
	}

	return CallRuleValidation(
		field,
		value,
		UniqueErrorMessage,
		&TypeValidation{
			Value: func() bool {
				count, err := db.Count(&orm.Statement{
					Table: args[0],
					Where: []interface{}{&orm.Where{
						Key:      field,
						Operator: orm.EQUALS,
						Value:    value,
					}},
				})

				if err != nil {
					return false
				}

				return count == 0
			},
			File: func() bool { return false },
		},
		append(args, field)...,
	)
}

/********************************** Email **********************************/
type DatetimeRule struct{}

// Comment
func (ctx *DatetimeRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	return CallRuleValidation(
		field,
		value,
		DatetimeErrorMessage,
		&TypeValidation{
			Value: func() bool {
				_, err := time.Parse(time.DateTime, value.(string))

				return err == nil
			},
			File: func() bool { return false },
		},
		args...,
	)
}

/********************************** Email **********************************/
type DateRule struct{}

// Comment
func (ctx *DateRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	return CallRuleValidation(
		field,
		value,
		DateErrorMessage,
		&TypeValidation{
			Value: func() bool {
				_, err := time.Parse(time.DateOnly, value.(string))

				return err == nil
			},
			File: func() bool { return false },
		},
		args...,
	)
}

/********************************** StringRule **********************************/
type StringRule struct{}

// Comment
func (ctx *StringRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	return CallRuleValidation(
		field,
		value,
		StringErrorMessage,
		&TypeValidation{
			Value: func() bool { return true },
			File:  func() bool { return false },
		},
		args...,
	)
}

/********************************** IntegerRule **********************************/
type IntegerRule struct{}

// Comment
func (ctx *IntegerRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	return CallRuleValidation(
		field,
		value,
		IntegerErrorMessage,
		&TypeValidation{
			Value: func() bool { return regexp.MustCompile(`^-?\d+$`).MatchString(value.(string)) },
			File:  func() bool { return false },
		},
		args...,
	)
}

/********************************** FloatRule **********************************/
type FloatRule struct{}

// Comment
func (ctx *FloatRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	return CallRuleValidation(
		field,
		value,
		FloatErrorMessage,
		&TypeValidation{
			Value: func() bool { return regexp.MustCompile(`^[+-]?(\d+\.\d+|\.\d+|\d+\.)$`).MatchString(value.(string)) },
			File:  func() bool { return false },
		},
		args...,
	)
}

/********************************** Float **********************************/
type NumberRule struct{}

// Comment
func (ctx *NumberRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	return CallRuleValidation(
		field,
		value,
		NumberErrorMessage,
		&TypeValidation{
			Value: func() bool {
				return regexp.MustCompile(`^[+-]?(\d+\.\d+|\.\d+|\d+\.)$|^[+-]?\d+$`).MatchString(value.(string))
			},
			File: func() bool { return false },
		},
		args...,
	)
}

/********************************** AcceptedRule **********************************/
type AcceptedRule struct{}

// Comment
func (ctx *AcceptedRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	return CallRuleValidation(
		field,
		value,
		AcceptedErrorMessage,
		&TypeValidation{
			Value: func() bool { return regexp.MustCompile(`^(?i)(yes|on|true|1)$`).MatchString(value.(string)) },
			File:  func() bool { return false },
		},
		args...,
	)
}

/********************************** NullableRule **********************************/
type NullableRule struct{}

// Comment
func (ctx *NullableRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {

	err := CallRuleValidation(
		field,
		value,
		NullableErrorMessage,
		&TypeValidation{
			Value: func() bool { return cast.ToString(value) != "" },
			File:  func() bool { return value != nil },
		},
	)

	if err != nil {
		return errors.New(strings.Split(err.Error(), "/")[0])
	}

	return err
}

/********************************** FileRule **********************************/
type FileRule struct{}

// Comment
func (ctx *FileRule) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	err := CallRuleValidation(
		field,
		value,
		FileErrorMessage,
		&TypeValidation{
			Value: func() bool { return false },
			File:  func() bool { return true },
		},
	)

	if err != nil {
		return errors.New(strings.Split(err.Error(), "/")[0])
	}

	return err
}

// Comment
var rules = map[string]RuleValidation{
	"required":  &RequiredRule{},
	"min":       &MinimumRule{},
	"max":       &MaximumRule{},
	"confirmed": &ConfirmedRule{},
	"email":     &EmailRule{},
	"unique":    &UniqueRule{},
	"exists":    &ExistsRule{},
	"datetime":  &DatetimeRule{},
	"date":      &DateRule{},
	"string":    &StringRule{},
	"integer":   &IntegerRule{},
	"float":     &FloatRule{},
	"number":    &NumberRule{},
	"accepted":  &AcceptedRule{},
	"nullable":  &NullableRule{},
	"file":      &FileRule{},
}

// Comment
func AddRule(name string, rule RuleValidation) {
	rules[name] = rule
}
