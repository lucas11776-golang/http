package connection

import (
	"net"
	"net/http"
	"strings"
)

type RequestCallback func(req *http.Request)

type Connection struct {
	Alive   bool
	conn    *net.Conn
	message []RequestCallback
	max     int
}

// Comment
func Init(conn *net.Conn, max int64) *Connection {
	return &Connection{
		Alive: true,
		conn:  conn,
		max:   int(max),
	}
}

// Comment
func (ctx *Connection) Conn() net.Conn {
	return *ctx.conn
}

// Comment
func (ctx *Connection) IP() string {
	return strings.Split((*ctx.conn).RemoteAddr().String(), ":")[0]
}

// Comment
func (ctx *Connection) Write(data []byte) error {
	_, err := (*ctx.conn).Write(data)

	return err
}

// Comment
func (ctx *Connection) Read(b []byte) ([]byte, error) {
	n, err := (*ctx.conn).Read(b)

	if err != nil {
		return nil, err
	}

	return b[:n], nil
}

// Comment
func (ctx *Connection) Message(callback RequestCallback) *Connection {
	ctx.message = append(ctx.message, callback)

	return ctx
}

// Comment
func (ctx *Connection) Close() error {
	return (*ctx.conn).Close()
}
