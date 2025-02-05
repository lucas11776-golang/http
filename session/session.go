package session

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/lucas11776-golang/http/request"
)

type Sessions struct {
	store *sessions.CookieStore
	name  string
}

type Session struct {
	session *sessions.Session
	request *request.Request
}

// Comment
func Init(name string, key []byte) *Sessions {
	return &Sessions{
		name:  name,
		store: sessions.NewCookieStore(key),
	}
}

// Comment
func (ctx *Sessions) Session(req *request.Request) (*Session, error) {
	session, err := ctx.store.Get(req.Request, ctx.name)

	if err != nil {
		return nil, err
	}

	return &Session{
		session: session,
		request: req,
	}, nil
}

type W http.ResponseWriter

// Comment
func (ctx *Session) Set(key string, value string) *Session {
	ctx.session.Values[key] = value

	// ctx.session.Save(ctx.request.Request, ctx.request.Response)

	return ctx
}
