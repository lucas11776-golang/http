package ws

import (
	"fmt"
	"http/ws/frame"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"
)

// [opcode, len, mask, data]

// Comment
func replyServerWsTest(t *testing.T, concat []byte) (net.Listener, error) {
	listener, err := net.Listen("tcp", ":0")

	if err != nil {
		return nil, err
	}

	go func() {
		for {
			conn, err := listener.Accept()

			if err != nil {
				break
			}

			ws := Create(conn)

			ws.OnReady(func(ws *Ws) {
				ws.OnMessage(func(data []byte) {
					err := ws.WriteText(append(data, concat...))

					if err != nil {
						listener.Close()
						t.Errorf("Something went wrong when trying ot send data: %s", err.Error())
					}
				})
			})

			ws.Emit(EVENT_READY, []byte{})

			ws.Listen()
		}
	}()

	return listener, nil
}

// Comment
func closeServer(t *testing.T, listener net.Listener) {
	err := listener.Close()

	if err != nil {
		t.Errorf("Something went wrong when closing server: %s", err.Error())
	}
}

func closeConnection(t *testing.T, conn net.Conn) {
	err := conn.Close()

	if err != nil {
		t.Errorf("Something went wrong when closing connection: %s", err.Error())
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

		listener, err := replyServerWsTest(t, concat)

		if err != nil {
			t.Errorf("Something went wrong starting server: %s", err.Error())
		}

		conn, err := net.Dial("tcp", listener.Addr().String())

		if err != nil {
			closeServer(t, listener)
			t.Errorf("Something went wrong when connecting to server: %s", err.Error())
		}

		go func() {
			time.Sleep(time.Millisecond * 50)

			_, err = conn.Write(payload)

			if err != nil {
				closeServer(t, listener)
				closeConnection(t, conn)
				t.Errorf("Something went wrong when connecting to server: %s", err.Error())
			}
		}()

		buffer := make([]byte, MAX_PAYLOAD_SIZE)

		_, err = conn.Read(buffer)

		if err != nil {
			closeServer(t, listener)
			closeConnection(t, conn)
			t.Errorf("Something went wrong when reading connection: %s", err.Error())
		}

		response := string(buffer[2:])
		expectedResponse := string(append(data, concat...))

		if response == expectedResponse {
			t.Errorf("Expected the response to be (%s) but go (%s)", expectedResponse, response)
		}
	})

	t.Run("TestPing", func(t *testing.T) {
		text := frame.OPCODE_TEXT + 128

		fmt.Println("TEXT: ", text, frame.OPCODE_TEXT, frame.OPCODE_TEXT&text)
	})
}
