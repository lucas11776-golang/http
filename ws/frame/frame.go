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
	InvalidPayloadError = errors.New("Invalid payload")
)

type Frame struct {
	fin        byte
	opcode     Opcode
	mask       byte
	size       uint64
	maskingKey byte
	data       []byte
	payload    []byte
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
		return nil, InvalidPayloadError
	}

	head := payload[:2]
	size := uint16(head[1] & 0x7F)
	frame := &Frame{payload: payload}

	if size < 126 {
		if len(payload) < int(size)+6 {
			return nil, InvalidPayloadError
		}

		frame.size = uint64(size)
		frame.data = unmask(payload[2:6], payload[6:frame.size+6])

		return frame, nil
	}

	if size == 126 {
		if len(payload) < 8 {
			return nil, InvalidPayloadError
		}

		frame.size = uint64(binary.BigEndian.Uint16(payload[2:4]))

		if len(payload) < int(frame.size)+8 {
			return nil, InvalidPayloadError
		}

		frame.data = unmask(payload[4:8], payload[8:frame.size+8])

		return frame, nil
	}

	if len(payload) < 10 {
		return nil, InvalidPayloadError
	}

	frame.size = uint64(binary.BigEndian.Uint64(payload[2:10]))

	if len(payload) < int(frame.size+14) {
		return nil, InvalidPayloadError
	}

	frame.data = unmask(payload[10:14], payload[14:frame.size+14])

	return frame, nil
}

// Comment
func Encode(opcode Opcode, data []byte) *Frame {
	size := len(data)
	frame := &Frame{data: data}
	opc := OPCODE_START + opcode

	frame.size = uint64(size)

	if size < 126 {
		payload := make([]byte, 2)
		payload[0] = byte(opc)
		payload[1] = byte(size)

		payload = append(payload, data...)

		frame.payload = payload

		return frame
	}

	if size >= 126 && size < int(math.Pow(2, 16)) {
		payload := make([]byte, 2)
		payload[0] = byte(opc)
		payload[1] = 126

		length := make([]byte, 2)

		binary.BigEndian.PutUint16(length, uint16(size))

		payload = append(payload, length...)
		payload = append(payload, data...)

		frame.payload = payload

		return frame
	}

	payload := make([]byte, 2)
	payload[0] = byte(opcode + OPCODE_START)
	payload[1] = 127

	length := make([]byte, 8)

	binary.BigEndian.PutUint64(length, uint64(size))

	payload = append(payload, length...)
	payload = append(payload, data...)

	frame.payload = payload

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
