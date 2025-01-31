package frame

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestFrame(t *testing.T) {
	t.Run("TestEncodeContinuation", func(t *testing.T) {
		opcode := OPCODE_CONTINUATION
		data := []byte("Continuation")
		frame := Encode(opcode, data)

		testFrame(t, frame, opcode, data)
	})

	t.Run("TestEncodeText", func(t *testing.T) {
		opcode := OPCODE_TEXT
		data := []byte("Text")
		frame := Encode(opcode, data)

		testFrame(t, frame, opcode, data)
	})

	t.Run("TestEncodeBinary", func(t *testing.T) {
		opcode := OPCODE_BINARY
		data := []byte("Binary")
		frame := Encode(opcode, data)

		testFrame(t, frame, opcode, data)
	})

	t.Run("TestEncodeClose", func(t *testing.T) {
		opcode := OPCODE_CONNECTION_CLOSE
		data := []byte("Close")
		frame := Encode(opcode, data)

		testFrame(t, frame, opcode, data)
	})

	t.Run("TestEncodePing", func(t *testing.T) {
		opcode := OPCODE_PING
		data := []byte("Ping")
		frame := Encode(opcode, data)

		testFrame(t, frame, opcode, data)
	})

	t.Run("TestEncodePong", func(t *testing.T) {
		opcode := OPCODE_PONG
		data := []byte("Pong")
		frame := Encode(opcode, data)

		testFrame(t, frame, opcode, data)
	})

	t.Run("TestEncodeDataLessThen126", func(t *testing.T) {
		opcode := OPCODE_TEXT
		data := generateData(125)
		frame := Encode(opcode, data)

		testFrame(t, frame, opcode, data)

		if frame.Payload()[0] != byte(opcode+OPCODE_CONST) {
			t.Errorf("Expected payload opcode flag to be (%d) but got (%d)", opcode+OPCODE_CONST, frame.Payload()[0])
		}

		if string(data) != string(frame.Payload()[2:]) {
			t.Errorf("Expected payload data to be (%s) but got (%s)", string(data), string(frame.Payload()[2:]))
		}
	})

	t.Run("TestEncodeDataGreaterThenEqualsTo126AndLessThen2**16", func(t *testing.T) {
		opcode := OPCODE_TEXT
		data := generateData(512)
		frame := Encode(opcode, data)

		testFrame(t, frame, opcode, data)

		fmt.Println("SIZE: ")

		if frame.Payload()[0] != byte(opcode+OPCODE_CONST) {
			t.Errorf("Expected payload opcode flag to be (%d) but got (%d)", opcode+OPCODE_CONST, frame.Payload()[0])
		}

		if frame.Payload()[1] != 126 {
			t.Errorf("Expected payload length flag to be (%d) but got (%d)", 126, frame.Payload()[1])
		}

		size := binary.BigEndian.Uint16(frame.Payload()[2:4])

		if size != uint16(len(data)) {
			t.Errorf("Expected payload length to be (%d) but got (%d)", len(data), size)
		}

		if string(data) != string(frame.Payload()[4:]) {
			t.Errorf("Expected payload data to be (%s) but got (%s)", string(data), string(frame.Payload()[4:]))
		}
	})

	t.Run("TestEncodeDataGreaterThen2**16", func(t *testing.T) {
		opcode := OPCODE_BINARY
		data := generateData(int(math.Pow(2, 16) + 1))
		frame := Encode(opcode, data)

		testFrame(t, frame, opcode, data)

		fmt.Println("SIZE: ")

		if frame.Payload()[0] != byte(opcode+OPCODE_CONST) {
			t.Errorf("Expected payload opcode flag to be (%d) but got (%d)", opcode+OPCODE_CONST, frame.Payload()[0])
		}

		if frame.Payload()[1] != 126 {
			t.Errorf("Expected payload length flag to be (%d) but got (%d)", 127, frame.Payload()[1])
		}

		size := binary.BigEndian.Uint64(frame.Payload()[2:10])

		if size != uint64(len(data)) {
			t.Errorf("Expected payload length to be (%d) but got (%d)", len(data), size)
		}

		if string(data) != string(frame.Payload()[10:]) {
			t.Errorf("Expected payload data to be (%s) but got (%s)", string(data), string(frame.Payload()[10:]))
		}
	})

	// t.Run("TestEncodeDataEqualGreatThe126", func(t *testing.T) {
	// 	opcode := OPCODE_PONG
	// 	data := []byte("Pong")
	// 	frame := Encode(opcode, data)

	// 	testingFrame(t, frame, opcode, data)
	// })

	t.Run("TestDecodeWithInvalidPayload", func(t *testing.T) {
		_, err := Decode([]byte{byte(OPCODE_CONST + OPCODE_TEXT)})

		if err == nil {
			t.Errorf("Expected the decode to have error because the payload was invalid")
		}
	})

	t.Run("TestDecodeDataLessThe126", func(t *testing.T) {
		opcode := OPCODE_CONST + OPCODE_TEXT
		size := 125
		mask, data, dataMask := testingPayload(size)

		payload := []byte{byte(opcode), byte(size)}
		payload = append(payload, mask...)
		payload = append(payload, dataMask...)

		frame, err := Decode(payload)

		if err != nil {
			t.Errorf("Something went wrong when trying to decode payload: %s", err.Error())
		}

		testFrame(t, frame, opcode, data)

		// Invalid payload
		payloadInvalid := []byte{byte(opcode), byte(size)}
		payloadInvalid = append(payloadInvalid, mask...)
		payloadInvalid = append(payloadInvalid, []byte{65, 66, 67}...)

		_, err = Decode(payloadInvalid)

		if err == nil {
			t.Errorf("Expected the decode to have error because the payload was invalid")
		}
	})

	t.Run("TestDecodeDataEquals126AndLessThen2**16", func(t *testing.T) {
		// opcode := OPCODE_CONST + OPCODE_BINARY
		// size := 128
		// mask, data, dataMask := testingPayload(size)

		// payload := []byte{byte(opcode), 127}
	})

	t.Run("TestDecodeDataGreaterOrEqualsTo2**10", func(t *testing.T) {
		// Write Test
	})
}

// Comment
func generateData(size int) []byte {
	data := make([]byte, size)

	for i := range data {
		data[i] = byte(65 + (rand.Float32() * 58))
	}

	return data
}

// Comment
func testingPayload(size int) (mask []byte, data []byte, dataMasked []byte) {
	msk := make([]byte, 4)
	dt := make([]byte, size)
	dtMask := make([]byte, size)

	for i := range mask {
		msk[i] = byte(rand.Int() * 255)
	}

	for i := range dt {
		dt[i] = byte(65 + (rand.Float32() * 58))
		dtMask[i] = dt[i] ^ msk[i%len(msk)]
	}

	return msk, dt, dtMask
}

// Comment
func testFrame(t *testing.T, frame *Frame, opcode Opcode, data []byte) {
	if frame.Length() != uint64(len(data)) {
		t.Errorf("Expected data len in frame to be (%d) but got (%d)", len(data), frame.Length())
	}

	if string(data) != string(frame.Data()) {
		t.Errorf("Expected data in frame to be (%s) but got (%s)", string(data), string(frame.Data()))
	}

	if opcode == OPCODE_CONTINUATION && frame.IsContinuation() == false {
		t.Errorf("Expected frame opcode to be continuation")
	}

	if opcode == OPCODE_TEXT && frame.IsText() == false {
		t.Errorf("Expected frame opcode to be text")
	}

	if opcode == OPCODE_BINARY && frame.IsBinary() == false {
		t.Errorf("Expected frame opcode to be binary")
	}

	if opcode == OPCODE_CONNECTION_CLOSE && frame.IsClose() == false {
		t.Errorf("Expected frame opcode to be close")
	}

	if opcode == OPCODE_PING && frame.IsPing() == false {
		t.Errorf("Expected frame opcode to be ping")
	}

	if opcode == OPCODE_PONG && frame.IsPong() == false {
		t.Errorf("Expected frame opcode to be pong")
	}
}
