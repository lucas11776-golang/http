package response

import "http/types"

type Response struct {
	status   int32
	protocol string
	headers  types.Headers
	body     []byte
}

// Comment
func Init() *Response {
	return &Response{}
}
