package testing

import (
	"testing"

	"github.com/lucas11776-golang/http"
)

type TestCase struct {
	HTTP    *http.HTTP
	Testing *testing.T
}

// Comment
func (ctx *TestCase) Request() *Request {
	// return NewRequest(ctx)
	return nil
}

type Ws struct {
}

// Comment
func (ctx *TestCase) Ws() *Ws {
	return nil
}

// Comment
func (ctx *TestCase) Cleanup() {
	ctx.HTTP.Close()
}
