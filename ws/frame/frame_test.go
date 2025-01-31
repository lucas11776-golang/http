package frame

import (
	"encoding/binary"
	"math"
	"math/rand"
	"testing"
)

func TestFrame(t *testing.T) {
	t.Run("TestEncodeDataLessThen126", func(t *testing.T) {
		opcode := OPCODE_TEXT
		data := generateData(125)
		frame := Encode(opcode, data)

		testFrame(t, frame, opcode, data)

		if frame.Payload()[0] != byte(opcode+OPCODE_START) {
			t.Fatalf("Expected payload opcode flag to be (%d) but got (%d)", opcode+OPCODE_START, frame.Payload()[0])
		}

		if frame.Payload()[1] != byte(len(data)) {
			t.Fatalf("Expected payload length flag to be (%d) but got (%d)", 126, frame.Payload()[1])
		}

		if string(data) != string(frame.Payload()[2:]) {
			t.Fatalf("Expected payload data to be (%s) but got (%s)", string(data), string(frame.Payload()[2:]))
		}
	})

	t.Run("TestEncodeDataGreaterThenEqualsTo126AndLessThen2**16", func(t *testing.T) {
		opcode := OPCODE_TEXT
		data := generateData(512)
		frame := Encode(opcode, data)

		testFrame(t, frame, opcode, data)

		if frame.Payload()[0] != byte(opcode+OPCODE_START) {
			t.Fatalf("Expected payload opcode flag to be (%d) but got (%d)", opcode+OPCODE_START, frame.Payload()[0])
		}

		if frame.Payload()[1] != 126 {
			t.Fatalf("Expected payload length flag to be (%d) but got (%d)", 126, frame.Payload()[1])
		}

		size := binary.BigEndian.Uint16(frame.Payload()[2:4])

		if size != uint16(len(data)) {
			t.Fatalf("Expected payload length to be (%d) but got (%d)", len(data), size)
		}

		if string(data) != string(frame.Payload()[4:]) {
			t.Fatalf("Expected payload data to be (%s) but got (%s)", string(data), string(frame.Payload()[4:]))
		}
	})

	t.Run("TestEncodeDataGreaterThen2**16", func(t *testing.T) {
		opcode := OPCODE_BINARY
		data := generateData(int(math.Pow(2, 16) + 1))
		frame := Encode(opcode, data)

		testFrame(t, frame, opcode, data)

		if frame.Payload()[0] != byte(opcode+OPCODE_START) {
			t.Fatalf("Expected payload opcode flag to be (%d) but got (%d)", opcode+OPCODE_START, frame.Payload()[0])
		}

		if frame.Payload()[1] != 126 {
			t.Fatalf("Expected payload length flag to be (%d) but got (%d)", 127, frame.Payload()[1])
		}

		size := binary.BigEndian.Uint64(frame.Payload()[2:10])

		if size != uint64(len(data)) {
			t.Fatalf("Expected payload length to be (%d) but got (%d)", len(data), size)
		}

		if string(data) != string(frame.Payload()[10:]) {
			t.Fatalf("Expected payload data to be (%s) but got (%s)", string(data), string(frame.Payload()[10:]))
		}
	})

	t.Run("TestEncodeContinuation", func(t *testing.T) {
		opcode := OPCODE_CONTINUATION
		data := []byte("Continuation")
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

	t.Run("TestDecodeWithInvalidPayload", func(t *testing.T) {
		_, err := Decode([]byte{byte(OPCODE_START + OPCODE_TEXT)})

		if err == nil {
			t.Fatalf("Expected the decode to have error because the payload was invalid")
		}
	})

	t.Run("TestDecodeDataLessThe126", func(t *testing.T) {
		opcode := OPCODE_START + OPCODE_TEXT
		mask, data, dataMask := testingPayload(125)

		payload := []byte{byte(opcode), byte(len(data))}
		payload = append(payload, mask...)
		payload = append(payload, dataMask...)

		frame, err := Decode(payload)

		if err != nil {
			t.Fatalf("Something went wrong when trying to decode payload: %s", err.Error())
		}

		testFrame(t, frame, opcode, data)
	})

	t.Run("TestDecodeDataGreaterThenEquals126AndLessThen2**16", func(t *testing.T) {
		opcode := OPCODE_START + OPCODE_TEXT
		mask, data, dataMask := testingPayload(1024)

		payload := []byte{byte(opcode), 126}

		size := make([]byte, 2)

		binary.BigEndian.PutUint16(size, uint16(len(data)))

		payload = append(payload, size...)
		payload = append(payload, mask...)
		payload = append(payload, dataMask...)

		frame, err := Decode(payload)

		if err != nil {
			t.Fatalf("Something went wrong when trying to decode payload: %s", err.Error())
		}

		testFrame(t, frame, opcode, data)
	})

	t.Run("TestDecodeDataGreaterOrEqualsTo2**16", func(t *testing.T) {
		opcode := OPCODE_START + OPCODE_TEXT
		mask, data, dataMask := testingPayload(int(math.Pow(2, 16) + 1))

		payload := []byte{byte(opcode), 127}

		size := make([]byte, 8)

		binary.BigEndian.PutUint64(size, uint64(len(data)))

		payload = append(payload, size...)
		payload = append(payload, mask...)
		payload = append(payload, dataMask...)

		frame, err := Decode(payload)

		if err != nil {
			t.Fatalf("Something went wrong when trying to decode payload: %s", err.Error())
		}

		testFrame(t, frame, opcode, data)
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
		t.Fatalf("Expected data len in frame to be (%d) but got (%d)", len(data), frame.Length())
	}

	if string(data) != string(frame.Data()) {
		t.Fatalf("Expected data in frame to be (%s) but got (%s)", string(data), string(frame.Data()))
	}

	if opcode == OPCODE_CONTINUATION && frame.IsContinuation() == false {
		t.Fatalf("Expected frame opcode to be continuation")
	}

	if opcode == OPCODE_TEXT && frame.IsText() == false {
		t.Fatalf("Expected frame opcode to be text")
	}

	if opcode == OPCODE_BINARY && frame.IsBinary() == false {
		t.Fatalf("Expected frame opcode to be binary")
	}

	if opcode == OPCODE_CONNECTION_CLOSE && frame.IsClose() == false {
		t.Fatalf("Expected frame opcode to be close")
	}

	if opcode == OPCODE_PING && frame.IsPing() == false {
		t.Fatalf("Expected frame opcode to be ping")
	}

	if opcode == OPCODE_PONG && frame.IsPong() == false {
		t.Fatalf("Expected frame opcode to be pong")
	}
}
