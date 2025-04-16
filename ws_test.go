package http

import (
	"bytes"
	"math/rand"
	"net"
	"strconv"
	"testing"

	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/types"
)

// TODO must disable test because first byte being read by request maybe
func TestWs(t *testing.T) {
	replyServerWsTest := func(concat []byte) (net.Listener, error) {
		// TODO Refactor ws test...
		listener, err := net.Listen("tcp", ":0")
		server := Server("127.0.0.1", 0)
		req, _ := NewRequest(METHOD_GET, "/", "HTTP/1.1", make(types.Headers), bytes.NewReader([]byte{}))

		req.Server = server

		if err != nil {
			return nil, err
		}

		go func() {
			for {
				conn, err := listener.Accept()

				if err != nil {
					break
				}

				req, err := NewRequest("GET", "/", "HTTP/1.1", types.Headers{}, bytes.NewReader([]byte{}))

				if err != nil {
					panic(err)
				}

				req.Server = server

				ws := InitWs(connection.Init(&conn), req)

				ws.Request = req

				ws.OnReady(func(ws *Ws) {
					ws.OnMessage(func(data []byte) {
						err := ws.Write(append(data, concat...))

						if err != nil {
							listener.Close()
						}
					})
				})

				ws.isReady()

				ws.Listen()

				server.Close()
			}
		}()

		return listener, nil
	}

	// Comment
	closeServer := func(t *testing.T, listener net.Listener) {
		err := listener.Close()

		if err != nil {
			t.Fatalf("Something went wrong when closing server: %s", err.Error())
		}
	}

	closeConnection := func(t *testing.T, conn net.Conn) {
		err := conn.Close()

		if err != nil {
			t.Fatalf("Something went wrong when closing connection: %s", err.Error())
		}
	}

	t.Run("TestSendMessage", func(t *testing.T) {
		data := []byte("Hello Number: ")
		concat := []byte(strconv.Itoa(int(rand.Float32() * 10000)))
		mask := []byte{}

		for i := 0; i < 4; i++ {
			mask = append(mask, byte(rand.Float32()*255))
		}

		maskData := []byte{}

		for i, b := range data {
			maskData = append(maskData, b^mask[i%len(mask)])
		}

		payload := []byte{129, byte(len(data))}
		payload = append(payload, mask...)
		payload = append(payload, maskData...)

		listener, err := replyServerWsTest(concat)

		if err != nil {
			t.Fatalf("Something went wrong starting server: %s", err.Error())
		}

		conn, err := net.Dial("tcp", listener.Addr().String())

		if err != nil {
			closeServer(t, listener)
			t.Fatalf("Something went wrong when connecting to server: %s", err.Error())
		}

		_, err = conn.Write(payload)

		if err != nil {
			closeServer(t, listener)
			closeConnection(t, conn)
			t.Fatalf("Something went wrong when connecting to server: %s", err.Error())
		}

		buffer := make([]byte, MAX_WEBSOCKET_PAYLOAD)

		_, err = conn.Read(buffer)

		if err != nil {
			closeServer(t, listener)
			closeConnection(t, conn)
			t.Fatalf("Something went wrong when reading connection: %s", err.Error())
		}

		response := string(buffer[2:])
		expectedResponse := string(append(data, concat...))

		if response == expectedResponse {
			t.Fatalf("Expected the response to be (%s) but go (%s)", expectedResponse, response)
		}
	})
}
