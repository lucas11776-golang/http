package request

import (
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/lucas11776-golang/http/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const MAX_RESPONSE_SIZE = 1024 * 1000

type Request struct {
	Conn            net.Conn
	protocal        string
	headers         types.Headers
	maxResponseSize int
}

// Comment
func CreateRequest() *Request {
	return &Request{
		protocal:        "HTTP/1.1",
		headers:         make(types.Headers),
		maxResponseSize: MAX_RESPONSE_SIZE,
	}
}

// Comment
func (ctx *Request) SetProtocal(protocal string) *Request {
	switch strings.ToUpper(protocal) {
	case "HTTP/1.1", "HTTP/2":
		ctx.protocal = strings.ToUpper(protocal)

	default:
		log.Fatalf("Request does not support protocal: %v", strings.ToUpper(protocal))
	}

	return ctx
}

// Comment
func (ctx *Request) SetHeader(key string, value string) *Request {
	ctx.headers[key] = value

	return ctx
}

// Comment
func (ctx *Request) SetHeaders(headers types.Headers) *Request {
	for k, v := range headers {
		ctx.headers[k] = v
	}

	return ctx
}

// Comment
func (ctx *Request) GetHeader(key string) string {
	header, ok := ctx.headers[key]

	if !ok {
		return ""
	}

	return header
}

// Comment
func (ctx *Request) Get(url string) (string, error) {
	http, _, err := ctx.Request("GET", url, []byte{})

	return http, err
}

// Comment
func (ctx *Request) Post(url string, data []byte) (string, error) {
	http, _, err := ctx.Request("POST", url, data)

	return http, err
}

// Comment
func (ctx *Request) PUT(url string, data []byte) (string, error) {
	http, _, err := ctx.Request("PUT", url, data)

	return http, err
}

// Comment
func (ctx *Request) Patch(url string, data []byte) (string, error) {
	http, _, err := ctx.Request("PATCH", url, data)

	return http, err
}

// Comment
func (ctx *Request) Delete(url string) (string, error) {
	http, _, err := ctx.Request("DELETE", url, []byte{})

	return http, err
}

// Comment
func (ctx *Request) Options(url string) (string, error) {
	http, _, err := ctx.Request("Options", url, []byte{})

	return http, err
}

// Comment
func (ctx *Request) Connect(url string, data []byte) (string, error) {
	http, _, err := ctx.Request("Connect", url, data)

	return http, err
}

// Comment
func (ctx *Request) parse(method string, path string, data []byte) string {
	if path == "" {
		path = "/"
	}

	arr := []string{
		strings.Join([]string{strings.ToUpper(method), path, "HTTP/1.1"}, " "),
	}

	for key, value := range ctx.headers {
		arr = append(arr, strings.Join([]string{cases.Title(language.English).String(key), value}, ": "))
	}

	if len(data) == 0 {
		return strings.Join(append(arr, "\r\n"), "\r\n")
	}

	arr = append(arr, strings.Join([]string{"Content-Length", strconv.Itoa(len(data))}, ": "))

	return strings.Join(append(arr, strings.Join([]string{"\r\n", string(data), "\r\n"}, "")), "\r\n")
}

type Stream struct {
	Conn *net.Conn
}

// Comment
func (ctx *Request) Http2Request(method string, address string, data []byte) (http string, stream *Stream, err error) {
	return "", nil, nil
}

// Comment
func (ctx *Request) Request(method string, address string, data []byte) (string, *Stream, error) {
	url, err := url.Parse(address)

	if err != nil {
		return "", nil, err
	}

	ctx.Conn, err = net.Dial("tcp", url.Host)

	if err != nil {
		return "", nil, err
	}

	_, err = ctx.Conn.Write([]byte(ctx.parse(method, url.Path, data)))

	if err != nil {
		return "", nil, err
	}

	err = ctx.Conn.SetDeadline(time.Now().Add(time.Second * 3))

	if err != nil {
		return "", nil, err
	}

	http := make([]byte, ctx.maxResponseSize)

	n, err := ctx.Conn.Read(http)

	if err != nil {
		return "", nil, err
	}

	return string(http[:n]), nil, nil
}
