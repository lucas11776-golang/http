package validation

import "net/http"

type Validator struct {
}

// Comment
func Validation(req *http.Request, rules any) *Validator {
	return &Validator{}
}
