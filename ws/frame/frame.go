package frame

import (
	"encoding/binary"
	"errors"
	"math"
)

type Opcode byte

const OPCODE_START Opcode = 128

// TODO check if error opcode exists
const (
	OPCODE_CONTINUATION     Opcode = 0x00
	OPCODE_TEXT             Opcode = 0x01
	OPCODE_BINARY           Opcode = 0x02
	OPCODE_CONNECTION_CLOSE Opcode = 0x08
	OPCODE_PING             Opcode = 0x09
	OPCODE_PONG             Opcode = 0x0A
)

var (
	ErrInvalidPayload = errors.New("invalid payload")
)

type Frame struct {
	opcode  Opcode
	size    uint64
	data    []byte
	payload []byte
}

// Comment
func unmask(mask []byte, data []byte) []byte {
	for i, masked := range data {
		data[i] = masked ^ mask[i%len(mask)]
	}

	return data
}

// Comment
func Decode(payload []byte) (*Frame, error) {
	if len(payload) < 2 {
		return nil, ErrInvalidPayload
	}

	head := payload[:2]
	size := uint16(head[1] & 0x7F)
	frame := &Frame{payload: payload}

	if size < 126 {
		if len(payload) < int(size)+6 {
			return nil, ErrInvalidPayload
		}

		frame.size = uint64(size)
		frame.data = unmask(payload[2:6], payload[6:frame.size+6])

		return frame, nil
	}

	if size == 126 {
		if len(payload) < 8 {
			return nil, ErrInvalidPayload
		}

		frame.size = uint64(binary.BigEndian.Uint16(payload[2:4]))

		if len(payload) < int(frame.size)+8 {
			return nil, ErrInvalidPayload
		}

		frame.data = unmask(payload[4:8], payload[8:frame.size+8])

		return frame, nil
	}

	if len(payload) < 10 {
		return nil, ErrInvalidPayload
	}

	frame.size = uint64(binary.BigEndian.Uint64(payload[2:10]))

	if len(payload) < int(frame.size+14) {
		return nil, ErrInvalidPayload
	}

	frame.data = unmask(payload[10:14], payload[14:frame.size+14])

	return frame, nil
}

// Comment
func Encode(opcode Opcode, data []byte) *Frame {
	frame := &Frame{
		opcode: OPCODE_START + opcode,
		data:   data,
		size:   uint64(len(data)),
	}

	frame.payload = []byte{byte(frame.opcode)}

	if frame.size < 126 {
		frame.payload = append(frame.payload, byte(frame.size))
	}

	if frame.size >= 126 && frame.size < uint64(math.Pow(2, 16)) {
		length := make([]byte, 2)

		binary.BigEndian.PutUint16(length, uint16(frame.size))

		frame.payload = append(frame.payload, 126)
		frame.payload = append(frame.payload, length...)
	}

	if frame.size >= uint64(math.Pow(2, 16)) {
		length := make([]byte, 8)

		binary.BigEndian.PutUint64(length, uint64(frame.size))

		frame.payload = append(frame.payload, 127)
		frame.payload = append(frame.payload, length...)
	}

	frame.payload = append(frame.payload, data...)

	return frame
}

// Comment
func (ctx *Frame) Opcode() Opcode {
	return Opcode(ctx.payload[0] - byte(OPCODE_START))
}

// Comment
func (ctx *Frame) IsContinuation() bool {
	return ctx.Opcode() == OPCODE_CONTINUATION
}

// Comment
func (ctx *Frame) IsBinary() bool {
	return ctx.Opcode() == OPCODE_BINARY
}

// Comment
func (ctx *Frame) IsText() bool {
	return ctx.Opcode() == OPCODE_TEXT
}

// Comment
func (ctx *Frame) IsClose() bool {
	return ctx.Opcode() == OPCODE_CONNECTION_CLOSE
}

// Comment
func (ctx *Frame) IsPing() bool {
	return ctx.Opcode() == OPCODE_PING
}

// Comment
func (ctx *Frame) IsPong() bool {
	return ctx.Opcode() == OPCODE_PONG
}

// Comment
func (ctx *Frame) Length() uint64 {
	return ctx.size
}

// Comment
func (ctx *Frame) Data() []byte {
	return ctx.data
}

// Comment
func (ctx *Frame) Payload() []byte {
	return ctx.payload
}
