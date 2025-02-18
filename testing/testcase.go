package testing

import (
	"fmt"
	"testing"

	"github.com/lucas11776-golang/http"
)

type TestCase struct {
	HTTP    *http.HTTP
	Testing *Testing
}

type Testing struct {
	T         *testing.T
	catchLogs bool
	errors    []string
}

// Comment
func NewTestCase(t *testing.T, HTTP *http.HTTP, catchLog bool) *TestCase {
	return &TestCase{
		HTTP: HTTP,
		Testing: &Testing{
			T:         t,
			catchLogs: catchLog,
		},
	}
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

// Comment
func (ctx *Testing) Log(args ...any) {
	if ctx.catchLogs {
		ctx.errors = append(ctx.errors, fmt.Sprintln(args...))

		return
	}

	ctx.T.Log(args...)
}

// Comment
func (ctx *Testing) Logf(format string, args ...any) {
	if ctx.catchLogs {
		ctx.errors = append(ctx.errors, fmt.Sprintf(format, args...))

		return
	}

	ctx.T.Logf(format, args...)
}

// Comment
func (ctx *Testing) Error(args ...any) {
	if ctx.catchLogs {
		ctx.errors = append(ctx.errors, fmt.Sprintln(args...))

		return
	}

	ctx.T.Error(args...)
}

// Comment
func (ctx *Testing) Errorf(format string, args ...any) {
	if ctx.catchLogs {
		ctx.errors = append(ctx.errors, fmt.Sprintf(format, args...))

		return
	}

	ctx.T.Errorf(format, args...)
}

// Comment
func (ctx *Testing) Fatal(args ...any) {
	if ctx.catchLogs {
		ctx.errors = append(ctx.errors, fmt.Sprintln(args...))

		return
	}

	ctx.T.Fatal(args...)
}

// Comment
func (ctx *Testing) Fatalf(format string, args ...any) {
	if ctx.catchLogs {
		ctx.errors = append(ctx.errors, fmt.Sprintf(format, args...))

		return
	}

	ctx.T.Fatalf(format, args...)
}

// Comment
func (ctx *Testing) hasError() bool {
	return len(ctx.errors) > 0
}

// Comment
func (ctx *Testing) popError() string {
	if !ctx.hasError() {
		return ""
	}

	err := ctx.errors[len(ctx.errors)-1]

	ctx.errors = ctx.errors[:len(ctx.errors)-1]

	return err
}
