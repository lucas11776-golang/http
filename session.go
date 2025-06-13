package http

import (
	"encoding/json"

	"github.com/gorilla/sessions"
)

const SESSION_DEFAULT_EXPIRE = (60 * 60) * 24

type Errors map[string]string

// TODO: temp session remove for better version.
const (
	ERROR_KEY_STORE   = "__ERROR__STORE__"
	ERROR_KEY_REQUEST = "__ERROR__REQUEST__"
)

type SessionManager interface {
	Set(key string, value string) SessionManager
	Get(key string) string
	Clear() SessionManager
	Path(path string) SessionManager
	Remove(key string) SessionManager
	CanSave() bool
	Save() SessionManager
	SetError(value string, error string) SessionManager
	Error(value string) string
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

// Comment
func (ctx *Sessions) Session(req *Request) SessionManager {
	session, _ := ctx.store.Get(req.Request, ctx.name)

	return (&Session{session: session, request: req}).decodeErrors()
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
func (ctx *Session) Set(key string, value string) SessionManager {
	ctx.session.Values[key] = value

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
	errors, _ := json.Marshal(ctx.session.Values[ERROR_KEY_STORE])

	ctx.session.Values[ERROR_KEY_STORE] = string(errors)

	delete(ctx.session.Values, ERROR_KEY_REQUEST)

	return ctx
}

// Comment
func (ctx *Session) decodeErrors() *Session {
	data, ok := ctx.session.Values[ERROR_KEY_STORE].(string)

	if !ok {
		data = ""
	}

	errs := make(Errors)

	json.Unmarshal([]byte(data), &errs)

	ctx.session.Values[ERROR_KEY_REQUEST] = errs

	return ctx
}

// Comment
func (ctx *Session) Save() SessionManager {
	if ctx.CanSave() {
		ctx.stringflyErrors().session.Save(ctx.request.Request, ctx.request.Response.Writer)
	}

	return ctx
}

// Comment
func (ctx *Session) SetError(key string, err string) SessionManager {
	if _, ok := ctx.session.Values[ERROR_KEY_STORE]; !ok {
		ctx.session.Values[ERROR_KEY_STORE] = make(Errors)
	}

	ctx.session.Values[ERROR_KEY_STORE].(Errors)[key] = err

	ctx.save = true

	return ctx
}

// Comment
func (ctx *Session) Error(value string) string {
	err, ok := ctx.session.Values[ERROR_KEY_REQUEST].(Errors)[value]

	if !ok {
		return ""
	}

	return err
}
