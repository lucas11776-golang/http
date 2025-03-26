package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"testing"

	"github.com/lucas11776-golang/http/utils/rsa"
	"golang.org/x/net/http2"
)

//

func TestHtt2(t *testing.T) {
	type Message struct {
		Message string `json:"message"`
	}

	tResponse := &Message{Message: fmt.Sprintf("Hello: %f", rand.Float32()*100000)}

	serve := func() *HTTP {
		cert, key, err := rsa.GenerateCertificate("localhost")

		if err != nil {
			t.Fatalf("TLS error: %v", err)
		}

		server := ServerTLS("127.0.0.1", 0, cert, key)

		server.Route().Get("/", func(req *Request, res *Response) *Response {
			return res.SetStatus(HTTP_RESPONSE_OK).Json(tResponse)
		})

		go server.Listen()

		return server
	}

	t.Run("HTTP2NegotiationRequest", func(t *testing.T) {
		server := serve()

		client := &http.Client{Transport: &http2.Transport{
			AllowHTTP: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

		res, err := client.Get(fmt.Sprintf("http://%s", server.Host()))

		if err != nil {
			t.Fatalf("Error making request: %v", err)
		}

		body, err := io.ReadAll(res.Body)

		if err != nil {
			t.Fatalf("Error reading response: %v", err)
		}

		if res.Proto != "HTTP/2.0" {
			t.Fatalf("Expectect response protocal to be (%s) but got (%s)", "HTTP/2.0", res.Proto)
		}

		tBody, _ := json.Marshal(tResponse)

		if string(tBody) != string(body) {
			t.Fatalf("Expectect response body to be (%s) but got (%s)", string(tBody), body)
		}

		server.Close()
	})
}
