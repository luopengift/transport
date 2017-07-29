package api

import (
	"github.com/luopengift/gohttp"
	"github.com/luopengift/transport/plugins/handler"
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
	ctx.Output("stats")
}

type StoreHandler struct {
	gohttp.HttpHandler
}

func (ctx *StoreHandler) GET() {
	ctx.Output(handler.Store)
}

func init() {
	app := gohttp.Init()
	app.Route("^/$", &RootHandler{})
	app.Route("^/stats$", &StatsHandler{})
	app.Route("^/store$", &StoreHandler{})
	go app.Run(":38888")
}
