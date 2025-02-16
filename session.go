package http

import (
	"log"

	"github.com/gorilla/sessions"
)

const SESSION_DEFAULT_EXPIRE = (60 * 60) * 24

type SessionManager interface {
	Set(key string, value string) SessionManager
	Get(key string) string
	Clear() SessionManager
	Path(path string) SessionManager
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
	session, err := ctx.store.Get(req.Request, ctx.name)

	if err != nil {
		log.Fatal(err.Error())
	}

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
	ctx.session.Options.MaxAge = -1

	return ctx
}

// Comment
func (ctx *Session) Save() SessionManager {
	ctx.session.Save(ctx.request.Request, ctx.request.Response.Writer)

	return ctx
}
