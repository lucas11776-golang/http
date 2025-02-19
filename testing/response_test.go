package testing

import (
	"strconv"
	"testing"

	"github.com/lucas11776-golang/http"
	"github.com/lucas11776-golang/http/types"
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
