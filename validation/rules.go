package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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
		Value: "the %s already exists in %s",
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

/********************************** Required **********************************/
type Required struct{}

// Comment
func (ctx *Required) Validate(validator *Validator, field string, value interface{}, args ...string) error {
	return CallRuleValidation(
		field,
		value,
		RequiredErrorMessage,
		&TypeValidation{
			Value: func() bool { return value.(string) != "" },
			File:  func() bool { return value.(string) != "" },
		},
		args...,
	)
}

/********************************** Minimum **********************************/
type Minimum struct{}

// Comment
func (ctx *Minimum) Validate(validator *Validator, field string, value interface{}, args ...string) error {
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

/********************************** Maximum **********************************/
type Maximum struct{}

// Comment
func (ctx *Maximum) Validate(validator *Validator, field string, value interface{}, args ...string) error {
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

/********************************** Confirmed **********************************/
type Confirmed struct{}

// Comment
func (ctx *Confirmed) Validate(validator *Validator, field string, value interface{}, args ...string) error {
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

/********************************** Email **********************************/
type Email struct{}

// Comment
func (ctx *Email) Validate(validator *Validator, field string, value interface{}, args ...string) error {
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

/********************************** Exists **********************************/
type Exists struct{}

// Comment
func (ctx *Exists) Validate(validator *Validator, field string, value interface{}, args ...string) error {
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

				return count == 0
			},
			File: func() bool { return false },
		},
		append(args, field)...,
	)
}

var rules = map[string]RuleValidation{
	"required":  &Required{},
	"min":       &Minimum{},
	"max":       &Maximum{},
	"confirmed": &Confirmed{},
	"email":     &Email{},
	"exists":    &Exists{},
}

// Comment
func AddRule(name string, rule RuleValidation) {
	rules[name] = rule
}
