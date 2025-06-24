package http

import (
	"strings"

	"github.com/lucas11776-golang/http"
	"github.com/lucas11776-golang/http/pages"
)

// Comment
func CsrfMiddleware(req *http.Request, res *http.Response, next http.Next) *http.Response {
	switch http.Method(strings.ToUpper(req.Method)) {
	case http.METHOD_POST, http.METHOD_PATCH, http.METHOD_PUT, http.METHOD_DELETE:
		if req.Session == nil {
			return next()
		}

		if req.FormValue(http.CSRF_INPUT_NAME) != req.Session.CsrfToken() {
			return res.Html(pages.CsrfExpired())
		}

		return next()

	default:
		return next()
	}

}
