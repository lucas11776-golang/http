package testing

import (
	"encoding/json"
	"math"
	"math/rand"
	"testing"

	"github.com/lucas11776-golang/http"
	"github.com/lucas11776-golang/http/utils/strings"
)

func TestTestingWs(t *testing.T) {
	ws := NewWs(NewTestCase(t, http.Server("127.0.0.1", 0).SetMaxWebsocketPayload(int(math.Pow(2, 17))), true))

	type coordinate struct {
		Longitude float32 `json:"longitude"`
		Latitude  float32 `json:"latitude"`
		Altitude  int     `json:"altitude"`
	}

	position := &coordinate{
		Longitude: rand.Float32() * 360,
		Latitude:  rand.Float32() * 180,
		Altitude:  rand.Int() * 10,
	}

	ws.testcase.http.Route().Ws("position", func(req *http.Request, ws *http.Ws) {
		ws.OnReady(func(ws *http.Ws) {
			ws.OnMessage(func(data []byte) {
				ws.Write(data)
			})
		})
	})

	go func() { ws.testcase.http.Listen() }()

	res := ws.Connect("position")

	fake := &coordinate{
		Longitude: rand.Float32() * 360,
		Latitude:  rand.Float32() * 180,
		Altitude:  rand.Int() * 10,
	}

	// Invalid position
	res.WriteJson(position)

	payloadFake, _ := json.Marshal(fake)

	res.AssertRead(payloadFake)

	if !ws.testing.hasError() {
		t.Fatalf("Expected assert read to log")
	}

	res.testcase.testing.popError()

	// Valid position
	res.WriteJson(position)

	payload, _ := json.Marshal(position)

	res.AssertRead(payload)

	if ws.testing.hasError() {
		t.Fatalf("Expected assert read to not log")
	}

	// Length greater then 126 less then 2^16
	data360 := strings.Random(360)

	res.WriteText([]byte(data360))

	res.AssertRead([]byte(data360))

	if ws.testing.hasError() {
		t.Fatalf("Expected assert read to not log")
	}

	// Length greater then 126 less then 2^16
	data2__16 := strings.Random(int(math.Pow(2, 16)))

	res.WriteText([]byte(data2__16))

	res.AssertRead([]byte(data2__16))

	if ws.testing.hasError() {
		t.Fatalf("Expected assert read to not log")
	}

	ws.testcase.Cleanup()
}

func TestWsSession(t *testing.T) {

}
