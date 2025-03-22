package http

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/lucas11776-golang/http/types"
	"github.com/lucas11776-golang/http/utils/ws"
)

func TestServerWebSocket(t *testing.T) {

	const (
		wsResponse          = "Hello World from :name !!!"
		authorization       = "test@123"
		unauthorizedMessage = "unauthorized access"
		authorizedMessage   = "Welcome to route"
	)

	serve := func() *HTTP {
		server := Server("127.0.0.1", 0).SetMaxWebsocketPayload(1024 * 10)

		auth := func(req *Request, res *Response, next Next) *Response {
			if req.GetHeader("Authorization") != authorization {
				return res
			}

			return next()
		}

		server.Route().Group("", func(route *Router) {
			route.Ws("/", func(req *Request, ws *Ws) {
				ws.OnReady(func(ws *Ws) {
					ws.OnMessage(func(data []byte) {
						err := ws.Write([]byte(strings.ReplaceAll(wsResponse, ":name", string(data))))

						if err != nil {
							t.Fatalf("Something went wrong when trying to send message: %s", err.Error())
						}

					})
				})
			})
			route.Ws("auth", func(req *Request, ws *Ws) {
				time.Sleep(time.Microsecond * 10)
				ws.WriteJson(map[string]string{"message": authorizedMessage})
			}).Middleware(auth)
		})

		go func() {
			server.Listen()
		}()

		return server
	}

	t.Run("TestHandshakeReplay", func(t *testing.T) {
		server := serve()
		conn, err := net.Dial("tcp", server.Host())

		if err != nil {
			t.Fatalf("Something went wrong when trying to connect to server: %s", err.Error())
		}

		htp := strings.Join([]string{
			"GET / HTTP/1.1",
			"Connection: Upgrade",
			"Sec-Websocket-Key: TnjNK5ivR7MUvlou4Ilj9g==",
			"Sec-Websocket-Version: 13",
			"Pragma: no-cache",
			"Upgrade: websocket",
			"\r\n",
		}, "\r\n")

		_, err = conn.Write([]byte(htp))

		if err != nil {
			t.Fatalf("Something went wrong when trying send request: %s", err.Error())
		}

		buf := make([]byte, server.MaxWebSocketPayloadSize)

		n, err := conn.Read(buf)

		if err != nil {
			t.Fatalf("Something went wrong when trying read connection: %s", err.Error())
		}

		alg := sha1.New()

		alg.Write([]byte(strings.Join([]string{"TnjNK5ivR7MUvlou4Ilj9g==", SEC_WEB_SOCKET_ACCEPT_STATIC}, "")))

		res := NewResponse("HTTP/1.1", HTTP_RESPONSE_SWITCHING_PROTOCOLS, types.Headers{
			"Connection":           "Upgrade",
			"Sec-Websocket-Accept": base64.StdEncoding.EncodeToString(alg.Sum(nil)),
			"Upgrade":              "websocket",
		}, []byte{})

		http := string(buf[:n])
		expectedHttp := ParseHttpResponse(res)

		if http != expectedHttp {
			t.Fatalf("Expected handshake response to be (%s) but got (%s)", expectedHttp, http)
		}

		name := strings.Join([]string{"user", strconv.Itoa(int(rand.Float32() * 1000))}, "")
		mask := []byte{34, 43, 56, 32}
		payload := []byte{129, byte(len(name))}
		payload = append(payload, mask...)

		for i, b := range []byte(name) {
			payload = append(payload, b^mask[i%4])
		}

		_, err = conn.Write(payload)

		if err != nil {
			t.Fatalf("Something went wrong when trying send payload: %s", err.Error())
		}

		buffNew := make([]byte, server.MaxWebSocketPayloadSize)

		n, err = conn.Read(buffNew)

		if err != nil {
			t.Fatalf("Something went wrong when trying read connection: %s", err.Error())
		}

		expectedResponse := strings.ReplaceAll(wsResponse, ":name", name)

		response := string(buffNew[2:n])

		if err != nil {
			t.Fatalf("Something went wrong when trying to decode payload: %s", err.Error())
		}

		if expectedResponse != response {
			t.Fatalf("Expected ws payload to be (%s) but got (%s)", expectedResponse, response)
		}

		conn.Close()
		server.Close()
	})

	t.Run("TestWebsocketMiddlewareUnauthorized", func(t *testing.T) {
		server := serve()
		w, err := ws.Connect(fmt.Sprintf("http://%s/%s", server.Host(), "auth"), types.Headers{})

		if err != nil {
			t.Fatalf("Something went wrong when connecting: %v", err)
		}

		_, err = w.Read()

		if err != io.EOF {
			t.Fatalf("Error must be type of %v", io.EOF)
		}

		w.Close()
		server.Close()
	})

	t.Run("TestWebsocketMiddlewareAuthorized", func(t *testing.T) {
		server := serve()
		w, err := ws.Connect(fmt.Sprintf("http://%s/%s", server.Host(), "auth"), types.Headers{
			"Authorization": authorization,
		})

		if err != nil {
			t.Fatalf("Something went wrong when connecting: %v", err)
		}

		data, err := w.Read()

		if err != nil {
			t.Fatalf("Something went wrong when reading response: %v", err)
		}

		message := make(map[string]string)

		err = json.Unmarshal(data, &message)

		if err != nil {
			t.Fatalf("Something went wrong when Unmarshal response data: %v", err)
		}

		if message["message"] != authorizedMessage {
			t.Fatalf("Expected response message to be (%s) but got (%s)", authorizedMessage, message["message"])
		}

		w.Close()
		server.Close()
	})
}
