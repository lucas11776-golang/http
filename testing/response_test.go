package testing

import (
	"testing"

	"github.com/lucas11776-golang/http"
	"github.com/lucas11776-golang/http/types"
)

// Comment
func TestResponseBody(t *testing.T) {
	req := NewRequest(&TestCase{
		HTTP:    http.Server("127.0.0.1", 0),
		Testing: t,
	})

	headers := types.Headers{"content-type": "application/json"}

	body := `{"message": "User has been created successfully"}`
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, headers, []byte(body))

	res := InitResponse(req, r, true)

	t.Run("TestAssertProtocol", func(t *testing.T) {
		res.AssertProtocol("HTTP/2.0")

		if !res.Testing.hasError() {
			t.Fatalf("Expected assert protocol log error")
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
			t.Fatalf("Expected assert status code log error")
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
			t.Fatalf("Expected assert has header log error")
		}

		res.Testing.popError()

		res.AssertHasHeader("content-type")

		if res.Testing.hasError() {
			t.Fatalf("Expected assert has header to not log error")
		}

		// header assert
		res.AssertHeader("content-type", "text/html")

		if !res.Testing.hasError() {
			t.Fatalf("Expected assert header log error")
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
			t.Fatalf("Expected assert body log error")
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
	req := NewRequest(&TestCase{
		HTTP:    http.Server("127.0.0.1", 0),
		Testing: t,
	})

	headers := types.Headers{"content-type": "application/json"}

	body := `{"message": "User has been created successfully"}`
	r := http.NewResponse("HTTP/1.1", http.HTTP_RESPONSE_OK, headers, []byte(body))

	res := InitResponse(req, r, true)

	res.AssertIsRedirect()

	if !res.Testing.hasError() {
		t.Fatalf("Expected assert body log error")
	}

	res.Testing.popError()

	res.AssertIsRedirect()

	res.Testing.popError()

	if res.Testing.hasError() {
		t.Fatalf("Expected assert body to not log error")
	}

}

func TestResponseView(t *testing.T) {

}

func TestResponseSession(t *testing.T) {

}
