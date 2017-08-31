package api

import (
	"github.com/luopengift/gohttp"
	"github.com/luopengift/transport"
)

type RootHandler struct {
	gohttp.HttpHandler
}

func (ctx *RootHandler) GET() {
	ctx.Output("root")
}

type StatsHandler struct {
	gohttp.HttpHandler
}

func (ctx *StatsHandler) GET() {
	ctx.Output(transport.T.Stat())
}

type StoreHandler struct {
	gohttp.HttpHandler
}

func ApiHttp(addr string) {
	app := gohttp.Init()
	app.Route("^/$", &RootHandler{})
	app.Route("^/stats$", &StatsHandler{})
	go app.Run(addr)
}
