package connection

import (
	"net"
)

type MessageCallback func(data []byte)

type Connection struct {
	Alive   bool
	conn    *net.Conn
	message []MessageCallback
	max     int64
}

// Comment
func Init(conn *net.Conn, max int64) *Connection {
	return &Connection{
		Alive: true,
		conn:  conn,
		max:   max,
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
func (ctx *Connection) Message(callback MessageCallback) *Connection {
	ctx.message = append(ctx.message, callback)

	return ctx
}

// Comment
func (ctx *Connection) Listen() {
	buffer := make([]byte, ctx.max)

	for {
		size, err := ctx.Conn().Read(buffer)

		if err != nil {
			ctx.Alive = false
			break
		}

		for _, callback := range ctx.message {
			go func() {
				callback(buffer[:size])
			}()
		}
	}
}

// Comment
func (ctx *Connection) Close() error {
	return ctx.Conn().Close()
}
