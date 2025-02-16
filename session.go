package http

import (
	"github.com/gorilla/sessions"
)

const SESSION_DEFAULT_EXPIRE = (60 * 60) * 24

type SessionManager interface {
	Set(key string, value string) SessionManager
	Get(key string) string
	Clear() SessionManager
	Path(path string) SessionManager
	Remove(key string) SessionManager
	CanSave() bool
	Save() SessionManager
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

	return &Sessions{
		name:  name,
		store: s,
	}
}

// Comment
func (ctx *Sessions) Session(req *Request) SessionManager {
	session, _ := ctx.store.Get(req.Request, ctx.name)

	return &Session{
		session: session,
		request: req,
	}
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
func (ctx *Session) Save() SessionManager {
	if ctx.CanSave() {
		ctx.session.Save(ctx.request.Request, ctx.request.Response.Writer)
	}

	return ctx
}
