package http

import (
	"bytes"
	"fmt"
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

		session.Set("user_id", userId)

		cookie, err := url.ParseQuery(strings.ReplaceAll(req.Response.Header["Set-Cookie"][0], "; ", "&"))

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

		fmt.Println("TEST:")
	})

	t.Run("TestOptionsSessionsManager", func(t *testing.T) {
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

	t.Run("TestOptionSessionManager", func(t *testing.T) {
		path := "dashboard"
		sessions := InitSession("session", []byte(str.Random(10)))

		req, err := NewRequest("POST", "/", "HTTP/1.1", make(types.Headers), bytes.NewReader([]byte{}))

		if err != nil {
			t.Fatalf("Something went wrong when trying to create request: %s", err.Error())
		}

		session := sessions.Session(req).Path(path)

		session.Set("user_id", "1")

		cookie, err := url.ParseQuery(strings.ReplaceAll(req.Response.Header["Set-Cookie"][0], "; ", "&"))

		if err != nil {
			t.Fatalf("Something went wrong when trying to parse cookie to query: %s", err)
		}

		if cookie.Get("Path") != path {
			t.Fatalf("Expected the path to be (%s) but got (%s)", path, cookie["path"][0])
		}
	})
}
