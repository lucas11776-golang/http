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

		testingFrame(t, frame, opcode, data)
	})

	t.Run("TestEncodeText", func(t *testing.T) {
		opcode := OPCODE_TEXT
		data := []byte("Text")
		frame := Encode(opcode, data)

		testingFrame(t, frame, opcode, data)
	})

	t.Run("TestEncodeBinary", func(t *testing.T) {
		opcode := OPCODE_BINARY
		data := []byte("Binary")
		frame := Encode(opcode, data)

		testingFrame(t, frame, opcode, data)
	})

	t.Run("TestEncodeClose", func(t *testing.T) {
		opcode := OPCODE_CONNECTION_CLOSE
		data := []byte("Close")
		frame := Encode(opcode, data)

		testingFrame(t, frame, opcode, data)
	})

	t.Run("TestEncodePing", func(t *testing.T) {
		opcode := OPCODE_PING
		data := []byte("Ping")
		frame := Encode(opcode, data)

		testingFrame(t, frame, opcode, data)
	})

	t.Run("TestEncodePong", func(t *testing.T) {
		opcode := OPCODE_PONG
		data := []byte("Pong")
		frame := Encode(opcode, data)

		testingFrame(t, frame, opcode, data)
	})

	// t.Run("TestEncodeDataEqualsTo126", func(t *testing.T) {
	// 	opcode := OPCODE_TEXT
	// 	data := generateData(512)
	// 	frame := Encode(opcode, data)

	// 	payload := []byte{byte(OPCODE_CONST + opcode), byte(len(data))}
	// 	payload = append(payload, data...)

	// 	testingFrame(t, frame, opcode, data)
	// })

	// t.Run("TestEncodeDataEqualGreatThe126", func(t *testing.T) {
	// 	opcode := OPCODE_PONG
	// 	data := []byte("Pong")
	// 	frame := Encode(opcode, data)

	// 	testingFrame(t, frame, opcode, data)
	// })

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

		testingFrame(t, frame, opcode, data)
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
func testingFrame(t *testing.T, frame *Frame, opcode Opcode, data []byte) {
	if frame.Length() != uint16(len(data)) {
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
