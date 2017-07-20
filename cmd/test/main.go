package main

import (
	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/transport"
	//_ "github.com/luopengift/transport/api"
	"github.com/luopengift/transport/config"
	_ "github.com/luopengift/transport/plugins"
	"runtime"
	//"flag"
)

const (
	VERSION = "0.0.1"
)

var t *transport.Transport

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger.Info("Transport starting...")
	cfg := config.NewConfig()
	logger.Info("%#v,%#v, %#v", transport.InputPlugins, transport.HandlePlugins, transport.OutputPlugins)

	input := cfg.Input()
	output := cfg.Output()
	handle := cfg.Handle()
	logger.Debug("%#v,%#v,%#v", input, handle, output)
	t = transport.NewTransport(input, handle, output)
	t.Run()
	logger.Debug("%#v", t)

}
