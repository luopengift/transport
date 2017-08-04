package api

import (
	"github.com/luopengift/gohttp"
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

func init() {
	app := gohttp.Init()
	app.Route("^/$", &RootHandler{})
	app.Route("^/stats$", &StatsHandler{})
	go app.Run(":38888")
}
