package frame

import (
	"math/rand"
	"testing"
)

func TestFrame(t *testing.T) {
	t.Run("TestEncodeContinuation", func(t *testing.T) {
		opcode := OPCODE_CONTINUATION
		data := []byte("Continuation")
		frame := Encode(opcode, data)
		payload := append([]byte{byte(OPCODE_CONST + opcode), byte(len(data))}, data...)

		testingFrame(t, frame, opcode, data, &payload)
	})

	t.Run("TestEncodeText", func(t *testing.T) {
		opcode := OPCODE_TEXT
		data := []byte("Text")
		frame := Encode(opcode, data)
		payload := append([]byte{byte(OPCODE_CONST + opcode), byte(len(data))}, data...)

		testingFrame(t, frame, opcode, data, &payload)
	})

	t.Run("TestEncodeBinary", func(t *testing.T) {
		opcode := OPCODE_BINARY
		data := []byte("Binary")
		frame := Encode(opcode, data)
		payload := append([]byte{byte(OPCODE_CONST + opcode), byte(len(data))}, data...)

		testingFrame(t, frame, opcode, data, &payload)
	})

	t.Run("TestEncodeClose", func(t *testing.T) {
		opcode := OPCODE_CONNECTION_CLOSE
		data := []byte("Close")
		frame := Encode(opcode, data)
		payload := append([]byte{byte(OPCODE_CONST + opcode), byte(len(data))}, data...)

		testingFrame(t, frame, opcode, data, &payload)
	})

	t.Run("TestEncodePing", func(t *testing.T) {
		opcode := OPCODE_PING
		data := []byte("Ping")
		frame := Encode(opcode, data)
		payload := append([]byte{byte(OPCODE_CONST + opcode), byte(len(data))}, data...)

		testingFrame(t, frame, opcode, data, &payload)
	})

	t.Run("TestEncodePong", func(t *testing.T) {
		opcode := OPCODE_PONG
		data := []byte("Pong")
		frame := Encode(opcode, data)
		payload := append([]byte{byte(OPCODE_CONST + opcode), byte(len(data))}, data...)

		testingFrame(t, frame, opcode, data, &payload)
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

		if string(frame.Data()) != string(data) {
			t.Errorf("Expected frame data to be (%s) but go (%s)", data, frame.Data())
		}

		if frame.Opcode() != OPCODE_TEXT {
			t.Errorf("Expected frame opcode to be (%b) but go (%b)", OPCODE_TEXT, frame.Opcode())
		}
	})

	t.Run("TestDecodeDataEquals126AndLessThen2**16", func(t *testing.T) {
		// Write test
	})

	t.Run("TestDecodeDataGreaterOrEqualsTo2**10", func(t *testing.T) {
		// Write Test
	})
}

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
func testingFrame(t *testing.T, frame *Frame, opcode Opcode, data []byte, payload *[]byte) {
	if string(data) != string(frame.Data()) {
		t.Errorf("Expected data in frame to be (%s) but got (%s)", string(data), string(frame.Data()))
	}

	if len(data) != len(frame.Data()) {
		t.Errorf("Expected data len in frame to be (%d) but got (%d)", len(data), len(frame.Data()))
	}

	if string(data) != string(frame.Data()) {
		t.Errorf("Expected frame data to be (%s) but got (%s)", string(data), string(frame.Data()))
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

	if payload != nil && string(*payload) != string(frame.Payload()) {
		t.Errorf("Expected frame payload to be (%s) but go (%s)", *payload, frame.Payload())
	}
}
