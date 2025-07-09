package testing

import (
	"strconv"
	"testing"
	"time"

	"github.com/lucas11776-golang/http"
	"github.com/lucas11776-golang/http/types"
	"github.com/lucas11776-golang/orm"
	"github.com/lucas11776-golang/orm/databases/sqlite"
)

// Comment
func TestResponseBody(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), true))
	headers := types.Headers{"content-type": "application/json"}
	body := `{"message": "User has been created successfully"}`
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, headers, []byte(body))

	res := NewResponse(req, r)

	t.Run("TestAssertProtocol", func(t *testing.T) {
		res.AssertProtocol("HTTP/2.0")

		if !res.testing.hasError() {
			t.Fatalf("Expected assert protocol to log error")
		}

		res.testing.popError()

		res.AssertProtocol("HTTP/1.1")

		if res.testing.hasError() {
			t.Fatalf("Expected assert protocol to not log error")
		}

		res.testcase.Cleanup()
	})

	t.Run("TestAssertStatus", func(t *testing.T) {
		res.AssertStatusCode(http.HTTP_RESPONSE_NOT_FOUND)

		if !res.testing.hasError() {
			t.Fatalf("Expected assert status code to log error")
		}

		res.testing.popError()

		res.AssertStatusCode(http.HTTP_RESPONSE_OK)

		if res.testing.hasError() {
			t.Fatalf("Expected assert status code to not log error")
		}

		res.testcase.Cleanup()
	})

	t.Run("TestAssertHeaders", func(t *testing.T) {
		res.AssertHeadersHas("x-forward")

		if !res.testing.hasError() {
			t.Fatalf("Expected assert has header to log error")
		}

		res.testing.popError()

		res.AssertHeadersHas("content-type")

		if res.testing.hasError() {
			t.Fatalf("Expected assert has header to not log error")
		}

		// header assert
		res.AssertHeader("content-type", "text/html")

		if !res.testing.hasError() {
			t.Fatalf("Expected assert header to log error")
		}

		res.testing.popError()

		res.AssertHeader("content-type", "application/json")

		if res.testing.hasError() {
			t.Fatalf("Expected assert header to not log error")
		}

		res.testcase.Cleanup()
	})

	t.Run("TestBody", func(t *testing.T) {
		res.AssertBody([]byte("Testing Fails Body"))

		if !res.testing.hasError() {
			t.Fatalf("Expected assert body to log error")
		}

		res.testing.popError()

		res.AssertBody([]byte(body))

		res.testing.popError()

		if res.testing.hasError() {
			t.Fatalf("Expected assert body to not log error")
		}
	})

	res.testcase.Cleanup()
}

func TestResponseRedirect(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), true))
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, make(types.Headers), []byte{})
	res := NewResponse(req, r)

	// Is Redirect
	res.AssertIsRedirect()

	if !res.testing.hasError() {
		t.Fatalf("Expected assert is redirect to log error")
	}

	res.testing.popError()

	bag := &http.RedirectBag{To: "authentication/login"}

	res.Response.Bag.Redirect = bag

	res.AssertIsRedirect()

	if res.testing.hasError() {
		t.Fatalf("Expected assert is redirect to not log error")
	}

	// Redirect To
	res.AssertRedirectTo("dashboard")

	if !res.testing.hasError() {
		t.Fatalf("Expected assert redirect to log error")
	}

	res.testing.popError()

	res.AssertRedirectTo(bag.To)

	if res.testing.hasError() {
		t.Fatalf("Expected assert is redirect to not log error")
	}

	res.testcase.Cleanup()
}

func TestResponseView(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), true))
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, make(types.Headers), []byte{})
	res := NewResponse(req, r)

	// Is View
	res.AssertIsView()

	if !res.testing.hasError() {
		t.Fatalf("Expected assert is view to log error")
	}

	res.testing.popError()

	bag := &http.ViewBag{
		Name: "authentication.password.reset",
		Data: http.ViewData{
			"error": "Reset password token has expired",
		},
	}

	res.Response.Bag.View = bag

	res.AssertIsView()

	if res.testing.hasError() {
		t.Fatalf("Expected assert is view to not log error")
	}

	// View
	res.AssertView("dashboard")

	if !res.testing.hasError() {
		t.Fatalf("Expected assert view to log error")
	}

	res.testing.popError()

	res.AssertView(bag.Name)

	if res.testing.hasError() {
		t.Fatalf("Expected assert is view to not log error")
	}

	// View Has
	res.AssertViewHas([]string{"message"})

	if !res.testing.hasError() {
		t.Fatalf("Expected assert view has to not log error")
	}

	res.testing.popError()

	res.AssertViewHas([]string{"error"})

	if res.testing.hasError() {
		t.Fatalf("Expected assert view has to not log error")
	}

	res.testcase.Cleanup()
}

func TestResponseSession(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), true))
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, make(types.Headers), []byte{})
	res := NewResponse(req, r)

	req.request.Response = res.Response
	res.request.request = req.request

	req.request.Session = req.testCase.http.Get("session").(http.SessionsManager).Session(req.request)
	res.Response.Session = req.request.Session

	// Has Session
	res.AssertSessionHas([]string{"user_id"})

	if !res.testing.hasError() {
		t.Fatalf("Expected assert has session to log error")
	}

	res.testing.popError()

	res.Response.Session.Set("user_id", strconv.Itoa(1))

	res.AssertSessionHas([]string{"user_id"})

	if res.testing.hasError() {
		t.Fatalf("Expected assert has session to not log error")
	}

	// Session
	res.AssertSession("is_admin", "1")

	if !res.testing.hasError() {
		t.Fatalf("Expected assert session to log error")
	}

	res.testing.popError()

	res.Response.Session.Set("is_admin", strconv.Itoa(1))

	res.AssertSession("is_admin", "1")

	if res.testing.hasError() {
		t.Fatalf("Expected assert session to not log error")
	}

	res.testcase.Cleanup()
}

func TestResponseSessionErrorsHas(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), true))
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, make(types.Headers), []byte{})
	res := NewResponse(req, r)

	req.request.Response = res.Response
	res.request.request = req.request

	req.request.Session = req.testCase.http.Get("session").(http.SessionsManager).Session(req.request)
	res.Response.Session = req.request.Session

	// Has Session errors
	res.AssertSessionErrorsHas([]string{"first_name"})

	if !res.testing.hasError() {
		t.Fatalf("Expected assert has session to log error")
	}

	res.testing.popError()

	res.Response.Session.SetError("email", "The email is required")

	res.AssertSessionErrorsHas([]string{"email"})

	if res.testing.hasError() {
		t.Fatalf("Expected assert has session to not log error")
	}

	res.testcase.Cleanup()
}

func TestResponseSessionError(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), true))
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, make(types.Headers), []byte{})
	res := NewResponse(req, r)

	req.request.Response = res.Response
	res.request.request = req.request

	req.request.Session = req.testCase.http.Get("session").(http.SessionsManager).Session(req.request)
	res.Response.Session = req.request.Session

	// Has Session error
	res.AssertSessionError("email", "The email is required")

	if !res.testing.hasError() {
		t.Fatalf("Expected assert has session to log error")
	}

	res.testing.popError()

	res.Response.Session.SetError("first_name", "The first name is required")

	res.AssertSessionError("first_name", "The first name is required")

	if res.testing.hasError() {
		t.Fatalf("Expected assert has session to not log error")
	}

	res.testcase.Cleanup()
}

func TestResponseJsonErrorsHas(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), true))
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, make(types.Headers), []byte{})
	res := NewResponse(req, r)

	req.request.Response = res.Response
	res.request.request = req.request

	res.Response.SetStatus(http.HTTP_RESPONSE_UNPROCESSABLE_CONTENT).Json(http.JsonErrorResponse{
		Message: "Form validation error",
		Errors: http.SessionErrorsBag{
			"email":    "The email is required",
			"password": "The password is required",
		},
	})

	// JSON error
	res.AssertJsonErrorsHas([]string{"first_name"})

	if !res.testing.hasError() {
		t.Fatalf("Expected assert has json error to log error")
	}

	res.testing.popError()

	res.AssertJsonErrorsHas([]string{"email"})

	if res.testing.hasError() {
		t.Fatalf("Expected assert has json error to not log error")
	}

	res.testcase.Cleanup()
}

func TestResponseJsonError(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), true))
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, make(types.Headers), []byte{})
	res := NewResponse(req, r)

	req.request.Response = res.Response
	res.request.request = req.request

	res.Response.SetStatus(http.HTTP_RESPONSE_UNPROCESSABLE_CONTENT).Json(http.JsonErrorResponse{
		Message: "Form validation error",
		Errors: http.SessionErrorsBag{
			"email":    "The email is required",
			"password": "The password is required",
		},
	})

	// JSON error
	res.AssertJsonError("first_name", "The first name is required")

	if !res.testing.hasError() {
		t.Fatalf("Expected assert has json error to log error")
	}

	res.testing.popError()

	res.AssertJsonError("password", "The password is required")

	if res.testing.hasError() {
		t.Fatalf("Expected assert has json error to not log error")
	}

	res.testcase.Cleanup()
}

func TestAssertDatabaseHas(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), true))
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, make(types.Headers), []byte{})
	res := NewResponse(req, r)

	req.request.Response = res.Response
	res.request.request = req.request

	type User struct {
		Connection string    `json:"-" connection:"sqlite" table:"users"`
		ID         int64     `json:"id" column:"id" type:"primary_key"`
		CreatedAt  time.Time `json:"created_at" column:"created_at" type:"datetime_current"`
		Email      string    `json:"email" column:"email" type:"string"`
		Name       string    `json:"name" column:"name" type:"string"`
	}

	db := sqlite.Connect(":memory:")

	orm.DB.Add("sqlite", db)

	if err := db.Migration().Migrate(orm.Models{User{}}); err != nil {
		t.Fatal(err)
	}

	user := &User{
		Email: "jeo@doe.com",
		Name:  "Jeo Deo",
	}

	_, err := db.Insert(&orm.Statement{
		Table: "users",
		Values: orm.Values{
			"email": user.Email,
			"name":  user.Name,
		},
		PrimaryKey: "id",
	})

	if err != nil {
		t.Fatal(err)
	}

	res.AssertDatabaseHas("sqlite", "users", map[string]interface{}{
		"email": "jane@doe.com",
		"name":  "Jane Deo",
	})

	if !res.testing.hasError() {
		t.Fatalf("Expected assert has database to log error")
	}

	res.testing.popError()

	res.AssertDatabaseHas("sqlite", "users", map[string]interface{}{
		"email": user.Email,
		"name":  user.Name,
	})

	if res.testing.hasError() {
		t.Fatalf("Expected assert has database to not log error")
	}

	orm.DB.Remove("sqlite")

	res.testcase.Cleanup()
}

func TestAssertDatabaseMissing(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), true))
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, make(types.Headers), []byte{})
	res := NewResponse(req, r)

	req.request.Response = res.Response
	res.request.request = req.request

	type User struct {
		Connection string    `json:"-" connection:"sqlite" table:"users"`
		ID         int64     `json:"id" column:"id" type:"primary_key"`
		CreatedAt  time.Time `json:"created_at" column:"created_at" type:"datetime_current"`
		Email      string    `json:"email" column:"email" type:"string"`
		Name       string    `json:"name" column:"name" type:"string"`
	}

	db := sqlite.Connect(":memory:")

	orm.DB.Add("sqlite", db)

	if err := db.Migration().Migrate(orm.Models{User{}}); err != nil {
		t.Fatal(err)
	}

	user := &User{
		Email: "john@doe.com",
		Name:  "John Deo",
	}

	_, err := db.Insert(&orm.Statement{
		Table: "users",
		Values: orm.Values{
			"email": user.Email,
			"name":  user.Name,
		},
		PrimaryKey: "id",
	})

	if err != nil {
		t.Fatal(err)
	}

	res.AssertDatabaseMissing("sqlite", "users", map[string]interface{}{
		"email": user.Email,
		"name":  user.Name,
	})

	if !res.testing.hasError() {
		t.Fatalf("Expected assert has database to log error")
	}

	res.testing.popError()

	res.AssertDatabaseMissing("sqlite", "users", map[string]interface{}{
		"email": "jane@deo.com",
		"name":  "Jane",
	})

	if res.testing.hasError() {
		t.Fatalf("Expected assert has database to not log error")
	}

	orm.DB.Remove("sqlite")

	res.testcase.Cleanup()
}

func TestAssertDatabaseCount(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), true))
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, make(types.Headers), []byte{})
	res := NewResponse(req, r)

	req.request.Response = res.Response
	res.request.request = req.request

	type User struct {
		Connection string    `json:"-" connection:"sqlite" table:"users"`
		ID         int64     `json:"id" column:"id" type:"primary_key"`
		CreatedAt  time.Time `json:"created_at" column:"created_at" type:"datetime_current"`
		Email      string    `json:"email" column:"email" type:"string"`
		Name       string    `json:"name" column:"name" type:"string"`
	}

	db := sqlite.Connect(":memory:")

	orm.DB.Add("sqlite", db)

	if err := db.Migration().Migrate(orm.Models{User{}}); err != nil {
		t.Fatal(err)
	}

	res.AssertDatabaseCount("sqlite", "users", 0)

	if res.testing.hasError() {
		t.Fatalf("Expected assert has database to log error")
	}

	res.testing.popError()

	_, err := db.Insert(&orm.Statement{
		Table: "users",
		Values: orm.Values{
			"email": "jeo@deo.com",
			"name":  "Jeo Deo",
		},
		PrimaryKey: "id",
	})

	if err != nil {
		t.Fatal(err)
	}

	res.AssertDatabaseCount("sqlite", "users", 0)

	if !res.testing.hasError() {
		t.Fatalf("Expected assert has database to log error")
	}

	res.testing.popError()

	res.AssertDatabaseCount("sqlite", "users", 1)

	if res.testing.hasError() {
		t.Fatalf("Expected assert has database to not log error")
	}

	orm.DB.Remove("sqlite")

	res.testcase.Cleanup()
}
