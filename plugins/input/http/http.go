package http

import (
	"github.com/luopengift/gohttp"
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
)

const (
	VERSION = "0.0.1"
)

var Ch chan []byte

type HttpInput struct {
	Addr string `json:"addr"`
	Path string `json:"path"`

	app *gohttp.Application
	gohttp.BaseHTTPHandler
}

func NewHttpInput() *HttpInput {
	return new(HttpInput)
}

func (in *HttpInput) Init(config transport.Configer) error {
	err := config.Parse(in)
	if err != nil {
		logger.Error("parse error:%v", err)
		return err
	}
	in.app = gohttp.Init()
	in.app.Route("^"+in.Path+"$", &HttpInput{})
	Ch = make(chan []byte, 100)
	return nil
}

func (in *HttpInput) POST() {
	body := in.BaseHTTPHandler.GetBodyArgs()
	Ch <- body
	logger.Error("body:%v,ok", string(body))
	in.Output("ok")
}

func (in *HttpInput) Read(p []byte) (int, error) {
	n := copy(p, <-Ch)
	return n, nil
}

func (in *HttpInput) Start() error {
	in.app.Run(in.Addr)
	return nil
}

func (in *HttpInput) Close() error {
	close(Ch)
	return nil
}

func (in *HttpInput) Version() string {
	return VERSION
}

func init() {
	transport.RegistInputer("http", NewHttpInput())
}
