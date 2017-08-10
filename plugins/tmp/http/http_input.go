package http

import (
	"github.com/luopengift/gohttp"
	"github.com/luopengift/transport/pipeline"
)

var ch = make(chan []byte, 100)

type HttpInput struct {
	app *gohttp.Application
	gohttp.HttpHandler
}

func (http *HttpInput) POST() {
	ch <- http.HttpHandler.GetBodyArgs()
	http.Output("ok")
}

func NewHttpInput() *HttpInput {
	http := new(HttpInput)
	http.app = gohttp.Init()
	return http
}

func (http *HttpInput) Init(cfg map[string]string) error {
	http.app.Config.SetAddress(cfg["addr"])
	http.app.Route("/post", &HttpInput{})
	return nil
}

func (http *HttpInput) Start() error {
	http.app.Run()
	return nil
}

func (http *HttpInput) Read(p []byte) (int, error) {
	n := copy(p, <-ch)
	return n, nil
}

func (http *HttpInput) Close() error {
	return nil
}

func init() {
	pipeline.RegistInputer("http", NewHttpInput())
}
