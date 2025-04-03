package http

import (
	"encoding/json"

	"github.com/lucas11776-golang/http/server/connection"
	"github.com/lucas11776-golang/http/ws/frame"
)

const (
	MAX_WEBSOCKET_PAYLOAD = 1024 * 2
)

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
	Alive   bool
	Request *Request
	conn    *connection.Connection
	events  Events
	ready   []ReadyCallback
}

// Comment
func InitWs(conn *connection.Connection, req *Request) *Ws {
	return &Ws{
		Alive:   true,
		conn:    conn,
		events:  make(Events),
		Request: req,
	}
}

// Comment
func (ctx *Ws) OnReady(callback ReadyCallback) {
	ctx.ready = append(ctx.ready, callback)
}

// Comment
func (ctx *Ws) appendEvent(event Event, callback EventCallback) {
	ctx.events[event] = append(ctx.events[event], callback)
}

// Comment
func (ctx *Ws) OnMessage(callback EventCallback) {
	ctx.appendEvent(EVENT_MESSAGE, callback)
}

// Comment
func (ctx *Ws) OnError(callback EventCallback) {
	ctx.appendEvent(EVENT_ERROR, callback)
}

// Comment
func (ctx *Ws) OnPing(callback EventCallback) {
	ctx.appendEvent(EVENT_PING, callback)
}

// Comment
func (ctx *Ws) OnPong(callback EventCallback) {
	ctx.appendEvent(EVENT_PONG, callback)
}

// Comment
func (ctx *Ws) OnClose(callback EventCallback) {
	ctx.appendEvent(EVENT_CLOSE, callback)
}

func (ctx *Ws) isReady() {
	for _, callback := range ctx.ready {
		callback(ctx)
	}
}

// Comment
func (ctx *Ws) Emit(event Event, data []byte) {
	if callbacks, ok := ctx.events[event]; ok {
		for i := range callbacks {
			callbacks[i](data)
		}
	}
}

// Comment
func (ctx *Ws) Write(data []byte) error {
	return ctx.conn.Write(frame.Encode(frame.OPCODE_TEXT, data).Payload())
}

// Comment
func (ctx *Ws) WriteBinary(data []byte) error {
	return ctx.conn.Write(frame.Encode(frame.OPCODE_BINARY, data).Payload())
}

// Comment
func (ctx *Ws) WriteJson(v any) error {
	json, err := json.Marshal(v)

	if err != nil {
		return err
	}

	return ctx.Write(json)
}

// Comment
func (ctx *Ws) emitter(opcode frame.Opcode, data []byte) {
	switch opcode {

	case frame.OPCODE_CONTINUATION:

	case frame.OPCODE_BINARY, frame.OPCODE_TEXT:
		ctx.Emit(EVENT_MESSAGE, data)

	case frame.OPCODE_CONNECTION_CLOSE:
		ctx.Emit(EVENT_CLOSE, data)

	case frame.OPCODE_PING:
		ctx.Emit(EVENT_PING, data)

	case frame.OPCODE_PONG:
		ctx.Emit(EVENT_PONG, data)

	}
}

// Comment
func (ctx *Ws) Listen() {
	// TODO temp fix must find out why in first request OPCODE is not there 1 == spc -> http2 reading fix byte of payload
	var requests int64 = 0

	for {
		requests++

		payload := make([]byte, ctx.Request.Server.MaxWebSocketPayloadSize)

		n, err := ctx.conn.Conn().Read(payload)

		if err != nil {
			ctx.Emit(EVENT_ERROR, []byte(err.Error()))

			break
		}

		// TODO :() NO!!!!!
		if requests == 1 {
			payload = append([]byte{129}, payload[:n]...)
		}

		frm, err := frame.Decode(payload[:n])

		// fmt.Println("---> ----->", payload[:n], err)

		if err != nil {
			continue
		}

		ctx.emitter(frm.Opcode(), frm.Data())
	}
}
