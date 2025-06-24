package http

import (
	"bytes"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/lucas11776-golang/http/types"
	str "github.com/lucas11776-golang/http/utils/strings"
	"github.com/spf13/cast"
)

func TestSession(t *testing.T) {
	t.Run("TestGetSession", func(t *testing.T) {
		sessions := InitSession("session", []byte(str.Random(10)))
		userId := strconv.Itoa(int(rand.Float32() * 10000))

		req, err := NewRequest("POST", "/login", "HTTP/1.1", make(types.Headers), bytes.NewReader([]byte{}))

		if err != nil {
			t.Fatalf("Something went wrong when trying to create request: %s", err.Error())
		}

		session := sessions.Session(req)

		session.Set("user_id", userId).Save()

		cookie, err := url.ParseQuery(strings.ReplaceAll(req.Response.GetHeader("set-cookie"), "; ", "&"))

		if err != nil {
			t.Fatalf("Something went wrong when trying to parse cookie to query: %s", err)
		}

		sessionKey := cookie.Get("session")

		headers := types.Headers{
			"Cookie": strings.Join([]string{"session", sessionKey}, "="),
		}

		reqRepeat, err := NewRequest("GET", "/dashboard", "HTTP/1.1", headers, strings.NewReader(""))

		if err != nil {
			t.Fatalf("Something went wrong when trying to create request: %s", err.Error())
		}

		session = sessions.Session(reqRepeat)

		if session.Get("user_id") != userId {
			t.Fatalf("Expected value session_id from user_id to be (%s) but got (%s)", userId, session.Get("user_id"))
		}
	})

	t.Run("TestSessionsManager", func(t *testing.T) {
		domain := "map.guarded.com"
		maxAge := ((60 * 60) * 24) * 5
		secure := true
		httpOnly := true
		sameSite := true
		path := "/authentication"

		sessions := InitSession("session", []byte(str.Random(10)))

		sessions.Domain(domain).MaxAge(maxAge).Secure(secure).HttpOnly(httpOnly).SameSite(sameSite).Path(path)

		if sessions.store.Options.Domain != domain {
			t.Fatalf("Expected the domain to be (%s) but got (%s)", domain, sessions.store.Options.Domain)
		}

		if sessions.store.Options.MaxAge != maxAge {
			t.Fatalf("Expected the max age to be (%d) but got (%d)", maxAge, sessions.store.Options.MaxAge)
		}

		if sessions.store.Options.Secure != secure {
			t.Fatalf("Expected the secure to be (%t) but got (%t)", secure, sessions.store.Options.Secure)
		}

		if sessions.store.Options.HttpOnly != httpOnly {
			t.Fatalf("Expected the http only to be (%t) but got (%t)", httpOnly, sessions.store.Options.HttpOnly)
		}

		if sessions.store.Options.SameSite != 1 {
			t.Fatalf("Expected the same site to be (%d) but got (%d)", 1, sessions.store.Options.SameSite)
		}

		if sessions.store.Options.Path != path {
			t.Fatalf("Expected the path to be (%s) but got (%s)", path, sessions.store.Options.Path)
		}
	})

	t.Run("TestSessionManager", func(t *testing.T) {
		userId := "1"
		userRole := "1"
		path := "dashboard"
		sessions := InitSession("session", []byte(str.Random(10)))

		req, err := NewRequest("POST", "/", "HTTP/1.1", make(types.Headers), bytes.NewReader([]byte{}))

		if err != nil {
			t.Fatalf("Something went wrong when trying to create request: %s", err.Error())
		}

		session := sessions.Session(req).Path(path)

		session.Set("user_id", userId).Set("role", userRole).Save()

		cookie, err := url.ParseQuery(strings.ReplaceAll(req.Response.GetHeader("Set-Cookie"), "; ", "&"))

		if err != nil {
			t.Fatalf("Something went wrong when trying to parse cookie to query: %s", err)
		}

		if cookie.Get("Path") != path {
			t.Fatalf("Expected the path to be (%s) but got (%s)", path, cookie.Get("Path"))
		}

		// Second Request
		headers := types.Headers{
			"cookie": strings.Join([]string{"session", cookie.Get("session")}, "="),
		}

		req, err = NewRequest("POST", "/", "HTTP/1.1", headers, bytes.NewReader([]byte{}))

		if err != nil {
			t.Fatal(err)
		}

		session = sessions.Session(req)

		if session.Get("user_id") != userId {
			t.Fatalf("Expected user id to be (%s) but got (%s)", userId, session.Get("user_id"))
		}

		if session.Get("role") != userRole {
			t.Fatalf("Expected role id to be (%s) but got (%s)", userRole, session.Get("role"))
		}

		session.Remove("user_id").Save()

		// Third Request
		cookie, err = url.ParseQuery(strings.ReplaceAll(req.Response.GetHeader("Set-Cookie"), "; ", "&"))

		if err != nil {
			t.Fatal(err)
		}

		headers = types.Headers{
			"cookie": strings.Join([]string{"session", cookie.Get("session")}, "="),
		}

		req, err = NewRequest("POST", "/", "HTTP/1.1", headers, bytes.NewReader([]byte{}))

		if err != nil {
			t.Fatal(err)
		}

		session = sessions.Session(req)

		if session.Get("user_id") != "" {
			t.Fatalf("Expected session user id to be empty but got (%s)", session.Get("user_id"))
		}

		if session.Get("role") != userRole {
			t.Fatalf("Expected role id to be (%s) but got (%s)", userRole, session.Get("role"))
		}

		session.Clear()

		// Fourth Request
		cookie, err = url.ParseQuery(strings.ReplaceAll(req.Response.GetHeader("Set-Cookie"), "; ", "&"))

		if err != nil {
			t.Fatal(err)
		}

		headers = types.Headers{
			"cookie": strings.Join([]string{"session", cookie.Get("session")}, "="),
		}

		if session.Get("role") != "" {
			t.Fatalf("Expected session user id to be empty but got (%s)", session.Get("role"))
		}
	})

	t.Run("TestSessionsManager", func(t *testing.T) {
		domain := "map.guarded.com"
		maxAge := ((60 * 60) * 24) * 5
		secure := true
		httpOnly := true
		sameSite := true

		sessions := InitSession("session", []byte(str.Random(10)))

		sessions.Domain(domain).MaxAge(maxAge).Secure(secure).HttpOnly(httpOnly).SameSite(sameSite)

		if sessions.store.Options.Domain != domain {
			t.Fatalf("Expected the domain to be (%s) but got (%s)", domain, sessions.store.Options.Domain)
		}

		if sessions.store.Options.MaxAge != maxAge {
			t.Fatalf("Expected the max age to be (%d) but got (%d)", maxAge, sessions.store.Options.MaxAge)
		}

		if sessions.store.Options.Secure != secure {
			t.Fatalf("Expected the secure to be (%t) but got (%t)", secure, sessions.store.Options.Secure)
		}

		if sessions.store.Options.HttpOnly != httpOnly {
			t.Fatalf("Expected the http only to be (%t) but got (%t)", httpOnly, sessions.store.Options.HttpOnly)
		}

		if sessions.store.Options.SameSite != 1 {
			t.Fatalf("Expected the same site to be (%d) but got (%d)", 1, sessions.store.Options.SameSite)
		}
	})

	t.Run("TestSessionManagerErrors", func(t *testing.T) {
		sessions := InitSession("session", []byte(str.Random(10)))

		req, err := NewRequest("POST", "/", "HTTP/1.1", make(types.Headers), bytes.NewReader([]byte{}))

		// First Request
		if err != nil {
			t.Fatalf("Something went wrong when trying to create request: %s", err.Error())
		}

		session := sessions.Session(req)

		emailError := "The email is required"
		passwordError := "The password is required"

		session.SetError("email", emailError).SetError("password", passwordError)

		session.Save() // SAVING SESSION

		// Second Request
		cookie, err := url.ParseQuery(strings.ReplaceAll(req.Response.GetHeader("Set-Cookie"), "; ", "&"))

		if err != nil {
			t.Fatal(err)
		}

		headers := types.Headers{
			"cookie": strings.Join([]string{"session", cookie.Get("session")}, "="),
		}

		req, err = NewRequest("POST", "/", "HTTP/1.1", headers, bytes.NewReader([]byte{}))

		if err != nil {
			t.Fatal(err)
		}

		session = sessions.Session(req)

		if session.Error("email") != emailError {
			t.Fatalf("Expected email error to be (%s) but got (%s)", emailError, session.Error("email"))
		}

		if session.Error("password") != passwordError {
			t.Fatalf("Expected password error to be (%s) but got (%s)", passwordError, session.Error("password"))
		}

		errors := session.Errors()

		if err := errors["email"]; err != emailError {
			t.Fatalf("Expected email error to be (%s) but got (%s)", err, err)
		}

		if err := errors["password"]; err != passwordError {
			t.Fatalf("Expected password error to be (%s) but got (%s)", err, err)
		}

		session.Save() // SAVING SESSION

		// Third Request
		cookie, err = url.ParseQuery(strings.ReplaceAll(req.Response.GetHeader("Set-Cookie"), "; ", "&"))

		if err != nil {
			t.Fatal(err)
		}

		headers = types.Headers{
			"cookie": strings.Join([]string{"session", cookie.Get("session")}, "="),
		}

		req, err = NewRequest("POST", "/", "HTTP/1.1", headers, bytes.NewReader([]byte{}))

		if err != nil {
			t.Fatal(err)
		}

		session = sessions.Session(req)

		if session.Error("email") != "" {
			t.Fatalf("Expected email error to be empty but got (%s)", session.Error("email"))
		}

		if session.Error("password") != "" {
			t.Fatalf("Expected email error to be empty but got (%s)", session.Error("password"))
		}

	})

	t.Run("TestSessionManagerCsrf", func(t *testing.T) {
		sessions := InitSession("session", []byte(str.Random(10)))

		req, err := NewRequest("POST", "/", "HTTP/1.1", make(types.Headers), bytes.NewReader([]byte{}))

		// First Request
		if err != nil {
			t.Fatalf("Something went wrong when trying to create request: %s", err.Error())
		}

		session := sessions.Session(req)

		session.Save() // SAVING SESSION

		// Second Request
		cookie, err := url.ParseQuery(strings.ReplaceAll(req.Response.GetHeader("Set-Cookie"), "; ", "&"))

		if err != nil {
			t.Fatal(err)
		}

		headers := types.Headers{
			"cookie": strings.Join([]string{"session", cookie.Get("session")}, "="),
		}

		req, err = NewRequest("POST", "/", "HTTP/1.1", headers, bytes.NewReader([]byte{}))

		if err != nil {
			t.Fatal(err)
		}

		req.Session = sessions.Session(req)

		if session.CsrfToken() == "" {
			t.Fatalf("Expected csrf to not be empty.")
		}

		if session.CsrfName() != CSRF_INPUT_NAME {
			t.Fatalf("Expected csrf name to be (%s) but got (%s)", CSRF_INPUT_NAME, session.CsrfName())
		}
	})

	t.Run("TestSessionManagerOld", func(t *testing.T) {
		sessions := InitSession("session", []byte(str.Random(10)))

		req, err := NewRequest("POST", "/", "HTTP/1.1", make(types.Headers), bytes.NewReader([]byte{}))

		// First Request
		if err != nil {
			t.Fatalf("Something went wrong when trying to create request: %s", err.Error())
		}

		session := sessions.Session(req).SetError("email", "The email is invalid")

		email := "jeo@doe.com"

		req.Form = url.Values{"email": []string{email}}

		session.Save() // SAVING SESSION

		// Second Request
		cookie, err := url.ParseQuery(strings.ReplaceAll(req.Response.GetHeader("Set-Cookie"), "; ", "&"))

		if err != nil {
			t.Fatal(err)
		}

		headers := types.Headers{
			"cookie": strings.Join([]string{"session", cookie.Get("session")}, "="),
		}

		req, err = NewRequest("POST", "/", "HTTP/1.1", headers, bytes.NewReader([]byte{}))

		if err != nil {
			t.Fatal(err)
		}

		req.Session = sessions.Session(req)

		if old := req.Session.Old("email"); old != email {
			t.Fatalf("Expected old email value to be (%s) but got (%s)", email, old)
		}

		if old := SessionOld(req)("email"); old != email {
			t.Fatalf("Expected old email value to be (%s) but got (%s)", email, old)
		}

		if old := SessionOld(req)("first_name"); old != "" {
			t.Fatalf("Expected old first name value to be empty but got (%s)", old)
		}
	})

	t.Run("TestSessionManagerHelperFunctions", func(t *testing.T) {
		userID := rand.Int()
		sessions := InitSession("session", []byte(str.Random(10)))

		req, err := NewRequest("POST", "/", "HTTP/1.1", make(types.Headers), bytes.NewReader([]byte{}))

		// First Request
		if err != nil {
			t.Fatalf("Something went wrong when trying to create request: %s", err.Error())
		}

		session := sessions.Session(req)

		firstNameError := "The first name is required"

		session.SetError("first_name", firstNameError).Set("user_id", userID)

		session.Save() // SAVING SESSION

		// Second Request
		cookie, err := url.ParseQuery(strings.ReplaceAll(req.Response.GetHeader("Set-Cookie"), "; ", "&"))

		if err != nil {
			t.Fatal(err)
		}

		headers := types.Headers{
			"cookie": strings.Join([]string{"session", cookie.Get("session")}, "="),
		}

		req, err = NewRequest("POST", "/", "HTTP/1.1", headers, bytes.NewReader([]byte{}))

		if err != nil {
			t.Fatal(err)
		}

		req.Session = sessions.Session(req)

		// Function - Has
		if !SessionHas(req)("first_name") {
			t.Fatalf("Expected session to have first_name error.")
		}

		if SessionHas(req)("email") {
			t.Fatalf("Expected session to not have email error.")
		}

		// Function - Error
		if err := SessionError(req)("first_name"); err != firstNameError {
			t.Fatalf("Expected first name error to be (%s) but got (%s).", firstNameError, err)
		}

		if err := SessionError(req)("email"); err != "" {
			t.Fatalf("Expected email error to be empty but got (%s).", err)
		}

		// Function - Error
		errors := SessionErrors(req)()

		if err := errors["first_name"]; err == "" {
			t.Fatalf("Expected first name error to be (%s) but got (%s),", err, firstNameError)
		}

		if err := errors["password"]; err != "" {
			t.Fatalf("Expected email error to be empty but got (%s).", err)
		}

		// Function - Session
		if id := SessionValue(req)("user_id"); id != cast.ToString(userID) {
			t.Fatalf("Expected session user id to be (%s) but got (%s).", cast.ToString(userID), id)
		}

		// Function - Csrf
		if token := SessionCsrfToken(req)(); token != req.Session.CsrfToken() {
			t.Fatalf("Expected csrf token to be (%s) but got (%s)", req.Session.CsrfToken(), token)
		}

		if name := SessionCsrfName(req)(); name == "" {
			t.Fatalf("Expected csrf name to be (%s) but got (%s)", req.Session.CsrfName(), name)
		}
	})
}
