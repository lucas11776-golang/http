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

		headers = types.Headers{
			"cookie": strings.Join([]string{"session", cookie.Get("session")}, "="),
		}

		req, err = NewRequest("POST", "/", "HTTP/1.1", headers, bytes.NewReader([]byte{}))

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

		headers = types.Headers{
			"cookie": strings.Join([]string{"session", cookie.Get("session")}, "="),
		}

		if session.Get("role") != "" {
			t.Fatalf("Expected session user id to be empty but got (%s)", session.Get("role"))
		}
	})
}
