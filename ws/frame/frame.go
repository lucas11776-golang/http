package frame

import (
	"encoding/binary"
)

type Opcode byte

const (
	OPCODE_CONTINUATION     Opcode = 0x00
	OPCODE_TEXT             Opcode = 0x01
	OPCODE_BINARY           Opcode = 0x02
	OPCODE_CONNECTION_CLOSE Opcode = 0x08
	OPCODE_PING             Opcode = 0x09
	OPCODE_PONG             Opcode = 0x0A
)

type Frame struct {
	fin        byte
	opcode     Opcode
	mask       byte
	size       uint16
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
	// check len

	head := payload[:2]
	size := uint16(head[1] & 0x7F)
	frame := &Frame{payload: payload}

	if frame.IsPing() || frame.IsPong() {
		return frame, nil
	}

	if size < 126 {
		// check len
		frame.size = size
		frame.data = unmask(payload[2:6], payload[6:frame.size+6])
		return frame, nil
	}

	if size == 126 {
		// check len
		frame.size = binary.BigEndian.Uint16(payload[2:4])
		frame.data = unmask(payload[4:8], payload[8:frame.size+8])
		return frame, nil
	}

	// check len

	frame.size = binary.BigEndian.Uint16(payload[2:10])
	frame.data = unmask(payload[10:14], payload[14:frame.size+14])

	return frame, nil
}

const OPCODE_CONST Opcode = 128

// Comment
func Encode(opcode Opcode, data []byte) *Frame {
	size := len(data)
	frame := &Frame{data: data}
	opc := OPCODE_CONST + opcode

	if size < 126 {
		payload := make([]byte, 2)
		payload[0] = byte(opc)
		payload[1] = byte(size)

		payload = append(payload, data...)

		frame.payload = payload

		return frame
	}

	// if size == 126 || size < int(math.Pow(2, 16)) {
	// 	payload := make([]byte, 2)
	// 	payload[0] = byte(opcode + OPCODE_CONST)
	// 	payload[1] = 126

	// 	length := make([]byte, 2)

	// 	binary.BigEndian.PutUint16(length, uint16(size))

	// 	payload = append(payload, length...)
	// 	payload = append(payload, data...)

	// 	frame.payload = payload

	// 	return frame
	// }

	// payload := make([]byte, 2)
	// payload[0] = byte(opcode + OPCODE_CONST)
	// payload[1] = 127

	// length := make([]byte, 8)

	// binary.BigEndian.PutUint64(length, uint64(size))

	// payload = append(payload, length...)
	// payload = append(payload, data...)

	// frame.payload = payload

	return frame
}

// Comment
func (ctx *Frame) Opcode() Opcode {
	return Opcode(ctx.payload[0] - byte(OPCODE_CONST))
}

// Comment
func (ctx *Frame) IsContinuation() bool {
	// return (Opcode(ctx.payload[0]) & OPCODE_CONTINUATION) == OPCODE_CONTINUATION

	return ctx.Opcode() == OPCODE_CONTINUATION
}

// Comment
func (ctx *Frame) IsBinary() bool {
	// return (Opcode(ctx.payload[0]) & OPCODE_BINARY) == OPCODE_BINARY

	return ctx.Opcode() == OPCODE_BINARY
}

// Comment
func (ctx *Frame) IsText() bool {
	// return (Opcode(ctx.payload[0]) & OPCODE_TEXT) == OPCODE_TEXT

	return ctx.Opcode() == OPCODE_TEXT
}

// Comment
func (ctx *Frame) IsClose() bool {
	// return (Opcode(ctx.payload[0]) & OPCODE_CONNECTION_CLOSE) == OPCODE_CONNECTION_CLOSE

	return ctx.Opcode() == OPCODE_CONNECTION_CLOSE
}

// Comment
func (ctx *Frame) IsPing() bool {
	// return (Opcode(ctx.payload[0]) & OPCODE_PING) == OPCODE_PING

	return ctx.Opcode() == OPCODE_PING
}

// Comment
func (ctx *Frame) IsPong() bool {
	// return (Opcode(ctx.payload[0]) & OPCODE_PONG) == OPCODE_PONG

	return ctx.Opcode() == OPCODE_PONG
}

// Comment
func (ctx *Frame) Length() uint16 {
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

// fmt.Println("IsContinuation: ", frame.IsContinuation())
// fmt.Println("IsBinary: ", frame.IsBinary())
// fmt.Println("IsText: ", frame.IsText())
// fmt.Println("IsPing: ", frame.IsPing())
// fmt.Println("IsPong: ", frame.IsPong())
// fmt.Println("IsClose: ", frame.IsClose())
// fmt.Println("-------------------------------------------------------------------\n\n")
