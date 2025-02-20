package testing

import (
	"fmt"
	"testing"

	"github.com/lucas11776-golang/http"
)

type Values map[string]string

type TestCase struct {
	http    *http.HTTP
	testing *Testing
}

type Testing struct {
	t         *testing.T
	catchLogs bool
	errors    []string
}

// Comment
func NewTestCase(t *testing.T, HTTP *http.HTTP, catchLog bool) *TestCase {
	return &TestCase{
		http: HTTP,
		testing: &Testing{
			t:         t,
			catchLogs: catchLog,
		},
	}
}

// Comment
func (ctx *TestCase) Request() *Request {
	return NewRequest(ctx)
}

// Comment
func (ctx *TestCase) Ws() *Ws {
	return NewWs(ctx)
}

// Comment
func (ctx *TestCase) Cleanup() {
	ctx.http.Close()
}

// Comment
func (ctx *Testing) Log(args ...any) {
	if ctx.catchLogs {
		ctx.errors = append(ctx.errors, fmt.Sprintln(args...))

		return
	}

	ctx.t.Log(args...)
}

// Comment
func (ctx *Testing) Logf(format string, args ...any) {
	if ctx.catchLogs {
		ctx.errors = append(ctx.errors, fmt.Sprintf(format, args...))

		return
	}

	ctx.t.Logf(format, args...)
}

// Comment
func (ctx *Testing) Error(args ...any) {
	if ctx.catchLogs {
		ctx.errors = append(ctx.errors, fmt.Sprintln(args...))

		return
	}

	ctx.t.Error(args...)
}

// Comment
func (ctx *Testing) Errorf(format string, args ...any) {
	if ctx.catchLogs {
		ctx.errors = append(ctx.errors, fmt.Sprintf(format, args...))

		return
	}

	ctx.t.Errorf(format, args...)
}

// Comment
func (ctx *Testing) Fatal(args ...any) {
	if ctx.catchLogs {
		ctx.errors = append(ctx.errors, fmt.Sprintln(args...))

		return
	}

	ctx.t.Fatal(args...)
}

// Comment
func (ctx *Testing) Fatalf(format string, args ...any) {
	if ctx.catchLogs {
		ctx.errors = append(ctx.errors, fmt.Sprintf(format, args...))

		return
	}

	ctx.t.Fatalf(format, args...)
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
