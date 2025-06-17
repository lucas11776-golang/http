package validation

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestValidation(t *testing.T) {
	t.Run("TestRequiredStringRule", func(t *testing.T) {
		request, err := http.NewRequest("POST", "/", strings.NewReader(""))

		if err != nil {
			t.Fatal(err)
		}

		request.Form = url.Values{}

		validator := Validation(request, RulesBag{
			"email": Rules{"required"},
		})

		if valid := validator.Validate(); valid == true {
			t.Fatalf("Expected validate to be (%t) but got (%t)", false, valid)
		}

		expectedErrMsg := "The email is required"
		errMsg := validator.Error("email")

		if errMsg != expectedErrMsg {
			t.Fatalf("Expected email error message to be (%s) but got (%s)", expectedErrMsg, errMsg)
		}

		email := "jane@doe.com"

		validator.request.Form.Set("email", email)

		if valid := validator.Validate(); valid == true {
			t.Fatalf("Expected validate to be (%t) but got (%t)", true, valid)
		}

		fields, _ := validator.Validated()

		if field := fields["email"]; field != email {
			t.Fatalf("Expected email to be (%s) but got (%s)", email, field)
		}
	})

	t.Run("TestRequiredCallbackRule", func(t *testing.T) {
		request, err := http.NewRequest("POST", "/", strings.NewReader(""))

		if err != nil {
			t.Fatal(err)
		}

		validator := Validation(request, RulesBag{
			"email": Rules{&Required{}},
		})

		if validator.Validate() != false {
			t.Fatal("Expected validate to be false but got true")
		}

		expectedErrMsg := "The email is required"
		errMsg := validator.Error("email")

		if errMsg != expectedErrMsg {
			t.Fatalf("Expected email error message to be (%s) but got (%s)", expectedErrMsg, errMsg)
		}

		email := "jane@doe.com"

		request.Form.Set("email", email)

		if valid := validator.Validate(); valid == true {
			t.Fatalf("Expected validate to be (%t) but got (%t)", true, valid)
		}
	})

}
