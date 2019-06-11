package squid

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/flosch/pongo2"
	"github.com/go-session/session"
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
	//option := session.SetCookieLifeTime(60 * 60 * 24 * 3)
	//manager := session.NewManager(option)
	store, err := session.Start(context.Background(), ctx.Response, ctx.Request)
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