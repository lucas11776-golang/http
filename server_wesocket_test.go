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

	AuthMiddleware := func(req *Request, res *Response, next Next) *Response {
		if req.GetHeader("Authorization") != authorization {
			return res
		}
		return next()
	}

	t.Run("TestHandshakeReplay", func(t *testing.T) {
		server := Server("127.0.0.1", 0).SetMaxWebsocketPayload(1024 * 10)

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
		})

		go server.Listen()

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
			"Host: 127.0.0.1:4567",
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

		res, err := HttpToResponse(string(buf[:n]))

		if err != nil {
			t.Fatalf("Invalid handshake response: %v", err)
		}

		alg := sha1.New()

		alg.Write([]byte(strings.Join([]string{"TnjNK5ivR7MUvlou4Ilj9g==", SEC_WEB_SOCKET_ACCEPT_STATIC}, "")))

		if res.GetHeader("Upgrade") != "websocket" {
			t.Fatalf("Expected Upgrage header to be (%s) but got (%s)", "websocket", res.GetHeader("Upgrade"))
		}

		if res.GetHeader("Connection") != "Upgrade" {
			t.Fatalf("Expected Connection header to be (%s) but got (%s)", "Upgrade", res.GetHeader("Connection"))
		}

		secWebsocketAccept := base64.StdEncoding.EncodeToString(alg.Sum(nil))

		if res.GetHeader("Sec-Websocket-Accept") != secWebsocketAccept {
			t.Fatalf(
				"Expected Sec-Websocket-Accept header to be (%s) but got (%s)",
				secWebsocketAccept,
				res.GetHeader("Sec-Websocket-Accept"),
			)
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

		// TODO: 1. Quatam State - (fails or pass) - Github workflow
		if expectedResponse != response {
			t.Fatalf("Expected ws payload to be (%s) but got (%s)", expectedResponse, response)
		}

		conn.Close()
		server.Close()
	})

	t.Run("TestWebsocketMiddlewareUnauthorized", func(t *testing.T) {
		server := Server("127.0.0.1", 0).SetMaxWebsocketPayload(1024 * 10)

		server.Route().Ws("auth", func(req *Request, ws *Ws) {
			ws.WriteJson(map[string]string{"message": authorizedMessage})
		}).Middleware(AuthMiddleware)

		go server.Listen()

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
		server := Server("127.0.0.1", 0).SetMaxWebsocketPayload(1024 * 10)

		server.Route().Ws("auth", func(req *Request, ws *Ws) {
			time.Sleep(time.Microsecond * 10)
			ws.WriteJson(map[string]string{"message": authorizedMessage})
		}).Middleware(AuthMiddleware)

		go server.Listen()

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
