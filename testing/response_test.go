package testing

import (
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

		if !res.Testing.hasError() {
			t.Fatalf("Expected assert protocol to log error")
		}

		res.Testing.popError()

		res.AssertProtocol("HTTP/1.1")

		if res.Testing.hasError() {
			t.Fatalf("Expected assert protocol to not log error")
		}
	})

	t.Run("TestAssertStatus", func(t *testing.T) {
		res.AssertStatusCode(http.HTTP_RESPONSE_NOT_FOUND)

		if !res.Testing.hasError() {
			t.Fatalf("Expected assert status code to log error")
		}

		res.Testing.popError()

		res.AssertStatusCode(http.HTTP_RESPONSE_OK)

		if res.Testing.hasError() {
			t.Fatalf("Expected assert status code to not log error")
		}
	})

	t.Run("TestAssertHeaders", func(t *testing.T) {
		res.AssertHasHeader("x-forward")

		if !res.Testing.hasError() {
			t.Fatalf("Expected assert has header to log error")
		}

		res.Testing.popError()

		res.AssertHasHeader("content-type")

		if res.Testing.hasError() {
			t.Fatalf("Expected assert has header to not log error")
		}

		// header assert
		res.AssertHeader("content-type", "text/html")

		if !res.Testing.hasError() {
			t.Fatalf("Expected assert header to log error")
		}

		res.Testing.popError()

		res.AssertHeader("content-type", "application/json")

		if res.Testing.hasError() {
			t.Fatalf("Expected assert header to not log error")
		}
	})

	t.Run("TestBody", func(t *testing.T) {
		res.AssertBody([]byte("Testing Fails Body"))

		if !res.Testing.hasError() {
			t.Fatalf("Expected assert body to log error")
		}

		res.Testing.popError()

		res.AssertBody([]byte(body))

		res.Testing.popError()

		if res.Testing.hasError() {
			t.Fatalf("Expected assert body to not log error")
		}
	})

	req.TestCase.Cleanup()
}

func TestResponseRedirect(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), true))
	headers := types.Headers{"content-type": "text/html"}
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, headers, []byte{})
	res := NewResponse(req, r)

	// Is Redirect
	res.AssertIsRedirect()

	if !res.Testing.hasError() {
		t.Fatalf("Expected assert is redirect to log error")
	}

	res.Testing.popError()

	bag := &http.RedirectBag{To: "authentication/login"}

	res.Response.Bag.Redirect = bag

	res.AssertIsRedirect()

	if res.Testing.hasError() {
		t.Fatalf("Expected assert is redirect to not log error")
	}

	// Redirect To
	res.AssertRedirectTo("dashboard")

	if !res.Testing.hasError() {
		t.Fatalf("Expected assert redirect to log error")
	}

	res.Testing.popError()

	res.AssertRedirectTo(bag.To)

	if res.Testing.hasError() {
		t.Fatalf("Expected assert is redirect to not log error")
	}

	req.TestCase.Cleanup()
}

func TestResponseView(t *testing.T) {
	req := NewRequest(NewTestCase(t, http.Server("127.0.0.1", 0), true))
	headers := types.Headers{"content-type": "text/html"}
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, headers, []byte{})
	res := NewResponse(req, r)

	// Is View
	res.AssertIsView()

	if !res.Testing.hasError() {
		t.Fatalf("Expected assert is view to log error")
	}

	res.Testing.popError()

	bag := &http.ViewBag{
		Name: "authentication.password.reset",
		Data: http.ViewData{
			"error": "Reset password token has expired",
		},
	}

	res.Response.Bag.View = bag

	res.AssertIsView()

	if res.Testing.hasError() {
		t.Fatalf("Expected assert is view to not log error")
	}

	// View
	res.AssertView("dashboard")

	if !res.Testing.hasError() {
		t.Fatalf("Expected assert view to log error")
	}

	res.Testing.popError()

	res.AssertView(bag.Name)

	if res.Testing.hasError() {
		t.Fatalf("Expected assert is view to not log error")
	}

	// View Has
	res.AssertViewHas([]string{"message"})

	if res.Testing.hasError() {
		t.Fatalf("Expected assert view has to not log error")
	}

	res.Testing.popError()

	res.AssertViewHas([]string{"error"})

	if res.Testing.hasError() {
		t.Fatalf("Expected assert view has to not log error")
	}

	req.TestCase.Cleanup()
}

func TestResponseSession(t *testing.T) {

}
