package squid

import (
	"encoding/json"
	"github.com/flosch/pongo2"
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