package http

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	str "strings"

	"github.com/gorilla/sessions"
	"github.com/lucas11776-golang/http/utils/strings"
	"github.com/spf13/cast"
)

// TODO: Many have errors saving session.Save

const SESSION_DEFAULT_EXPIRE = (60 * 60) * 24

type SessionErrorsBag map[string]string
type SessionOldBag map[string]string

// TODO: temp session remove for better version.
const (
	ERROR_KEY_STORE_KEY   = "__ERROR__STORE__KEY__"
	ERROR_KEY_REQUEST_KEY = "__ERROR__REQUEST__KEY__"
	CSFR_KEY              = "__CSRF__KEY__"
	OLD_STORE_KEY         = "__OLD__FORM__VALUES__STORE_KEY__"
	OLD_REQUEST_KEY       = "__OLD__FORM__VALUES__REQUEST__KEY__"
	CSRF_INPUT_NAME       = "CSRF_TOKEN"
)

type SessionManager interface {
	Set(key string, value interface{}) SessionManager
	Get(key string) string
	Clear() SessionManager
	Path(path string) SessionManager
	Remove(key string) SessionManager
	CanSave() bool
	Save() SessionManager
	SetError(key string, value string) SessionManager
	SetErrors(errors SessionErrorsBag) SessionManager
	Errors() SessionErrorsBag
	Error(key string) string
	Csrf() string
	Old(key string) string
}

type SessionsManager interface {
	Session(req *Request) SessionManager
	MaxAge(seconds int) SessionsManager
	Secure(secure bool) SessionsManager
	Domain(domain string) SessionsManager
	HttpOnly(httpOnly bool) SessionsManager
	SameSite(sameSite bool) SessionsManager
}

type Sessions struct {
	store *sessions.CookieStore
	name  string
}

type Session struct {
	session *sessions.Session
	request *Request
	save    bool
}

// Comment
func InitSession(name string, key []byte) *Sessions {
	s := sessions.NewCookieStore(key)

	s.Options = &sessions.Options{
		MaxAge: SESSION_DEFAULT_EXPIRE,
	}

	return &Sessions{name: name, store: s}
}

func (ctx *Session) newCsrf() *Session {
	ctx.session.Values[CSFR_KEY] = fmt.Sprintf("%d-%s", time.Now().Add(time.Minute*10).Unix(), strings.Random(50))

	ctx.save = true

	return ctx
}

// Comment
func (ctx *Session) initCsrf() *Session {
	csrf, ok := ctx.session.Values[CSFR_KEY].(string)

	if !ok {
		return ctx.newCsrf()
	}

	token := str.Split(csrf, "-")

	if len(token) != 2 {
		return ctx.newCsrf()
	}

	if t := cast.ToInt64(token[0]); t == 0 || time.Now().Unix() > t {
		return ctx.newCsrf()
	}

	return ctx
}

// Comment
func (ctx *Session) initErrors() *Session {
	data, ok := ctx.session.Values[ERROR_KEY_STORE_KEY].(string)

	if !ok {
		data = ""
	}

	errs := make(SessionErrorsBag)

	json.Unmarshal([]byte(data), &errs)

	ctx.session.Values[ERROR_KEY_REQUEST_KEY] = errs

	if len(errs) != 0 {
		ctx.save = true
	}

	return ctx
}

// Comment
func (ctx *Session) initOld() *Session {
	values, ok := ctx.session.Values[OLD_STORE_KEY]

	if !ok {
		return ctx
	}

	form := SessionOldBag{}

	json.Unmarshal([]byte(values.(string)), &form)

	ctx.session.Values[OLD_REQUEST_KEY] = form

	if len(form) != 0 {
		ctx.save = true
	}

	return ctx
}

// Comment
func (ctx *Sessions) Session(req *Request) SessionManager {
	session, _ := ctx.store.Get(req.Request, ctx.name)

	return (&Session{session: session, request: req}).initCsrf().initErrors().initOld()
}

// Comment
func (ctx *Sessions) MaxAge(seconds int) SessionsManager {
	ctx.store.Options.MaxAge = seconds

	return ctx
}

// Comment
func (ctx *Sessions) Secure(secure bool) SessionsManager {
	ctx.store.Options.Secure = secure

	return ctx
}

// Comment
func (ctx *Sessions) Domain(domain string) SessionsManager {
	ctx.store.Options.Domain = domain

	return ctx
}

// Comment
func (ctx *Sessions) HttpOnly(httpOnly bool) SessionsManager {
	ctx.store.Options.HttpOnly = httpOnly

	return ctx
}

// Comment
func (ctx *Sessions) SameSite(sameSite bool) SessionsManager {
	if sameSite {
		ctx.store.Options.SameSite = 1
	} else {
		ctx.store.Options.SameSite = 0
	}

	return ctx
}

// Comment
func (ctx *Session) Path(path string) SessionManager {
	ctx.session.Options.Path = path

	return ctx
}

type SessionBag map[string]interface{}

// Comment
func (ctx *Session) Set(key string, value interface{}) SessionManager {
	ctx.session.Values[key] = cast.ToString(value)

	ctx.save = true

	return ctx
}

// Comment
func (ctx *Session) Get(key string) string {
	value, ok := ctx.session.Values[key]

	if !ok {
		return ""
	}

	return value.(string)
}

// Comment
func (ctx *Session) Clear() SessionManager {
	for k := range ctx.session.Values {
		delete(ctx.session.Values, k)
	}

	ctx.save = true

	return ctx
}

// Comment
func (ctx *Session) Remove(key string) SessionManager {
	delete(ctx.session.Values, key)

	ctx.save = true

	return ctx
}

// Comment
func (ctx *Session) CanSave() bool {
	return ctx.save
}

// Comment
func (ctx *Session) stringflyErrors() *Session {
	errors, _ := json.Marshal(ctx.session.Values[ERROR_KEY_STORE_KEY])

	ctx.session.Values[ERROR_KEY_STORE_KEY] = string(errors)

	delete(ctx.session.Values, ERROR_KEY_REQUEST_KEY)

	return ctx
}

func (ctx *Session) saveFormValues(values url.Values) {
	formValues := map[string]string{}

	for k, v := range values {
		formValues[k] = v[0]
	}

	form, _ := json.Marshal(formValues)
	ctx.session.Values[OLD_STORE_KEY] = string(form)
}

// comment
func (ctx *Session) clearCache() *Session {
	delete(ctx.session.Values, ERROR_KEY_REQUEST_KEY)
	delete(ctx.session.Values, OLD_REQUEST_KEY)

	return ctx
}

// Comment
func (ctx *Session) Save() SessionManager {
	if !ctx.CanSave() {
		return ctx
	}

	if errors, ok := ctx.session.Values[ERROR_KEY_STORE_KEY].(SessionErrorsBag); ok && len(errors) != 0 {
		if values := ctx.request.Form; values != nil {
			ctx.saveFormValues(values)
		}
	}

	ctx.stringflyErrors().session.Save(ctx.request.Request, ctx.request.Response.Writer)

	return ctx.clearCache()
}

// Comment
func (ctx *Session) SetError(key string, value string) SessionManager {
	if _, ok := ctx.session.Values[ERROR_KEY_STORE_KEY]; !ok {
		ctx.session.Values[ERROR_KEY_STORE_KEY] = make(SessionErrorsBag)
	}

	ctx.session.Values[ERROR_KEY_STORE_KEY].(SessionErrorsBag)[key] = value

	ctx.save = true

	return ctx
}

// Comment
func (ctx *Session) SetErrors(errors SessionErrorsBag) SessionManager {
	for k, v := range errors {
		ctx.SetError(k, v)
	}

	return ctx
}

// Comment
func (ctx *Session) Errors() SessionErrorsBag {
	return ctx.session.Values[ERROR_KEY_REQUEST_KEY].(SessionErrorsBag)
}

// Comment
func (ctx *Session) Error(key string) string {
	err, ok := ctx.session.Values[ERROR_KEY_REQUEST_KEY].(SessionErrorsBag)[key]

	if !ok {
		return ""
	}

	return err
}

// Comment
func (ctx *Session) Csrf() string {
	csrf, ok := ctx.session.Values[CSFR_KEY].(string)

	if !ok {
		return ""
	}

	token := str.Split(csrf, "-")

	if len(token) != 2 {
		return ""
	}

	return fmt.Sprintf(`<input name="%s" value="%s">`, CSRF_INPUT_NAME, token[1])
}

// Comment
func (ctx *Session) Old(key string) string {
	old, ok := ctx.session.Values[OLD_REQUEST_KEY].(SessionOldBag)[key]

	if !ok {
		return ""
	}

	return old
}

// Comment
func SessionValue(req *Request) func(key string) string {
	return func(key string) string {
		return req.Session.Get(key)
	}
}

// Comment
func SessionHas(req *Request) func(key string) bool {
	return func(key string) bool {
		return req.Session.Error(key) != ""
	}
}

// Comment
func SessionError(req *Request) func(key string) string {
	return func(key string) string {
		return req.Session.Error(key)
	}
}

// Comment
func SessionErrors(req *Request) func() SessionErrorsBag {
	return func() SessionErrorsBag {
		return req.Session.Errors()
	}
}

// Comment
func SessionCsrf(req *Request) func() string {
	return func() string {
		return req.Session.Csrf()
	}
}

// Comment
func SessionOld(req *Request) func(key string) string {
	return func(key string) string {
		return req.Session.Old(key)
	}
}
