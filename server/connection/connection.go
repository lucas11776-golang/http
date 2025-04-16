package connection

import (
	"net"
	"net/http"
	"strings"
)

type RequestCallback func(req *http.Request)

type Connection struct {
	conn *net.Conn
}

type Conn interface {
	Conn() interface{}
	IP() string
	Write(data []byte) error
	Read(data []byte) error
	Close() error
}

// Comment
func Init(conn *net.Conn) *Connection {
	return &Connection{conn: conn}
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
func (ctx *Connection) Close() error {
	return (*ctx.conn).Close()
}
