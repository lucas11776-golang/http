package connection

import (
	"bufio"
	"net"
	"net/http"
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
func (ctx *Connection) Write(data []byte) error {
	_, err := ctx.Conn().Write(data)

	return err
}

// Comment
func (ctx *Connection) Message(callback RequestCallback) *Connection {
	ctx.message = append(ctx.message, callback)

	return ctx
}

// Comment
func (ctx *Connection) Listen() {
	for {
		req, err := http.ReadRequest(bufio.NewReader(bufio.NewReaderSize(ctx.Conn(), ctx.max)))

		if err != nil {
			ctx.Alive = false
			break
		}

		for _, callback := range ctx.message {
			go func() {
				callback(req)
			}()
		}
	}

	defer ctx.Conn().Close()
}

// Comment
func (ctx *Connection) Close() error {
	return ctx.Conn().Close()
}
