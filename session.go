package squid

import (
	"net/http"
	"time"
)

const (
	cookieName = "squid_session_name"
	expired = 60 * 60 * 24 * 7
)

type Session struct {
	Context		Context
}

type list map[string]string
var storage = map[string]list{}

func (s Session) Set(key, value string) {
	cookie, err := s.Context.Request.Cookie(cookieName)
	if err == nil && len(cookie.Value) > 0 {
		s.setCookie(cookie.Value)
		storage[cookie.Value] = list{key: value}
		return
	}
	sid := newUUID()
	s.setCookie(sid)
	storage[sid] = list{key: value}
}

func (s Session) Get(key string) (string, error) {
	cookie, err := s.Context.Request.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	return storage[cookie.Value][key], nil
}

func (s Session) setCookie(sid string) {
	c := &http.Cookie{
		Name: cookieName,
		Value: sid,
		Path: "/",
		HttpOnly: true,
		Secure: false,
		MaxAge: int(expired),
		Expires: time.Now().Add(time.Duration(expired) * time.Second),
	}
	http.SetCookie(s.Context.Response, c)
	s.Context.Request.AddCookie(c)
}

func (s Session) Flush() {
	cookie, err := s.Context.Request.Cookie(cookieName)
	if err != nil {
		return
	}
	delete(storage, cookie.Value)
}