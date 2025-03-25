package http

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"testing"
	"time"

	"golang.org/x/net/http2"
)

//

func GenerateCert(host string) (cert string, key string, err error) {
	k, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		// return "", "", err
	}

	keyBytes := x509.MarshalPKCS1PrivateKey(k)

	// PEM encoding of private key
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyBytes,
		},
	)

	// 	openssl genrsa 2048 > host.key
	// chmod 400 host.key
	// openssl req -new -x509 -nodes -sha256 -days 365 -key host.key -out host.cert

	notBefore := time.Now().Add(-(time.Second * 20))
	notAfter := notBefore.Add(365 * 24 * 10 * time.Hour)

	//Create certificate templet
	template := x509.Certificate{
		SerialNumber:          big.NewInt(0),
		Subject:               pkix.Name{CommonName: "testing"},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	//Create certificate using templet
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &k.PublicKey, k)

	if err != nil {
		// return "", "", err

	}

	//pem encoding of certificate
	certPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: derBytes,
		},
	))

	return certPem, string(keyPEM), nil
}

func TestHtt2(t *testing.T) {
	response := struct {
		Message string `json:"message"`
	}{Message: fmt.Sprintf("Hello: %f", 2321412.3214)}

	serve := func() *HTTP {
		cert, key, err := GenerateCert("127.0.0.1")

		if err != nil {
			t.Fatalf("TLS error: %v", err)
		}

		fmt.Println(cert, "\r\n", key)

		server := ServerTLS("127.0.0.1", 0, cert, key)

		server.Route().Get("/", func(req *Request, res *Response) *Response {
			return res.SetStatus(HTTP_RESPONSE_OK).Json(response)
		})

		go server.Listen()

		return server
	}

	t.Run("HTTP2NegotiationRequest", func(t *testing.T) {
		server := serve()

		// req := request.CreateRequest().
		// 	SetProtocal("HTTP/1.1").
		// 	SetHeaders(types.Headers{
		// 		"Connection":     "Upgrade, HTTP2-Settings",
		// 		"Upgrade":        "h2c",
		// 		"HTTP2-Settings": "AAAACCCCAAAAZZZZZ"})

		// http, err := req.Get(fmt.Sprintf("http://%s:%d", server.Address(), server.Port()))

		// if err != nil {
		// 	t.Fatalf("Failed to send request: %v", err)
		// }

		// res, err := HttpToResponse(http)

		// if err != nil {
		// 	t.Fatalf("Failed to parse http response: %v", err)
		// }

		// if Status(res.StatusCode) != HTTP_RESPONSE_SWITCHING_PROTOCOLS {
		// 	t.Fatalf("Expected status code to be (%d) but got (%d)", HTTP_RESPONSE_SWITCHING_PROTOCOLS, res.StatusCode)
		// }

		// if res.GetHeader("Connection") != "Upgrade" {
		// 	t.Fatalf("Expected connection header to be (%s) but got (%s)", "Upgrade", res.GetHeader("Connection"))
		// }

		// if res.GetHeader("Upgrade") != "HTTP/2.0" {
		// 	t.Fatalf("Expected upgrade header to be (%s) but got (%s)", "HTTP/2.0", res.GetHeader("Upgrade"))
		// }

		// Create an HTTP/2 Transport
		// transport :=

		fmt.Println()

		fmt.Println("ADDRESS", fmt.Sprintf("http://%s", server.Host()))

		// Create an HTTP client with the HTTP/2 transport
		client := &http.Client{Transport: &http2.Transport{
			AllowHTTP: true,
		}}

		// Make an HTTP/2 GET request
		resp, err := client.Get(fmt.Sprintf("http://%s", server.Host()))

		if err != nil {
			fmt.Println("Error making request:", err, resp)
			return
		}

		// defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		fmt.Println("Response status:", resp.Status)
		fmt.Println("Response body:", string(body))

		server.Close()
	})

}

// > GET / HTTP/1.1
// > Host: 127.0.0.1:6666
// > Connection: Upgrade, HTTP2-Settings
// > Upgrade: h2c
// > HTTP2-Settings: ZZZZZZZZZZZZZZZZZ

// Connection: Upgrade
// Upgrade: HTTP/2.0
