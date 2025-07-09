package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/lucas11776-golang/http"
	"github.com/lucas11776-golang/http/types"
	"github.com/lucas11776-golang/orm"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Response struct {
	testcase *TestCase
	request  *Request
	Response *http.Response
	testing  *Testing
}

// Comment
func NewResponse(req *Request, res *http.Response) *Response {
	return &Response{
		testcase: req.testCase,
		request:  req,
		Response: res,
		testing:  req.testCase.testing,
	}
}

// Comment
func (ctx *Response) AssertProtocol(protocol string) *Response {
	if ctx.Response.Protocol() != protocol {
		ctx.testing.Fatalf("Expected response protocol to be (%s) but got (%s)", ctx.Response.Protocol(), protocol)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertStatusCode(status http.Status) *Response {
	if ctx.Response.StatusCode != int(status) {
		ctx.testing.Fatalf("Expected response status code to be (%d) but got (%d)", status, ctx.Response.StatusCode)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertOk() *Response {
	if !(ctx.Response.StatusCode >= int(http.HTTP_RESPONSE_OK) && ctx.Response.StatusCode <= int(http.HTTP_RESPONSE_ACCEPTED)) {
		ctx.testing.Fatalf("Expected response status code to be (200, 201, 202) but got (%d)", ctx.Response.StatusCode)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertNotFound() *Response {
	return ctx.AssertStatusCode(http.HTTP_RESPONSE_NOT_FOUND)
}

// Comment
func (ctx *Response) AssertUnauthorized() *Response {
	return ctx.AssertStatusCode(http.HTTP_RESPONSE_UNAUTHORIZED)
}

// Comment
func (ctx *Response) AssertHeadersHas(header string) *Response {
	_, ok := ctx.Response.Header[cases.Title(language.English).String(header)]

	if !ok {
		ctx.testing.Fatalf("Expected response to contain header (%s)", header)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertHeader(header string, value string) *Response {
	if ctx.Response.GetHeader(header) != value {
		ctx.testing.Fatalf("Expected response header (%s) to be (%s) but got (%s)", header, value, ctx.Response.GetHeader(header))
	}

	return ctx
}

// Comment
func (ctx *Response) AssertHeaders(headers types.Headers) *Response {
	for k, v := range headers {
		ctx.AssertHeader(k, v)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertBody(body []byte) *Response {
	tBody, err := io.ReadAll(ctx.Response.Body)

	if err != nil {
		ctx.testing.Fatalf("Something went wrong when trying to read body: %v", err)
	}

	if string(tBody) != string(body) {
		ctx.testing.Fatalf("Expected response body to be (%s) but got (%s)", string(body), string(tBody))
	}

	return ctx
}

// Comment
func (ctx *Response) AssertIsRedirect() *Response {
	if ctx.Response.Bag.Redirect == nil {
		ctx.testing.Fatalf("Expected response to be redirect")
	}

	return ctx
}

// Comment
func (ctx *Response) AssertRedirectTo(path string) *Response {
	ctx.AssertIsRedirect()

	if strings.Trim(path, "/") != strings.Trim(ctx.Response.Bag.Redirect.To, "/") {
		ctx.testing.Fatalf(
			"Expected redirect path to be (%s) but go (%s)",
			strings.Trim(path, "/"),
			strings.Trim(ctx.Response.Bag.Redirect.To, "/"),
		)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertIsView() *Response {
	if ctx.Response.Bag.View == nil {
		ctx.testing.Fatalf("Expected response to be view")
	}

	return ctx
}

// Comment
func (ctx *Response) AssertView(view string) *Response {
	ctx.AssertIsView()

	if view != ctx.Response.Bag.View.Name {
		ctx.testing.Fatalf("Expected view to be (%s) but go (%s)", view, ctx.Response.Bag.View.Name)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertViewHas(keys []string) *Response {
	ctx.AssertIsView()

	for _, key := range keys {
		_, ok := ctx.Response.Bag.View.Data[key]

		if !ok {
			ctx.testing.Fatalf("Expected view data to have (%s)", key)
		}
	}

	return ctx
}

// Comment
func (ctx *Response) AssertSessionHas(keys []string) *Response {
	for _, key := range keys {
		if ctx.Response.Session.Get(key) == "" {
			ctx.testing.Fatalf("Expected session to have (%s)", key)
		}
	}

	return ctx
}

// Comment
func (ctx *Response) AssertSession(key string, value string) *Response {
	if ctx.Response.Session.Get(key) != value {
		ctx.testing.Fatalf(
			"Expected session %s to but (%s) but got (%s)",
			key,
			value,
			ctx.Response.Session.Get(key),
		)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertSessionErrorsHas(keys []string) *Response {
	for _, key := range keys {
		if ctx.Response.Session.Error(key) == "" {
			ctx.testing.Fatalf("Expected session errors to have (%s)", key)
		}
	}

	return ctx
}

// Comment
func (ctx *Response) AssertSessionError(key string, value string) *Response {
	if err := ctx.Response.Session.Error(key); err != value {
		ctx.testing.Fatalf("Expected session error to be (%s) but got (%s)", value, err)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertSessionErrors(errs map[string]string) *Response {
	for k, v := range errs {
		ctx.AssertSessionError(k, v)
	}

	return ctx
}

// Comment
func (ctx *Response) parseResponseJsonErrorsBody() (*http.JsonErrorResponse, error) {
	if contentType := ctx.Response.GetHeader("content-type"); contentType != "application/json" {
		return nil, fmt.Errorf("Response content type is not application/json but is %s", contentType)
	}

	if code := ctx.Response.StatusCode; code < int(http.HTTP_RESPONSE_BAD_REQUEST) {
		return nil, fmt.Errorf("Response status code is not error code status %d", code)
	}

	body, err := io.ReadAll(ctx.Response.Body)

	if err != nil {
		return nil, err
	}

	var errs http.JsonErrorResponse

	if err := json.Unmarshal(body, &errs); err != nil {
		return nil, err
	}

	ctx.Response.Body = io.NopCloser(bytes.NewBuffer(body))

	return &errs, nil
}

// Comment
func (ctx *Response) AssertJsonErrorsHas(keys []string) *Response {
	errs, err := ctx.parseResponseJsonErrorsBody()

	if err != nil {
		ctx.testing.Fatalf("Failed to read json response body - %v", err)

		return ctx
	}

	for _, key := range keys {
		if _, ok := errs.Errors[key]; !ok {
			ctx.testing.Fatalf("Expected json errors to have (%s)", key)
		}
	}

	return ctx
}

// Comment
func (ctx *Response) AssertJsonError(key string, value string) *Response {
	errs, err := ctx.parseResponseJsonErrorsBody()

	if err != nil {
		ctx.testing.Fatalf("Failed to read json response body - %v", err)

		return ctx
	}

	if err := errs.Errors[key]; err != value {
		ctx.testing.Fatalf("Expected json error to be (%s) but got (%s)", value, err)
	}

	return ctx
}

// Comment
func (ctx *Response) AssertJsonErrors(errors map[string]string) *Response {
	for k, v := range errors {
		ctx.AssertJsonError(k, v)
	}

	return ctx
}

// Comment
func (ctx *Response) existsInDatabase(connection string, table string, values map[string]interface{}) (int64, error) {
	db := orm.DB.Database(connection)

	if db == nil {
		return 0, fmt.Errorf("database connection %s does not exists", connection)
	}

	where := []interface{}{}

	index := 0

	for k, v := range values {
		index += 1

		if index%2 == 0 {
			where = append(where, orm.AND)
		}

		where = append(where, &orm.Where{Key: k, Operator: orm.EQUALS, Value: v})
	}

	return db.Count(&orm.Statement{
		Table: table,
		Where: where,
	})
}

// Comment
func (ctx *Response) AssertDatabaseHas(connection string, table string, values map[string]interface{}) *Response {
	count, err := ctx.existsInDatabase(connection, table, values)

	if err != nil {
		ctx.testing.Fatal(err)

		return ctx
	}

	if count == 0 {
		data, _ := json.Marshal(values)

		ctx.testing.Fatalf("Record (%s) does not exists in database", string(data))

		return ctx
	}

	return ctx
}

// Comment
func (ctx *Response) AssertDatabaseMissing(connection string, table string, values map[string]interface{}) *Response {
	count, err := ctx.existsInDatabase(connection, table, values)

	if err != nil {
		ctx.testing.Fatal(err)

		return ctx
	}

	if count != 0 {
		data, _ := json.Marshal(values)

		ctx.testing.Fatalf("Record (%s) exists in database", string(data))

		return ctx
	}

	return ctx
}

// Comment
func (ctx *Response) AssertDatabaseCount(connection string, table string, size int64) *Response {
	count, err := ctx.existsInDatabase(connection, table, map[string]interface{}{})

	if err != nil {
		ctx.testing.Fatal(err)

		return ctx
	}

	if count != size {
		ctx.testing.Fatalf("Expected table to have (%d) records but got (%d)", size, count)

		return ctx
	}

	return ctx
}
