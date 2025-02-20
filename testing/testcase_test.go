package testing

import (
	"encoding/json"
	"testing"

	"github.com/lucas11776-golang/http"
)

func TestTestCaseRequest(t *testing.T) {
	testcase := NewTestCase(t, http.Server("127.0.0.1", 0), false)

	body := "<h1>Home Page</h1>"

	testcase.http.Route().Get("/", func(req *http.Request, res *http.Response) *http.Response {
		return res.Html(body)
	})

	res := testcase.Request().Get("/")

	res.AssertOk()
	res.AssertHeader("content-type", "text/html")
	res.AssertBody([]byte(body))

	testcase.Cleanup()
}

func TestTestCaseWs(t *testing.T) {
	testcase := NewTestCase(t, http.Server("127.0.0.1", 0), false)

	type DeviceCommand struct {
		Name      string   `json:"name"`
		Arguments []string `json:"arguments"`
	}

	type DeviceCommandResponse struct {
		Command string `json:"command"`
		Message string `json:"message"`
	}

	testcase.http.Route().Ws("/devices/{device}/commands", func(req *http.Request, ws *http.Ws) {
		ws.OnReady(func(ws *http.Ws) {
			ws.OnMessage(func(data []byte) {
				command := &DeviceCommand{}

				if json.Unmarshal(data, command) != nil {
					return
				}

				switch command.Name {
				case "protection-lock":
					ws.WriteJson(&DeviceCommandResponse{
						Command: command.Name,
						Message: "Devices has been locked",
					})
					break
				default:
				}
			})
		})
	})

	ws := testcase.Ws().Connect("/devices/mobile-4235646342335/commands")

	ws.WriteJson(&DeviceCommand{
		Name:      "protection-lock",
		Arguments: []string{"unlock-face-recognition-only", "geolocation-on"},
	})

	ws.AssertJson(&DeviceCommandResponse{
		Command: "protection-lock",
		Message: "Devices has been locked",
	})

	testcase.Cleanup()
}
