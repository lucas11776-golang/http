package validation

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/lucas11776-golang/orm"
	"github.com/lucas11776-golang/orm/databases/sqlite"
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

	errorMsg := func(err string) string {
		return strings.ToUpper(err[:1]) + err[1:]
	}

	testValidator := func(validator *Validator, valid bool, errors Errors) {
		if passed := validator.Validate(); passed != valid {
			fmt.Printf("Expected validate to be %t but got %t\r\n", valid, passed)
			t.Fatal("")
		}

		for k, v := range errors {
			if err := validator.Error(k); err != v {
				fmt.Printf("Expected %s to be (%s) but got (%s)\r\n", k, v, err)
				t.Fatal("")
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

	t.Run("TestConfirmed", func(t *testing.T) {
		request, validator := validation(RulesBag{
			"password": Rules{"confirmed"},
		})

		// Fail
		request.Form.Set("password", "test@123")

		testValidator(validator, false, Errors{
			"password": errorMsg(fmt.Sprintf(ConfirmedErrorMessage.Value, "password", "password")),
		})

		// Pass
		request.Form.Set("password_confirmation", "test@123")

		testValidator(validator.Reset(), true, Errors{})
	})

	t.Run("TestEmail", func(t *testing.T) {
		request, validator := validation(RulesBag{
			"email": Rules{"email"},
		})

		// Fail
		request.Form.Set("email", "jane#doe.com")

		testValidator(validator, false, Errors{
			"email": errorMsg(fmt.Sprintf(EmailErrorMessage.Value, "email")),
		})

		// Pass
		request.Form.Set("email", "jane@doe.com")

		testValidator(validator.Reset(), true, Errors{})
	})

	t.Run("TestExists", func(t *testing.T) {
		type User struct {
			Connection string    `json:"-" connection:"sqlite"`
			ID         int64     `json:"id" column:"id" type:"primary_key"`
			CreatedAt  time.Time `json:"created_at" column:"created_at" type:"datetime_current"`
			Email      string    `json:"email" column:"email" type:"string"`
		}

		orm.DB.Add("sqlite", sqlite.Connect(":memory:"))

		db := orm.DB.Database("sqlite").Migration()

		if err := db.Migrate(orm.Models{User{}}); err != nil {
			t.Fatal(err)
		}

		user, err := orm.Model(User{}).Insert(orm.Values{"email": "jeo@doe.com"})

		if err != nil {
			t.Fatal(err)
		}

		request, validator := validation(RulesBag{
			"email": Rules{"exists:users,sqlite"},
		})

		// // Pass
		request.Form.Set("email", "jane@deo.com")

		testValidator(validator.Reset(), true, Errors{})

		// Fail
		request.Form.Set("email", user.Email)

		testValidator(validator, false, Errors{
			"email": errorMsg(fmt.Sprintf(ExistsErrorMessage.Value, "email", "users")),
		})

		orm.DB.Remove("sqlite")
	})

}
