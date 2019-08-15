package squid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/flosch/pongo2"
	"github.com/go-session/session"
	"io"
	"net/http"
)

type Context struct {
	Response 	http.ResponseWriter
	Request 	*http.Request
}

type Controller map[string]func(ctx Context)

func (ctx Context) Get(name string) string {
	if ctx.Request.Method == "POST" {
		return ctx.Request.PostFormValue(name)
	}
	return ctx.Request.FormValue(name)
}

func (ctx Context) Render(file string, data Params) error {
	tpl := pongo2.Must(pongo2.FromFile(file))
	if data == nil {
		data = Params{}
	}
	return tpl.ExecuteWriter(data.Update(), ctx.Response)
}

func (ctx Context) Json(data Params) {
	str, err := json.Marshal(data)
	if err != nil {
		_, _ = ctx.Response.Write([]byte(err.Error()))
		return
	}
	ctx.Response.Header().Set("Content-Type", "application/json")
	_, _ = ctx.Response.Write(str)
}

func (ctx Context) Redirect(url string) {
	http.Redirect(ctx.Response, ctx.Request, url, http.StatusTemporaryRedirect)
}

func (ctx Context) SetSession(key string, value interface{}) error {
	manager := session.NewManager(session.SetExpired(60 * 60 * 24 * 7))
	store, err := manager.Start(context.Background(), ctx.Response, ctx.Request)
	//store, err := session.Start(context.Background(), ctx.Response, ctx.Request)
	if err != nil {
		return err
	}
	store.Set(key, value)
	err = store.Save()
	if err != nil {
		return err
	}
	return nil
}

func (ctx Context) GetSession(key string) (string, error) {
	store, err := session.Start(context.Background(), ctx.Response, ctx.Request)
	if err != nil {
		return err.Error(), err
	}
	value, ok := store.Get(key)
	if ok {
		return value.(string), nil
	}
	return "", errors.New("session does not exist")
}

func (ctx Context) FlushSession() {
	store, err := session.Start(context.Background(), ctx.Response, ctx.Request)
	if err != nil {
		return
	}
	_ = store.Flush()
}

func newUUID() string {
	var buf [16]byte
	_, _ = io.ReadFull(rand.Reader, buf[:])
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80

	dst := make([]byte, 36)
	hex.Encode(dst, buf[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], buf[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], buf[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], buf[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], buf[10:])

	return string(dst)
}