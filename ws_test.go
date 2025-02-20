package http

import (
	"bytes"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/types"
)

// [opcode, len, mask, data]

// Comment
func replyServerWsTest(concat []byte) (net.Listener, error) {
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

			ws := InitWs(connection.Init(&conn, MAX_WEBSOCKET_PAYLOAD))

			ws.Request = req

			ws.OnReady(func(ws *Ws) {
				ws.OnMessage(func(data []byte) {
					err := ws.Write(append(data, concat...))

					if err != nil {
						listener.Close()
					}
				})
			})

			ws.Emit(EVENT_READY, []byte{})

			ws.Listen()

			server.Close()
		}
	}()

	return listener, nil
}

// Comment
func closeServer(t *testing.T, listener net.Listener) {
	err := listener.Close()

	if err != nil {
		t.Fatalf("Something went wrong when closing server: %s", err.Error())
	}
}

func closeConnection(t *testing.T, conn net.Conn) {
	err := conn.Close()

	if err != nil {
		t.Fatalf("Something went wrong when closing connection: %s", err.Error())
	}
}

func TestWs(t *testing.T) {
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

		time.Sleep(time.Millisecond * 50)

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
