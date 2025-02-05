package session

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/lucas11776-golang/http"
	"github.com/lucas11776-golang/http/types"
	"github.com/lucas11776-golang/http/utils/strings"
)

func TestSession(t *testing.T) {
	sessions := Init("session", []byte(strings.Random(10)))

	t.Run("TestGetSession", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", "HTTP/1.1", make(types.Headers), bytes.NewReader([]byte{}))

		if err != nil {
			t.Fatalf("Something went wrong when trying to create request: %s", err.Error())
		}

		session, err := sessions.Session(req)

		if err != nil {
			t.Fatalf("Something went wrong when trying to get session: %s", err.Error())
		}

		session.Set("user_id", "1")

		fmt.Println("TEST:", session)
	})
}
