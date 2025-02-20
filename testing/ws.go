package testing

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net"
	"strings"

	"github.com/lucas11776-golang/http"
	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/ws/frame"
)

type Ws struct {
	testcase *TestCase
	testing  *Testing
	sessions Values
}

// comment
func NewWs(testcase *TestCase) *Ws {
	return &Ws{
		testcase: testcase,
		testing:  testcase.testing,
		sessions: make(Values),
	}
}

// Comment
func (ctx *Ws) handshake(uri string, conn *connection.Connection) error {
	req := strings.Join([]string{
		strings.Join([]string{"GET", fmt.Sprintf("/%s", strings.Trim(uri, "/")), "HTTP/1.1"}, " "),
		"Connection: Upgrade",
		"Sec-Websocket-Key: TnjNK5ivR7MUvlou4Ilj9g==",
		"Sec-Websocket-Version: 13",
		"Pragma: no-cache",
		"Upgrade: websocket",
		"\r\n",
	}, "\r\n")

	err := conn.Write([]byte(req))

	if err != nil {
		return err
	}

	http := make([]byte, 1024)

	_, err = conn.Conn().Read(http)

	return err
}

// Comment
func (ctx *Ws) Connect(uri string) *WsResponse {
	conn, err := net.Dial("tcp", ctx.testcase.http.Host())

	if err != nil {
		ctx.testing.Fatalf("Something went wrong when trying to establish connection: %v", err)
	}

	connection := connection.Init(&conn, http.MAX_WEBSOCKET_PAYLOAD)

	err = ctx.handshake(uri, connection)

	if err != nil {
		ctx.testing.Fatalf("Something went wrong when trying to establish handshake: %v", err)
	}

	return &WsResponse{
		testcase: ctx.testcase,
		testing:  ctx.testing,
		conn:     connection,
	}
}

type WsResponse struct {
	testcase *TestCase
	testing  *Testing
	conn     *connection.Connection
}

// Comment
func (ctx *WsResponse) mask(data []byte) (mask []byte, masked []byte) {
	msk := make([]byte, 4)

	for i := range msk {
		msk[i] = byte(rand.Float32() * 255)
	}

	for i, d := range data {
		data[i] = d ^ msk[i%len(msk)]
	}

	return msk, data
}

// Comment
func (ctx *WsResponse) write(opcode frame.Opcode, data []byte) error {
	opc := frame.OPCODE_START + opcode
	mask, data := ctx.mask(data)
	size := len(data)

	payload := []byte{byte(opc)}

	if size < 126 {
		payload = append(payload, byte(size))
		payload = append(payload, mask...)
	}

	if size >= 126 && size < int(math.Pow(2, 16)) {
		length := make([]byte, 2)

		binary.BigEndian.PutUint16(length, uint16(size))

		payload = append(payload, 126)
		payload = append(payload, length...)
		payload = append(payload, mask...)
	}

	if size > int(math.Pow(2, 16)) {
		length := make([]byte, 8)

		binary.BigEndian.PutUint64(length, uint64(size))

		payload = append(payload, 127)
		payload = append(payload, length...)
		payload = append(payload, mask...)
	}

	payload = append(payload, data...)

	return ctx.conn.Write(payload)
}

// Comment
func (ctx *WsResponse) WriteText(data []byte) error {
	return ctx.write(frame.OPCODE_TEXT, data)
}

// Comment
func (ctx *WsResponse) WriteBinary(data []byte) error {
	return ctx.write(frame.OPCODE_BINARY, data)
}

// Comment
func (ctx *WsResponse) WriteJson(v any) error {
	data, err := json.Marshal(v)

	if err != nil {
		return err
	}

	return ctx.write(frame.OPCODE_TEXT, data)
}

// Comment
func (ctx WsResponse) Read() []byte {
	data := make([]byte, int(math.Pow(2, 18)))

	n, err := ctx.conn.Conn().Read(data)

	if err != nil {
		ctx.testcase.testing.Fatalf("Something went wrong when trying to read payload: %v", err)
	}

	size := data[1]

	if size < 126 {
		return data[2 : size+2]
	}

	if size == 126 {
		return data[4 : binary.BigEndian.Uint16(data[2:4])+4]
	}

	if size == 127 {
		return data[10 : binary.BigEndian.Uint64(data[2:10])+10]
	}

	return data[:n]
}

// Comment
func (ctx *WsResponse) AssertRead(payload []byte) *WsResponse {
	data := ctx.Read()

	if string(data) != string(payload) {
		ctx.testing.Fatalf("Expected payload to be (%s) but got (%s)", string(payload), string(data))
	}

	return ctx
}
