package ws

import (
	"fmt"
	"http/ws/frame"
	"net"
)

const MAX_PAYLOAD_SIZE = 2048

type Event string

const (
	EVENT_READY   Event = "ready"
	EVENT_MESSAGE Event = "message"
	EVENT_PING    Event = "ping"
	EVENT_PONG    Event = "pong"
	EVENT_ERROR   Event = "error"
	EVENT_CLOSE   Event = "close"
)

type ReadyCallback func(ws *Ws)

type EventCallback func(data []byte)

type Events map[Event][]EventCallback

type Ws struct {
	Alive bool
	conn  net.Conn
	event Events
	ready []ReadyCallback
}

// Comment
func Create(conn net.Conn) *Ws {
	return &Ws{
		Alive: true,
		conn:  conn,
		event: make(Events),
	}
}

// Comment
func (ctx *Ws) OnReady(callback ReadyCallback) {
	ctx.ready = append(ctx.ready, callback)
}

// Comment
func (ctx *Ws) OnMessage(callback EventCallback) {
	ctx.event[EVENT_MESSAGE] = append(ctx.event[EVENT_MESSAGE], callback)
}

// Comment
func (ctx *Ws) OnPing(callback EventCallback) {

}

// Comment
func (ctx *Ws) OnError(callback EventCallback) {

}

// Comment
func (ctx *Ws) OnClose(callback EventCallback) {
}

// Comment
func (ctx *Ws) Emit(event Event, data []byte) {
	switch event {
	case EVENT_READY:
		for _, callback := range ctx.ready {
			go func() {
				callback(ctx)
			}()
		}
		break
	default:
		for _, callback := range ctx.event[event] {
			callback(data)
		}
	}
}

// Comment
func (ctx *Ws) Write(data []byte) error {
	_, err := ctx.conn.Write(data)

	return err
}

// Comment
func (ctx *Ws) WriteText(data []byte) error {
	return ctx.Write(frame.Encode(frame.OPCODE_TEXT, data).Payload())
}

// Comment
func (ctx *Ws) WriteBinary(data []byte) error {
	return nil
}

// Comment
func (ctx *Ws) Listen() {
	payload := make([]byte, MAX_PAYLOAD_SIZE)

	for {

		_, err := ctx.conn.Read(payload)

		if err != nil {
			ctx.Emit(EVENT_ERROR, []byte(err.Error()))

			break
		}

		frame, err := frame.Decode(payload)

		if err != nil {
			continue
		}

		if frame.IsBinary() || frame.IsText() {
			ctx.Emit(EVENT_MESSAGE, frame.Data())

			continue
		}

		fmt.Println(string(frame.Data()))
	}
}
