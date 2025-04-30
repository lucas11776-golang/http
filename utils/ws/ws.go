package ws

import (
	"errors"
	"fmt"
	"net"
	u "net/url"
	"strings"

	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/types"
	str "github.com/lucas11776-golang/http/utils/strings"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	MaxWebSocketPayloadSize int = 1024 * 4
)

type Ws struct {
	Conn        *connection.Connection
	PayloadSize int
}

// Comment
func (ctx *Ws) Read() ([]byte, error) {
	payload := make([]byte, ctx.PayloadSize)

	n, err := ctx.Conn.Conn().Read(payload)

	if err != nil {
		return nil, err
	}

	return ctx.decode(payload[:n])

	// if payload[1] < 126 {
	// 	return payload[2:n], nil
	// } else if payload[1] == 126 {
	// 	return payload[4:n], nil
	// } else if payload[1] == 127 {
	// 	return payload[10:n], nil
	// } else {
	// 	return nil, errors.New("invalid payload size")
	// }
}

// Comment
func (ctx *Ws) decode(payload []byte) ([]byte, error) {
	if payload[1] < 126 {
		return payload[2:], nil
	} else if payload[1] == 126 {
		return payload[4:], nil
	} else if payload[1] == 127 {
		return payload[10:], nil
	} else {
		return nil, errors.New("invalid payload size")
	}
}

// Comment
func (ctx *Ws) Close() error {
	return ctx.Conn.Close()
}

// Comment
func Connect(url string, headers types.Headers) (*Ws, error) {
	path, err := u.Parse(url)

	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", path.Host)

	if err != nil {
		return nil, err
	}

	request := []string{
		fmt.Sprintf("GET %s HTTP/1.1", path.Path),
		"Connection: Upgrade",
		fmt.Sprintf("Sec-Websocket-Key: %s", str.Random(20)),
		"Sec-Websocket-Version: 13",
		"Pragma: no-cache",
		"Upgrade: websocket",
		fmt.Sprintf("Host: %s", conn.LocalAddr().String()),
	}

	for k, v := range headers {
		request = append(request, fmt.Sprintf("%s: %s", cases.Title(language.English).String(k), v))
	}

	_, err = conn.Write([]byte(strings.Join(append(request, "\r\n"), "\r\n")))

	if err != nil {
		return nil, err
	}

	ws := &Ws{
		Conn:        connection.Init(&conn),
		PayloadSize: MaxWebSocketPayloadSize,
	}

	handshake := make([]byte, 1024)

	n, err := conn.Read(handshake)

	if err != nil {
		return nil, err
	}

	// TODO: need to validate Sec-Websocket-Accept, Upgrade and Status code - (this is lazy work)
	if strings.Split(string(handshake[:n]), "\r\n")[0] != "HTTP/1.1 101 Switching Protocols" {
		return nil, errors.New("invalid handshake response")
	}

	return ws, nil
}
