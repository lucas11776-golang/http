package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/lucas11776-golang/http/pages"
	"github.com/lucas11776-golang/http/types"
	h "github.com/lucas11776-golang/http/utils/headers"
	"github.com/lucas11776-golang/http/utils/response"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Status int

const (
	HTTP_RESPONSE_CONTINUE                        Status = 100
	HTTP_RESPONSE_SWITCHING_PROTOCOLS             Status = 101
	HTTP_RESPONSE_PROCESSING                      Status = 102
	HTTP_RESPONSE_EARLY_HINTS                     Status = 103
	HTTP_RESPONSE_OK                              Status = 200
	HTTP_RESPONSE_CREATED                         Status = 201
	HTTP_RESPONSE_ACCEPTED                        Status = 202
	HTTP_RESPONSE_NON_AUTHORITATIVE_INFORMATION   Status = 203
	HTTP_RESPONSE_NO_CONTENT                      Status = 204
	HTTP_RESPONSE_RESET_CONTENT                   Status = 205
	HTTP_RESPONSE_PARTIAL_CONTENT                 Status = 206
	HTTP_RESPONSE_MULTI_STATUS                    Status = 207
	HTTP_RESPONSE_ALREADY_REPORTED                Status = 208
	HTTP_RESPONSE_IM_USED                         Status = 226
	HTTP_RESPONSE_MULTIPLE_CHOICES                Status = 300
	HTTP_RESPONSE_MOVE_PERMANENTLY                Status = 301
	HTTP_RESPONSE_FOUND                           Status = 302
	HTTP_RESPONSE_SEE_OTHER                       Status = 303
	HTTP_RESPONSE_NOT_MODIFIED                    Status = 304
	HTTP_RESPONSE_USE_PROXY                       Status = 305
	HTTP_RESPONSE_UNUSED                          Status = 306
	HTTP_RESPONSE_TEMPORARY_REDIRECT              Status = 307
	HTTP_RESPONSE_PERMANENT_REDIRECT              Status = 308
	HTTP_RESPONSE_BAD_REQUEST                     Status = 400
	HTTP_RESPONSE_UNAUTHORIZED                    Status = 401
	HTTP_RESPONSE_PAYMENT_REQUIRED                Status = 402
	HTTP_RESPONSE_FORBIDDEN                       Status = 403
	HTTP_RESPONSE_NOT_FOUND                       Status = 404
	HTTP_RESPONSE_METHOD_NOT_ALLOWED              Status = 405
	HTTP_RESPONSE_NOT_ACCEPTABLE                  Status = 406
	HTTP_RESPONSE_PROXY_AUTHENTICATION_REQUIRED   Status = 407
	HTTP_RESPONSE_REQUEST_TIMEOUT                 Status = 408
	HTTP_RESPONSE_CONFLICT                        Status = 409
	HTTP_RESPONSE_GONE                            Status = 410
	HTTP_RESPONSE_LENGTH_REQUIRED                 Status = 411
	HTTP_RESPONSE_PRECONDITION_FAILED             Status = 412
	HTTP_RESPONSE_CONTENT_TOO_LARGER              Status = 413
	HTTP_RESPONSE_URI_TOO_LARGE                   Status = 414
	HTTP_RESPONSE_UNSUPPORTED_MEDIA_TYPE          Status = 415
	HTTP_RESPONSE_RANGE_NOT_SATISFIABLE           Status = 416
	HTTP_RESPONSE_EXPECTATION_FAILED              Status = 417
	HTTP_RESPONSE_IM_A_TEAPOT                     Status = 418
	HTTP_RESPONSE_MISDIRECTED_REQUEST             Status = 421
	HTTP_RESPONSE_UNPROCESSABLE_CONTENT           Status = 422
	HTTP_RESPONSE_LOCKED                          Status = 423
	HTTP_RESPONSE_FAILED_DEPENDENCY               Status = 424
	HTTP_RESPONSE_TOO_EARLY                       Status = 425
	HTTP_RESPONSE_UPGRADE_REQUIRED                Status = 426
	HTTP_RESPONSE_PRECONDITION_REQUIRED           Status = 428
	HTTP_RESPONSE_TOO_MANY_REQUIRED               Status = 429
	HTTP_RESPONSE_REQUEST_HEADER_FIELD_TOO_LARGE  Status = 431
	HTTP_RESPONSE_UNAVAILABLE_FOR_LEGAL_REASONS   Status = 451
	HTTP_RESPONSE_INTERNAL_SERVER_ERROR           Status = 500
	HTTP_RESPONSE_NOT_IMPLEMENTED                 Status = 501
	HTTP_RESPONSE_BAD_GATEWAY                     Status = 502
	HTTP_RESPONSE_SERVICE_UNAVAILABLE             Status = 503
	HTTP_RESPONSE_GATEWAY_TIMEOUT                 Status = 504
	HTTP_RESPONSE_HTTP_VERSION_NOT_SUPPORTED      Status = 505
	HTTP_RESPONSE_VARIANT_ALSO_NEGOTIATES         Status = 506
	HTTP_RESPONSE_INSUFFICIENT_STORAGE            Status = 507
	HTTP_RESPONSE_LOOP_DETECTED                   Status = 508
	HTTP_RESPONSE_NOT_EXTENDED                    Status = 510
	HTTP_RESPONSE_NETWORK_AUTHENTICATION_REQUIRED Status = 511
)

var (
	ErrHttpResponse = errors.New("invalid http response")
)

type RedirectBag struct {
	To string
}

type ViewBag struct {
	Name string
	Data ViewData
}

type Bag struct {
	View     *ViewBag
	Redirect *RedirectBag
}

type Response struct {
	*http.Response
	Writer  http.ResponseWriter
	Request *Request
	Session SessionManager
	Bag     *Bag
	Ws      *Ws
}

type Writer struct {
	response *Response
}

// Comment
func (ctx *Response) Protocol() string {
	return ctx.Proto
}

// Comment
func (ctx *Writer) Header() http.Header {
	return ctx.response.Header
}

// Comment
func (ctx *Writer) Write([]byte) (int, error) {
	return 0, nil
}

// Comment
func (ctx *Writer) WriteHeader(status int) {
	ctx.response.Status = strings.Join([]string{strconv.Itoa(status), StatusText(Status(status))}, "")
}

// Comment
func HttpResponse(protocol string, status Status, headers types.Headers, body []byte) *http.Response {
	return &http.Response{
		Proto:      protocol,
		StatusCode: int(status),
		Status:     strings.Join([]string{strconv.Itoa(int(status)), StatusText(status)}, " "),
		Header:     h.ToHeader(headers),
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
}

// Comment
func NewResponse(protocol string, status Status, headers types.Headers, body []byte) *Response {
	res := &Response{
		Bag:      &Bag{},
		Response: HttpResponse(protocol, status, headers, body),
	}

	res.Writer = &Writer{response: res}

	return res
}

// Comment
func InitResponse() *Response {
	return NewResponse("HTTP/1.1", HTTP_RESPONSE_OK, make(types.Headers), []byte{})
}

// Comment
func (ctx *Response) SetStatus(status Status) *Response {
	ctx.Status = strings.Join([]string{strconv.Itoa(int(status)), StatusText(status)}, " ")
	ctx.StatusCode = int(status)

	return ctx
}

// Comment
func (ctx *Response) SetHeader(key string, value string) *Response {
	ctx.Header.Set(cases.Title(language.English).String(key), value)

	return ctx
}

// Comment
func (ctx *Response) SetHeaders(headers types.Headers) *Response {
	for k, v := range headers {
		ctx.SetHeader(k, v)
	}

	return ctx
}

// Comment
func (ctx *Response) GetHeader(key string) string {
	header, ok := ctx.Header[cases.Title(language.English).String(key)]

	if !ok {
		return ""
	}

	return strings.Join(header, ",")
}

// Comment
func (ctx *Response) WithError(key string, value string) *Response {
	if ctx.Session != nil {
		ctx.Session.SetError(key, value)
	}

	return ctx
}

// Comment
func (ctx *Response) WithErrors(errors SessionErrorsBag) *Response {
	if ctx.Session != nil {
		ctx.Session.SetErrors(errors)
	}

	return ctx
}

// Comment
func (ctx *Response) SetBody(body []byte) *Response {
	ctx.Body = io.NopCloser(bytes.NewReader(body))

	return ctx
}

// Comment
func (ctx *Response) Html(html string) *Response {
	return ctx.SetHeader("content-type", "text/html").SetBody([]byte(html))
}

// Comment
func (ctx *Response) Json(v any) *Response {
	ctx.SetHeader("content-type", "application/json")

	data, err := json.Marshal(v)

	if err != nil {
		data, _ := json.Marshal(map[string]string{
			"message": "parse error",
		})

		return ctx.SetBody(data)
	}

	return ctx.SetBody(data)
}

// Comment
func (ctx *Response) Back() *Response {
	ctx.Bag.Redirect = &RedirectBag{To: ctx.Request.Header.Get("Referer")}

	return ctx.SetStatus(HTTP_RESPONSE_TEMPORARY_REDIRECT).Html(pages.Redirect(ctx.Bag.Redirect.To))
}

// Comment
func (ctx *Response) Redirect(path string) *Response {
	ctx.Bag.Redirect = &RedirectBag{To: strings.Trim(path, "/")}

	return ctx.SetStatus(HTTP_RESPONSE_TEMPORARY_REDIRECT).Html(pages.Redirect(ctx.Bag.Redirect.To))
}

// Comment
func (ctx *Response) Download(contentType string, filename string, binary []byte) *Response {
	return ctx.SetHeaders(types.Headers{
		"Content-Disposition": "attachment; filename=\"" + filename + "\"",
		"Content-Type":        contentType}).
		SetBody(binary)
}

// Comment
func (ctx *Response) View(view string, data ViewData) *Response {
	ctx.Bag.View = &ViewBag{
		Name: view,
		Data: data,
	}

	html, err := ctx.Request.
		Server.
		Get("view").(*View).
		Read(ctx.Bag.View.Name, ctx.Bag.View.Data, ctx.Request)

	if err != nil {
		return ctx.SetStatus(HTTP_RESPONSE_INTERNAL_SERVER_ERROR).Html(err.Error())
	}

	return ctx.Html(string(html))
}

// Comment
func HttpToResponse(text string) (*Response, error) {
	protocol, status, headers, body, err := response.ParseHttpToResponse(text)

	if err != nil {
		return nil, err
	}

	return NewResponse(protocol, Status(status), headers, body), nil
}

// Comment
func StatusText(status Status) string {
	switch status {
	case HTTP_RESPONSE_CONTINUE:
		return "Continue"
	case HTTP_RESPONSE_SWITCHING_PROTOCOLS:
		return "Switching Protocols"
	case HTTP_RESPONSE_PROCESSING:
		return "Processing"
	case HTTP_RESPONSE_EARLY_HINTS:
		return "Early Hints"
	case HTTP_RESPONSE_OK:
		return "Ok"
	case HTTP_RESPONSE_CREATED:
		return "Created"
	case HTTP_RESPONSE_ACCEPTED:
		return "Accepted"
	case HTTP_RESPONSE_NON_AUTHORITATIVE_INFORMATION:
		return "Non-Authoritative Information"
	case HTTP_RESPONSE_NO_CONTENT:
		return "No Content"
	case HTTP_RESPONSE_RESET_CONTENT:
		return "Reset Content"
	case HTTP_RESPONSE_PARTIAL_CONTENT:
		return "Partial Content"
	case HTTP_RESPONSE_MULTI_STATUS:
		return "Multi-Status"
	case HTTP_RESPONSE_ALREADY_REPORTED:
		return "Already Reported"
	case HTTP_RESPONSE_IM_USED:
		return "IM Used"
	case HTTP_RESPONSE_MULTIPLE_CHOICES:
		return "Multiple Choices"
	case HTTP_RESPONSE_MOVE_PERMANENTLY:
		return "Moved Permanently"
	case HTTP_RESPONSE_FOUND:
		return "Found"
	case HTTP_RESPONSE_SEE_OTHER:
		return "See Other"
	case HTTP_RESPONSE_NOT_MODIFIED:
		return "Not Modified"
	case HTTP_RESPONSE_USE_PROXY:
		return "Use Proxy"
	case HTTP_RESPONSE_UNUSED:
		return "unused"
	case HTTP_RESPONSE_TEMPORARY_REDIRECT:
		return "Temporary Redirect"
	case HTTP_RESPONSE_PERMANENT_REDIRECT:
		return "Permanent Redirect"
	case HTTP_RESPONSE_BAD_REQUEST:
		return "Bad Request"
	case HTTP_RESPONSE_UNAUTHORIZED:
		return "Unauthorized"
	case HTTP_RESPONSE_PAYMENT_REQUIRED:
		return "Payment Required"
	case HTTP_RESPONSE_FORBIDDEN:
		return "Forbidden"
	case HTTP_RESPONSE_NOT_FOUND:
		return "Not Found"
	case HTTP_RESPONSE_METHOD_NOT_ALLOWED:
		return "Method Not Allowed"
	case HTTP_RESPONSE_NOT_ACCEPTABLE:
		return "Not Acceptable"
	case HTTP_RESPONSE_PROXY_AUTHENTICATION_REQUIRED:
		return "Proxy Authentication Required"
	case HTTP_RESPONSE_REQUEST_TIMEOUT:
		return "Request Timeout"
	case HTTP_RESPONSE_CONFLICT:
		return "Conflict"
	case HTTP_RESPONSE_GONE:
		return "Gone"
	case HTTP_RESPONSE_LENGTH_REQUIRED:
		return "Length Required"
	case HTTP_RESPONSE_PRECONDITION_FAILED:
		return "Precondition Failed"
	case HTTP_RESPONSE_CONTENT_TOO_LARGER:
		return "Content Too Large"
	case HTTP_RESPONSE_URI_TOO_LARGE:
		return "URI Too Long"
	case HTTP_RESPONSE_UNSUPPORTED_MEDIA_TYPE:
		return "Unsupported Media Type"
	case HTTP_RESPONSE_RANGE_NOT_SATISFIABLE:
		return "Range Not Satisfiable"
	case HTTP_RESPONSE_EXPECTATION_FAILED:
		return "Expectation Failed"
	case HTTP_RESPONSE_IM_A_TEAPOT:
		return "I`m a teapot"
	case HTTP_RESPONSE_MISDIRECTED_REQUEST:
		return "Misdirected Request"
	case HTTP_RESPONSE_UNPROCESSABLE_CONTENT:
		return "Unprocessable Content"
	case HTTP_RESPONSE_LOCKED:
		return "Locked"
	case HTTP_RESPONSE_FAILED_DEPENDENCY:
		return "Failed Dependency"
	case HTTP_RESPONSE_TOO_EARLY:
		return "Too Early"
	case HTTP_RESPONSE_UPGRADE_REQUIRED:
		return "Upgrade Required"
	case HTTP_RESPONSE_PRECONDITION_REQUIRED:
		return "Precondition Required"
	case HTTP_RESPONSE_TOO_MANY_REQUIRED:
		return "Too Many Requests"
	case HTTP_RESPONSE_REQUEST_HEADER_FIELD_TOO_LARGE:
		return "Request Header Field Too Large"
	case HTTP_RESPONSE_UNAVAILABLE_FOR_LEGAL_REASONS:
		return "Unavailable For Legal Reasons"
	case HTTP_RESPONSE_INTERNAL_SERVER_ERROR:
		return "Internal Server Error"
	case HTTP_RESPONSE_NOT_IMPLEMENTED:
		return "Not Implemented"
	case HTTP_RESPONSE_BAD_GATEWAY:
		return "Bad Gateway"
	case HTTP_RESPONSE_SERVICE_UNAVAILABLE:
		return "Service Unavailable"
	case HTTP_RESPONSE_GATEWAY_TIMEOUT:
		return "Gateway Timeout"
	case HTTP_RESPONSE_HTTP_VERSION_NOT_SUPPORTED:
		return "HTTP Version Not Supported"
	case HTTP_RESPONSE_VARIANT_ALSO_NEGOTIATES:
		return "Variant Also Negotiates"
	case HTTP_RESPONSE_INSUFFICIENT_STORAGE:
		return "Insufficient Storage"
	case HTTP_RESPONSE_LOOP_DETECTED:
		return "Loop Detected"
	case HTTP_RESPONSE_NOT_EXTENDED:
		return "Not Extended"
	case HTTP_RESPONSE_NETWORK_AUTHENTICATION_REQUIRED:
		return "Network Authentication Required"
	default:
		return "Ok"
	}
}
