package request

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"

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
	return ctx.Request("GET", url, []byte{})
}

// Comment
func (ctx *Request) Post(url string, data []byte) (string, error) {
	return ctx.Request("POST", url, data)
}

// Comment
func (ctx *Request) PUT(url string, data []byte) (string, error) {
	return ctx.Request("PUT", url, data)
}

// Comment
func (ctx *Request) Patch(url string, data []byte) (string, error) {
	return ctx.Request("PATCH", url, data)
}

// Comment
func (ctx *Request) Delete(url string) (string, error) {
	return ctx.Request("DELETE", url, []byte{})
}

// Comment
func (ctx *Request) Options(url string) (string, error) {
	return ctx.Request("Options", url, []byte{})
}

// Comment
func (ctx *Request) Connect(url string, data []byte) (string, error) {
	return ctx.Request("Connect", url, data)
}

type Stream struct {
	Conn *net.Conn
}

// Comment
func (ctx *Request) Http2Request(method string, address string, data []byte) (http string, stream *Stream, err error) {
	return "", nil, nil
}

// Comment
func (ctx *Request) Request(method string, address string, data []byte) (string, error) {
	// TODO: Remove stream check if request is http2 or not...
	// TODO: Need to return http.Response...
	request, err := http.NewRequest(method, address, bytes.NewBuffer(data))

	if err != nil {
		return "", err
	}

	for k, v := range ctx.headers {
		request.Header.Set(k, v)
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return "", err
	}

	return ParseHttpResponse(response), nil
}

// Comment
func ParseHttpResponse(res *http.Response) string {
	http := []string{}

	http = append(http, strings.Join([]string{res.Proto, res.Status}, " "))

	keys := make([]string, 0, len(res.Header))

	body, err := io.ReadAll(res.Body)

	if err != nil {
		body = []byte{}
	}

	res.Header.Set("Content-Length", strconv.Itoa(len(body)))

	for k := range res.Header {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		k := cases.Title(language.English).String(key)
		v := strings.Join(res.Header[key], ";")

		http = append(http, strings.Join([]string{k, v}, ": "))
	}

	if len(body) == 0 {
		return strings.Join(append(http, "\r\n"), "\r\n")
	}

	return strings.Join(append(http, strings.Join([]string{"\r\n", string(body), "\r\n"}, "")), "\r\n")
}
