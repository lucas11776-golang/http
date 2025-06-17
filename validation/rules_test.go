package validation

import (
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestRules(t *testing.T) {
	validation := func(bag RulesBag) (*http.Request, *Validator) {
		request, err := http.NewRequest("POST", "/", strings.NewReader(""))

		request.Form = url.Values{}
		request.MultipartForm = &multipart.Form{}

		if err != nil {
			t.Fatal(err)
		}

		return request, Validation(request, bag)
	}

	testValidator := func(validator *Validator, valid bool, errors Errors) {
		if passed := validator.Validate(); passed != valid {
			t.Fatalf("Expected validate to be %t but got %t", valid, passed)
		}

		for k, v := range errors {
			if err := validator.Error(k); err != v {
				t.Fatalf("Expected %s to be (%s) but got (%s)", k, v, err)
			}
		}

	}

	t.Run("TestRequired", func(t *testing.T) {
		request, validator := validation(RulesBag{
			"email": Rules{"required"},
		})

		// Fail
		testValidator(validator, false, Errors{
			"email": "The email is required",
		})

		// Pass
		request.Form.Set("email", "jeo@doe.com")

		testValidator(validator.Reset(), true, Errors{})
	})

	t.Run("TestMin", func(t *testing.T) {
		request, validator := validation(RulesBag{
			"first_name": Rules{"min:3"},
		})

		// Fail
		request.Form.Set("first_name", "Je")

		testValidator(validator, false, Errors{
			"first_name": "The first name must have minimum length of 3 character",
		})

		// Pass
		request.Form.Set("first_name", "Jeo")

		testValidator(validator.Reset(), true, Errors{})
	})

	t.Run("TestMax", func(t *testing.T) {
		request, validator := validation(RulesBag{
			"first_name": Rules{"max:5"},
		})

		// Fail
		request.Form.Set("first_name", "Peterson")

		testValidator(validator, false, Errors{
			"first_name": "The first name must have maximum length of 5 character",
		})

		// Pass
		request.Form.Set("first_name", "Jeo")

		testValidator(validator.Reset(), true, Errors{})
	})

}
