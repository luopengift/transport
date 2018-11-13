package api

import (
	"github.com/luopengift/gohttp"
	"github.com/luopengift/log"
	"github.com/luopengift/transport"
)

type RootHandler struct {
	gohttp.BaseHTTPHandler
}

func (ctx *RootHandler) GET() {
	ctx.Output("root")
}

type StatsHandler struct {
	gohttp.BaseHTTPHandler
}

func (ctx *StatsHandler) GET() {
	stats := transport.T.Stat()
	log.Debug("[API] %v", stats)
	ctx.Output(stats)
}

type StoreHandler struct {
	gohttp.BaseHTTPHandler
}

func APIHttp(addr string) {
	app := gohttp.Init()
	app.Route("^/$", &RootHandler{})
	app.Route("^/stats$", &StatsHandler{})
	go app.Run(addr)
}
