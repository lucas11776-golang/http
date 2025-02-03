package request

import (
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/lucas11776-golang/http/types"
)

const MAX_REQUEST_SIZE = 1024 * 1000

type Request struct {
	method  string
	headers types.Headers
	data    []byte
}

// Comment
func CreateRequest() *Request {
	return &Request{
		headers: make(types.Headers),
	}
}

// Comment
func (ctx *Request) Header(key string, value string) *Request {
	ctx.headers[key] = value

	return ctx
}

// Comment
func (ctx *Request) Post(url string, data []byte) (string, error) {
	return ctx.Request("POST", url, data)
}

// Comment
func (ctx *Request) Get(url string) (string, error) {
	return ctx.Request("GET", url, []byte{})
}

// Comment
func (ctx *Request) parse(method string, path string, data []byte) string {
	arr := []string{
		strings.Join([]string{strings.ToUpper(method), path, "github.com/lucas11776-golang/http/1.1"}, " "),
	}

	for key, value := range ctx.headers {
		arr = append(arr, strings.Join([]string{key, value}, ": "))
	}

	if len(data) == 0 {
		return strings.Join(append(arr, "\r\n"), "\r\n")
	}

	arr = append(arr, strings.Join([]string{"Content-Length", strconv.Itoa(len(data))}, ": "))

	return strings.Join(append(arr, strings.Join([]string{"\r\n", string(data), "\r\n"}, "")), "\r\n")
}

// Comment
func (ctx *Request) Request(method string, address string, data []byte) (string, error) {
	url, err := url.Parse(address)

	if err != nil {
		return "", err
	}

	listener, err := net.Dial("tcp", url.Host)

	if err != nil {
		return "", err
	}

	_, err = listener.Write([]byte(ctx.parse(method, url.Path, data)))

	if err != nil {
		return "", err
	}

	err = listener.SetDeadline(time.Now().Add(time.Second * 3))

	if err != nil {
		return "", err
	}

	http := make([]byte, MAX_REQUEST_SIZE)

	n, err := listener.Read(http)

	if err != nil {
		return "", err
	}

	listener.Close()

	return string(http[:n]), nil
}
