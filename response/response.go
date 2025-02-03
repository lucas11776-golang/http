package response

import (
	"sort"
	"strconv"
	"strings"

	"github.com/lucas11776-golang/http/request"
	"github.com/lucas11776-golang/http/types"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Status int16

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

type ResponseType string

const (
// RESPONSE_TYPE_NEXT     ResponseType = "next"
// RESPONSE_TYPE_DATA     ResponseType = "data"
// RESPONSE_TYPE_VIEW     ResponseType = "view"
// RESPONSE_TYPE_REDIRECT ResponseType = "redirect"
// RESPONSE_TYPE_DOWNLOAD ResponseType = "download"
)

type Response struct {
	format     ResponseType
	protocol   string
	status     Status
	statusText string
	headers    types.Headers
	body       []byte
	Request    *request.Request
}

// Comment
func Create(protocol string, status Status, headers types.Headers, body []byte) *Response {
	return &Response{
		protocol: protocol,
		status:   status,
		headers:  headers,
		body:     body,
	}
}

// Comment
func Init() *Response {
	return &Response{
		status:   200,
		protocol: "HTTP/1.1",
		headers:  make(types.Headers),
	}
}

// Comment
func (ctx *Response) Status(status Status) *Response {
	ctx.status = status

	return ctx
}

// Comment
func (ctx *Response) Header(key string, value string) *Response {
	ctx.headers[key] = value

	return ctx
}

// Comment
func (ctx *Response) Body(body []byte) *Response {
	return BodyDefault(ctx, body)
}

// Comment
func (ctx *Response) Html(html string) *Response {
	return BodyHtml(ctx, html)
}

// Comment
func (ctx *Response) Json(v any) *Response {
	return BodyJson(ctx, v)
}

// Comment
func (ctx *Response) Redirect(path string) *Response {
	return BodyRedirect(ctx, path)
}

// Comment
func (ctx *Response) Download(contentType string, filename string, binary []byte) *Response {
	return BodyDownload(ctx, contentType, filename, binary)
}

// Comment
func ParseHttp(res *Response) string {
	http := []string{}

	http = append(http, strings.Join([]string{res.protocol, strconv.Itoa(int(res.status)), getStatusText(res.status)}, " "))

	keys := make([]string, 0, len(res.headers))

	for k := range res.headers {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		http = append(http, strings.Join([]string{cases.Title(language.English).String(key), res.headers[key]}, ": "))
	}

	http = append(http, strings.Join([]string{"Content-Length", strconv.Itoa(len(res.body))}, ": "))

	if len(res.body) == 0 {
		return strings.Join(append(http, "\r\n"), "\r\n")
	}

	return strings.Join(append(http, strings.Join([]string{"\r\n", string(res.body), "\r\n"}, "")), "\r\n")
}

// Comment
func getStatusText(status Status) string {
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
		return ""
	}
}
