package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/lucas11776-golang/http/request"
	"github.com/lucas11776-golang/http/types"
)

type Header struct {
	method   string
	path     string
	protocol string
}

type Content struct {
	host    string
	headers types.Headers
	body    []byte
}

// Comment
func HttpHeader(http []string) (*Header, error) {
	header := strings.Split(http[0], " ")

	if len(header) != 3 {
		return nil, fmt.Errorf("Invalid http header: %s", http[0])
	}

	return &Header{
		method:   strings.ToUpper(header[0]),
		path:     header[1],
		protocol: strings.ToUpper(header[2]),
	}, nil
}

// Comment
func HttpContent(http []string) (*Content, error) {
	content := Content{
		headers: make(types.Headers),
	}

	for i, line := range http[1:] {
		if line == "" {
			content.body = []byte(strings.Trim(strings.Join(http[i:], "\r\n"), "\r\n"))
			break
		}

		header := strings.Split(line, ":")

		if len(header) < 2 {
			return nil, fmt.Errorf("Invalid header %s", header[0])
		}

		key := strings.ToLower(header[0])
		value := strings.Trim(strings.Join(header[1:], ":"), " ")

		if key == "host" {
			content.host = value
			continue
		}

		content.headers[key] = value
	}

	return &content, nil
}

// Comment
func ParseHttp(http string) (*request.Request, error) {
	hp := strings.Split(http, "\r\n")

	header, err := HttpHeader(hp)

	if err != nil {
		return nil, err
	}

	content, err := HttpContent(hp)

	if err != nil {
		return nil, err
	}

	req, err := request.Create(header.method, header.path, header.protocol, content.headers, bytes.NewReader(content.body))

	if err != nil {
		return nil, err
	}

	req.Host = content.host

	err = req.ParseForm()

	if err != nil {
		return nil, err
	}

	return req, nil
}
